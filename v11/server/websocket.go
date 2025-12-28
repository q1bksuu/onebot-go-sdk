package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/q1bksuu/onebot-go-sdk/v11/internal/util"

	"github.com/q1bksuu/onebot-go-sdk/v11/entity"
)

// WSConfig 正向 WebSocket 服务配置.
type WSConfig struct {
	Addr         string                     // 监听地址，例 ":6700"
	PathPrefix   string                     // 路径前缀，可为空或"/"，最终用于 /api、/event、/ 路由
	AccessToken  string                     // 可选鉴权，若为空则不校验
	CheckOrigin  func(r *http.Request) bool // 可选跨域校验，默认全放行
	ReadTimeout  time.Duration              // 读取超时（可选），默认 0
	WriteTimeout time.Duration              // 写入超时（可选），默认 0
	IdleTimeout  time.Duration              // 空闲超时（可选），默认 0
}

type wsConn struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func (c *wsConn) writeJSON(v any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.conn.WriteJSON(v)
	if err != nil {
		return fmt.Errorf("write json failed: %w", err)
	}

	return nil
}

// WebSocketServer 实现 OneBot 正向 WebSocket 传输层.
type WebSocketServer struct {
	srv      *http.Server
	cfg      WSConfig
	handler  ActionRequestHandler
	upgrader websocket.Upgrader

	mu            sync.Mutex
	eventConns    map[*wsConn]struct{}
	universalConn map[*wsConn]struct{}
}

