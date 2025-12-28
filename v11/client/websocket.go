package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/q1bksuu/onebot-go-sdk/v11/entity"
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
	actionHandler server.ActionRequestHandler

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
	// 合并外部 context 和内部 context
	mergedCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		<-mergedCtx.Done()
		c.cancel()
	}()

	url := c.getURL()
	if url == "" {
		return server.ErrUniversalClientURLEmpty
	}

	c.wg.Add(1)

	go c.run(mergedCtx, url)

	// 等待 context 取消
	<-mergedCtx.Done()

	// 等待所有 goroutine 完成
	c.wg.Wait()

	return nil
}

// Shutdown 优雅关闭所有连接.
func (c *WebSocketClient) Shutdown(ctx context.Context) error {
	c.cancel()

	// 关闭所有连接
	c.mu.Lock()

	conns := []*websocket.Conn{}
	if c.conn != nil {
		conns = append(conns, c.conn)
	}

	c.mu.Unlock()

	for _, conn := range conns {
		_ = conn.Close()
	}

	// 等待所有 goroutine 完成
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

// getURL 获取 Universal 客户端 URL.
func (c *WebSocketClient) getURL() string {
	return c.cfg.URL
}

// buildHeaders 构建连接请求头.
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
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("dial context canceled: %w", ctx.Err())
		default:
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

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		headers := c.buildHeaders("Universal")

		conn, err := c.dialWithReconnect(ctx, url, headers)
		if err != nil {
			return
		}

		c.mu.Lock()
		c.conn = conn
		c.mu.Unlock()

		// 在独立 goroutine 中处理 API 请求
		apiCtx, apiCancel := context.WithCancel(ctx)
		apiDone := make(chan struct{})

		go func() {
			defer close(apiDone)

			c.serveActionConn(apiCtx, conn)
		}()

		// 等待连接关闭或 context 取消
		select {
		case <-ctx.Done():
			apiCancel()
			<-apiDone
			c.mu.Lock()

			if c.conn == conn {
				c.conn = nil
			}

			c.mu.Unlock()

			_ = conn.Close()

			return
		case <-apiDone:
			// 连接断开，准备重连
			apiCancel()
		}

		c.mu.Lock()

		if c.conn == conn {
			c.conn = nil
		}

		c.mu.Unlock()

		_ = conn.Close()
	}
}

// serveActionConn 处理 API 连接的消息.
func (c *WebSocketClient) serveActionConn(ctx context.Context, conn *websocket.Conn) {
	if c.actionHandler == nil {
		c.serveActionConnWithoutHandler(ctx, conn)

		return
	}

	c.serveActionConnWithHandler(ctx, conn)
}

// serveActionConnWithoutHandler 处理没有 handler 的情况，只读取消息但不处理.
func (c *WebSocketClient) serveActionConnWithoutHandler(ctx context.Context, conn *websocket.Conn) {
	for {
		if ctx.Err() != nil {
			return
		}

		_, _, err := conn.ReadMessage()
		if err != nil {
			return
		}
	}
}

// serveActionConnWithHandler 处理有 handler 的情况.
func (c *WebSocketClient) serveActionConnWithHandler(ctx context.Context, conn *websocket.Conn) {
	for {
		if ctx.Err() != nil {
			return
		}

		if !c.setReadDeadline(conn) {
			return
		}

		data, err := c.readMessage(conn)
		if err != nil {
			return
		}

		resp := c.handleActionMessage(ctx, data)

		if !c.setWriteDeadline(conn) {
			return
		}

		err = conn.WriteJSON(resp)
		if err != nil {
			return
		}
	}
}

// setReadDeadline 设置读取超时.
func (c *WebSocketClient) setReadDeadline(conn *websocket.Conn) bool {
	if c.cfg.ReadTimeout > 0 {
		err := conn.SetReadDeadline(time.Now().Add(c.cfg.ReadTimeout))
		if err != nil {
			return false
		}
	}

	return true
}

// setWriteDeadline 设置写入超时.
func (c *WebSocketClient) setWriteDeadline(conn *websocket.Conn) bool {
	if c.cfg.WriteTimeout > 0 {
		err := conn.SetWriteDeadline(time.Now().Add(c.cfg.WriteTimeout))
		if err != nil {
			return false
		}
	}

	return true
}

// readMessage 读取消息.
func (c *WebSocketClient) readMessage(conn *websocket.Conn) ([]byte, error) {
	_, data, err := conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("read message failed: %w", err)
	}

	return data, nil
}

// handleActionMessage 处理动作请求消息.
func (c *WebSocketClient) handleActionMessage(_ context.Context, _ []byte) any {
	panic("unimplemented")
}
