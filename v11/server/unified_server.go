package server

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/q1bksuu/onebot-go-sdk/v11/dispatcher"
)

// UnifiedServer 支持同时运行 HTTP 和 WebSocket 的统一服务器.
type UnifiedServer struct {
	*BaseServer

	httpSrv *HTTPServer
	wsSrv   *WebSocketServer
}

// UnifiedConfig 统一服务器配置.
type UnifiedConfig struct {
	Addr string // 统一监听地址

	// HTTP 配置 (Addr 字段将被忽略)
	HTTP HTTPConfig
	// WebSocket 配置 (Addr 字段将被忽略)
	WS WSConfig

	// Server 通用配置
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

// NewUnifiedServer 创建统一服务器.
func NewUnifiedServer(
	cfg UnifiedConfig,
	httpOpts []HTTPServerOption,
	wsHandler dispatcher.ActionRequestHandler,
) *UnifiedServer {
	// 覆盖子配置的 Addr，虽然不会被底层 Server 用来监听，但为了保持一致性
	cfg.HTTP.Addr = cfg.Addr
	cfg.WS.Addr = cfg.Addr

	// 创建各个服务器实例
	httpSrv := NewHTTPServer(cfg.HTTP, httpOpts...)
	wsSrv := NewWebSocketServer(cfg.WS, wsHandler)

	// 使用 combinedHandler 分发请求
	// 这允许 HTTP 和 WebSocket 共用同一个端口和路径（例如 "/"），
	// 通过 Upgrade 头来区分协议。
	mainMux := &combinedHandler{
		http: httpSrv.Handler(),
		ws:   wsSrv.Handler(),
	}

	server := &UnifiedServer{
		httpSrv: httpSrv,
		wsSrv:   wsSrv,
	}

	baseCfg := ServerConfig{
		Addr:              cfg.Addr,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}
	server.BaseServer = NewBaseServer(baseCfg, mainMux)

	return server
}

// Start 启动统一服务器.
func (s *UnifiedServer) Start(ctx context.Context) error {
	return s.BaseServer.Start(ctx, func(ctx context.Context) error {
		// 同时确保 WS 连接被关闭 (WS Server 的 Shutdown 主要是关闭连接)
		// 注意：WS Server 的 Shutdown 也会尝试关闭它自己的 srv，但因为它的 srv 没有监听，所以应该没问题。
		// 不过，为了更干净，我们只做必要的清理：关闭连接。
		// 由于 WebSocketServer.Shutdown 内部调用了 closeAllConns 然后 srv.Shutdown
		// 我们可以直接调用它，反正 srv.Shutdown 对未启动的 server 是 no-op。
		return s.wsSrv.Shutdown(ctx)
	})
}

// Shutdown 关闭服务器.
func (s *UnifiedServer) Shutdown(ctx context.Context) error {
	_ = s.wsSrv.Shutdown(ctx)

	return s.BaseServer.Shutdown(ctx)
}

// combinedHandler 简单的组合处理器.
type combinedHandler struct {
	http http.Handler
	ws   http.Handler
}

func (h *combinedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 简单的协议探测
	// 如果是 WebSocket 握手请求，优先交给 WS 处理器
	if strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
		h.ws.ServeHTTP(w, r)

		return
	}

	// 否则默认交给 HTTP 处理器
	// 注意：如果 WS 配置了 Universal 路径 "/" 且 HTTP 也配置了 "/"，
	// 这里非 Upgrade 请求会去 HTTP。
	// 这通常是期望的行为：同一个端口，普通 HTTP 请求走 REST API，WS 握手走 WS。
	h.http.ServeHTTP(w, r)
}