// NewWebSocketServer 创建 WebSocketServer. 若传入 CheckOrigin 为 nil，则允许任意来源.
func NewWebSocketServer(cfg WSConfig, handler ActionRequestHandler) *WebSocketServer {
	prefix := util.NormalizePath(cfg.PathPrefix)

	upgrader := websocket.Upgrader{CheckOrigin: cfg.CheckOrigin}
	if upgrader.CheckOrigin == nil {
		upgrader.CheckOrigin = func(*http.Request) bool { return true }
	}

	server := &WebSocketServer{
		cfg:           cfg,
		handler:       handler,
		upgrader:      upgrader,
		eventConns:    make(map[*wsConn]struct{}),
		universalConn: make(map[*wsConn]struct{}),
	}

	mux := http.NewServeMux()
	mux.HandleFunc(prefix+"/api", server.handleAPI)
	mux.HandleFunc(prefix+"/api/", server.handleAPI)
	mux.HandleFunc(prefix+"/event", server.handleEvent)
	mux.HandleFunc(prefix+"/event/", server.handleEvent)

	universalPath := prefix
	if universalPath == "" {
		universalPath = "/"
	}

	mux.HandleFunc(universalPath, server.handleUniversal)

	server.srv = &http.Server{
		Addr:         cfg.Addr,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return server
}

// Start 启动 WebSocket 服务器（异步监听）.
func (s *WebSocketServer) Start(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		err := s.srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		// 使用独立的上下文执行 shutdown，确保即使原上下文已取消也能完成关闭
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		//nolint:contextcheck // 需要独立的上下文来执行 shutdown
		return s.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}

// Shutdown 优雅关闭，关闭所有连接.
func (s *WebSocketServer) Shutdown(ctx context.Context) error {
	s.closeAllConns()

	err := s.srv.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	return nil
}

// BroadcastEvent 推送事件.
func (s *WebSocketServer) BroadcastEvent(event entity.Event) {
	s.mu.Lock()

	conns := make([]*wsConn, 0, len(s.eventConns)+len(s.universalConn))
	for conn := range s.eventConns {
		conns = append(conns, conn)
	}

	for conn := range s.universalConn {
		conns = append(conns, conn)
	}

	s.mu.Unlock()

	for _, conn := range conns {
		_ = conn.writeJSON(event)
	}
}

func (s *WebSocketServer) handleAPI(w http.ResponseWriter, r *http.Request) {
	if !s.matchPath(r.URL.Path, "/api") {
		http.NotFound(w, r)

		return
	}

	if errResp := s.checkAccess(r); errResp != nil {
		s.writeHandshakeError(w, errResp)

		return
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	s.serveActionConn(r.Context(), &wsConn{conn: conn}, false)
}

func (s *WebSocketServer) handleUniversal(w http.ResponseWriter, r *http.Request) {
	if !s.matchUniversalPath(r.URL.Path) {
		http.NotFound(w, r)

		return
	}

	if errResp := s.checkAccess(r); errResp != nil {
		s.writeHandshakeError(w, errResp)

		return
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	wsC := &wsConn{conn: conn}

	s.mu.Lock()
	s.universalConn[wsC] = struct{}{}
	s.mu.Unlock()

	s.serveActionConn(r.Context(), wsC, true)

	s.mu.Lock()
	delete(s.universalConn, wsC)
	s.mu.Unlock()

	_ = conn.Close()
}

func (s *WebSocketServer) handleEvent(w http.ResponseWriter, r *http.Request) {
	if !s.matchPath(r.URL.Path, "/event") {
		http.NotFound(w, r)

		return
	}

	if errResp := s.checkAccess(r); errResp != nil {
		s.writeHandshakeError(w, errResp)

		return
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	wsC := &wsConn{conn: conn}

	s.mu.Lock()
	s.eventConns[wsC] = struct{}{}
	s.mu.Unlock()

	// 仅保活，读取直到关闭.
	for {
		_, _, readErr := conn.ReadMessage()
		if readErr != nil {
			break
		}
	}

	s.mu.Lock()
	delete(s.eventConns, wsC)
	s.mu.Unlock()

	_ = conn.Close()
}

func (s *WebSocketServer) serveActionConn(ctx context.Context, wsC *wsConn, track bool) {
	for {
		_, data, err := wsC.conn.ReadMessage()
		if err != nil {
			break
		}

		resp := s.handleActionMessage(ctx, data)

		writeErr := wsC.writeJSON(resp)
		if writeErr != nil {
			break
		}
	}

	_ = wsC.conn.Close()

	if !track {
		return
	}
}

type actionRequestEnvelope struct {
	Action string          `json:"action"`
	Params map[string]any  `json:"params"`
	Echo   json.RawMessage `json:"echo,omitempty"`
}

type actionResponseEnvelope struct {
	Status  entity.ActionResponseStatus  `json:"status"`
	Retcode entity.ActionResponseRetcode `json:"retcode"`
	Data    json.RawMessage              `json:"data,omitempty"`
	Message string                       `json:"message,omitempty"`
	Echo    json.RawMessage              `json:"echo,omitempty"`
}

func (s *WebSocketServer) handleActionMessage(ctx context.Context, data []byte) *actionResponseEnvelope {
	var reqEnv actionRequestEnvelope

	err := json.Unmarshal(data, &reqEnv)
	if err != nil {
		return &actionResponseEnvelope{
			Status:  entity.StatusFailed,
			Retcode: entity.ActionResponseRetcode(1400),
			Message: "invalid json",
		}
	}

	req := &entity.ActionRequest{Action: reqEnv.Action, Params: reqEnv.Params}

	resp, err := s.handler.HandleActionRequest(ctx, req)
	if err != nil {
		mapped := mapHandlerError(err)

		return &actionResponseEnvelope{
			Status:  mapped.Status,
			Retcode: mapped.Retcode,
			Data:    mapped.Data,
			Message: mapped.Message,
			Echo:    reqEnv.Echo,
		}
	}

	if resp == nil {
		resp = &entity.ActionRawResponse{Status: entity.StatusFailed, Retcode: -1, Message: "empty response"}
	}

	return &actionResponseEnvelope{
		Status:  resp.Status,
		Retcode: resp.Retcode,
		Data:    resp.Data,
		Message: resp.Message,
		Echo:    reqEnv.Echo,
	}
}

func (s *WebSocketServer) checkAccess(r *http.Request) *actionResponseEnvelope {
	if s.cfg.AccessToken == "" {
		return nil
	}

	token := r.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	} else if token == "" {
		token = r.URL.Query().Get("access_token")
	}

	if token == "" {
		return &actionResponseEnvelope{Status: entity.StatusFailed, Retcode: 1401, Message: "missing access token"}
	}

	if token != s.cfg.AccessToken {
		return &actionResponseEnvelope{Status: entity.StatusFailed, Retcode: 1403, Message: "forbidden"}
	}

	return nil
}

func (s *WebSocketServer) matchPath(path, suffix string) bool {
	normalized := util.NormalizePath(s.cfg.PathPrefix)
	base := normalized + suffix

	return path == base || path == base+"/"
}

func (s *WebSocketServer) matchUniversalPath(path string) bool {
	normalized := util.NormalizePath(s.cfg.PathPrefix)
	if normalized == "" {
		return path == "/"
	}

	return path == normalized || path == normalized+"/"
}

func (s *WebSocketServer) closeAllConns() {
	s.mu.Lock()

	conns := make([]*wsConn, 0, len(s.eventConns)+len(s.universalConn))
	for c := range s.eventConns {
		conns = append(conns, c)
	}

	for c := range s.universalConn {
		conns = append(conns, c)
	}

	s.mu.Unlock()

	for _, c := range conns {
		_ = c.conn.Close()
	}
}

func (s *WebSocketServer) writeHandshakeError(w http.ResponseWriter, env *actionResponseEnvelope) {
	status := http.StatusUnauthorized
	if env.Retcode == 1403 {
		status = http.StatusForbidden
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(env)
	if err != nil {
		// 如果编码失败，响应头已经发送，无法返回错误
		// 记录错误或忽略
		_ = err
	}
}

func mapHandlerError(err error) *entity.ActionRawResponse {
	switch {
	case errors.Is(err, ErrActionNotFound):
		return &entity.ActionRawResponse{
			Status:  entity.StatusFailed,
			Retcode: 1404,
			Message: err.Error(),
		}
	case errors.Is(err, ErrBadRequest):
		return &entity.ActionRawResponse{
			Status:  entity.StatusFailed,
			Retcode: 1400,
			Message: err.Error(),
		}
	default:
		return &entity.ActionRawResponse{
			Status:  entity.StatusFailed,
			Retcode: 1500,
			Message: err.Error(),
		}
	}
}
