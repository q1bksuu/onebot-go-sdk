package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/q1bksuu/onebot-go-sdk/v11/dispatcher"
	"github.com/q1bksuu/onebot-go-sdk/v11/entity"
	wsinternal "github.com/q1bksuu/onebot-go-sdk/v11/internal/ws"
	"github.com/q1bksuu/onebot-go-sdk/v11/server"
)

// WSClientConfig WebSocket 客户端配置.
type WSClientConfig struct {
	URL               string
	ReconnectInterval time.Duration // 断线重连间隔
	SelfID            int64         // 机器人 QQ 号（用于 X-Self-ID 请求头）
	AccessToken       string        // 可选鉴权令牌
	ReadTimeout       time.Duration // 读取超时（可选），默认 0
	WriteTimeout      time.Duration // 写入超时（可选），默认 0
}

// WSClientOption 用于配置 WebSocketClient 的选项函数类型.
type WSClientOption func(*WebSocketClient)

// WebSocketClient 实现 OneBot 反向 WebSocket 传输层.
type WebSocketClient struct {
	cfg           WSClientConfig
	actionHandler dispatcher.ActionRequestHandler

	// 连接管理
	mu   sync.Mutex
	conn *websocket.Conn

	// 控制
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewWebSocketClient 创建反向 WebSocket 客户端.
func NewWebSocketClient(cfg WSClientConfig, opts ...WSClientOption) *WebSocketClient {
	_, cancel := context.WithCancel(context.Background())

	client := &WebSocketClient{
		cfg:    cfg,
		cancel: cancel,
	}

	// 应用选项
	for _, opt := range opts {
		opt(client)
	}

	return client
}

// Start 启动客户端，建立连接并开始处理消息.
func (c *WebSocketClient) Start(ctx context.Context) error {
	mergedCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		<-mergedCtx.Done()
		c.cancel()
	}()

	if c.cfg.URL == "" {
		return server.ErrUniversalClientURLEmpty
	}

	c.wg.Add(1)

	go c.run(mergedCtx, c.cfg.URL)

	<-mergedCtx.Done()
	c.wg.Wait()

	return nil
}

// Shutdown 优雅关闭所有连接.
func (c *WebSocketClient) Shutdown(ctx context.Context) error {
	c.cancel()

	c.mu.Lock()
	conn := c.conn
	c.mu.Unlock()

	if conn != nil {
		_ = conn.Close()
	}

	done := make(chan struct{})

	go func() {
		c.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("shutdown timeout: %w", ctx.Err())
	}
}

// BroadcastEvent 推送事件到 Event 或 Universal 连接.
func (c *WebSocketClient) BroadcastEvent(event entity.Event) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var conns []*websocket.Conn

	if len(conns) == 0 {
		return nil
	}

	var lastErr error

	for _, conn := range conns {
		err := conn.WriteJSON(event)
		if err != nil {
			lastErr = err
		}
	}

	return lastErr
}

func (c *WebSocketClient) buildHeaders(clientRole string) http.Header {
	headers := make(http.Header)
	headers.Set("X-Self-Id", strconv.FormatInt(c.cfg.SelfID, 10))
	headers.Set("X-Client-Role", clientRole)

	if c.cfg.AccessToken != "" {
		headers.Set("Authorization", "Bearer "+c.cfg.AccessToken)
	}

	return headers
}

// dialWithReconnect 建立连接，支持断线重连.
func (c *WebSocketClient) dialWithReconnect(
	ctx context.Context, url string, headers http.Header,
) (*websocket.Conn, error) {
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	reconnectInterval := c.cfg.ReconnectInterval
	if reconnectInterval == 0 {
		reconnectInterval = 3 * time.Second
	}

	for {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("dial context canceled: %w", ctx.Err())
		}

		conn, resp, err := dialer.Dial(url, headers)
		if err == nil {
			_ = resp.Body.Close()

			return conn, nil
		}

		// 等待重连间隔
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("reconnect context canceled: %w", ctx.Err())
		case <-time.After(reconnectInterval):
		}
	}
}

// run 运行客户端.
func (c *WebSocketClient) run(ctx context.Context, url string) {
	defer c.wg.Done()

	for ctx.Err() == nil {
		headers := c.buildHeaders("Universal")

		conn, err := c.dialWithReconnect(ctx, url, headers)
		if err != nil {
			return
		}

		c.mu.Lock()
		c.conn = conn
		c.mu.Unlock()

		apiCtx, apiCancel := context.WithCancel(ctx)
		apiDone := make(chan struct{})

		go func() {
			defer close(apiDone)

			c.serveActionConn(apiCtx, conn)
		}()

		select {
		case <-ctx.Done():
			apiCancel()
			<-apiDone

			c.clearConn(conn)
			_ = conn.Close()

			return
		case <-apiDone:
			apiCancel()
		}

		c.clearConn(conn)
		_ = conn.Close()
	}
}

func (c *WebSocketClient) serveActionConn(ctx context.Context, conn *websocket.Conn) {
	if c.actionHandler == nil {
		<-ctx.Done()

		return
	}

	for ctx.Err() == nil {
		if c.cfg.ReadTimeout > 0 {
			err := conn.SetReadDeadline(time.Now().Add(c.cfg.ReadTimeout))
			if err != nil {
				return
			}
		}

		_, data, err := conn.ReadMessage()
		if err != nil {
			return
		}

		resp := c.handleActionMessage(ctx, data)

		if c.cfg.WriteTimeout > 0 {
			err := conn.SetWriteDeadline(time.Now().Add(c.cfg.WriteTimeout))
			if err != nil {
				return
			}
		}

		err = conn.WriteJSON(resp)
		if err != nil {
			return
		}
	}
}

func (c *WebSocketClient) handleActionMessage(ctx context.Context, data []byte) *entity.ActionResponseEnvelope {
	return wsinternal.HandleActionMessage(ctx, data, c.actionHandler, server.ErrBadRequest)
}

func (c *WebSocketClient) clearConn(conn *websocket.Conn) {
	c.mu.Lock()

	if c.conn == conn {
		c.conn = nil
	}

	c.mu.Unlock()
}
