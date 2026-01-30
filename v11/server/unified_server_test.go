package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/q1bksuu/onebot-go-sdk/v11/dispatcher"
	"github.com/q1bksuu/onebot-go-sdk/v11/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helper to quickly build UnifiedServer backed by httptest.Server.
func newTestUnifiedServer(
	t *testing.T,
	cfg UnifiedConfig,
	httpHandler dispatcher.ActionRequestHandler,
	wsHandler dispatcher.ActionRequestHandler,
) *httptest.Server {
	t.Helper()

	cfg.HTTP.ActionHandler = httpHandler
	cfg.WS.ActionHandler = wsHandler

	unifiedServer := NewUnifiedServer(cfg)
	testServer := httptest.NewServer(unifiedServer.Srv.Handler)

	return testServer
}

func TestUnifiedServer_Initialization(t *testing.T) {
	t.Parallel()

	cfg := UnifiedConfig{
		ServerConfig: ServerConfig{
			Addr: ":1234",
		},
		HTTP: UnifiedHTTPConfig{
			APIPathPrefix: "/api",
		},
		WS: UnifiedWSConfig{
			PathPrefix: "/ws",
		},
	}

	server := NewUnifiedServer(cfg)

	// 验证配置覆盖
	assert.Equal(t, cfg.Addr, server.httpSrv.cfg.Addr)
	assert.Equal(t, cfg.Addr, server.wsSrv.cfg.Addr)

	// 验证前缀归一化 (NewHTTPServer 和 NewWebSocketServer 会做这个)
	assert.Equal(t, "/api/", server.httpSrv.cfg.APIPathPrefix)
	// WS PathPrefix 归一化逻辑在 matchPath 中，这里直接看 cfg
	assert.Equal(t, "/ws", server.wsSrv.cfg.PathPrefix)
}

func TestUnifiedServer_Routing_HTTP(t *testing.T) {
	t.Parallel()

	cfg := UnifiedConfig{
		HTTP: UnifiedHTTPConfig{
			APIPathPrefix: "/api",
		},
	}

	httpHandler := dispatcher.ActionRequestHandlerFunc(
		func(_ context.Context, _ *entity.ActionRequest) (*entity.ActionRawResponse, error) {
			return &entity.ActionRawResponse{
				Status:  entity.StatusOK,
				Retcode: 0,
				Message: "http response",
			}, nil
		})

	testServer := newTestUnifiedServer(t, cfg, httpHandler, nil)
	defer testServer.Close()

	// 发送普通 HTTP 请求
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, testServer.URL+"/api/test_action", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer func() {
		_ = resp.Body.Close()
	}()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var apiResp entity.ActionRawResponse

	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	require.NoError(t, err)
	assert.Equal(t, "http response", apiResp.Message)
}

func TestUnifiedServer_Routing_WebSocket(t *testing.T) {
	t.Parallel()

	cfg := UnifiedConfig{
		WS: UnifiedWSConfig{
			PathPrefix: "/ws",
		},
	}

	wsHandler := dispatcher.ActionRequestHandlerFunc(
		func(_ context.Context, _ *entity.ActionRequest) (*entity.ActionRawResponse, error) {
			return &entity.ActionRawResponse{
				Status:  entity.StatusOK,
				Retcode: 0,
				Message: "ws response",
			}, nil
		})

	testServer := newTestUnifiedServer(t, cfg, nil, wsHandler)
	defer testServer.Close()

	// 连接 WebSocket
	wsURL := strings.Replace(testServer.URL, "http", "ws", 1) + "/ws"
	conn, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)

	if resp != nil {
		_ = resp.Body.Close()
	}

	defer func() {
		_ = conn.Close()
	}()

	// 发送 Action 请求
	reqJSON := `{"action": "test_action", "params": {}, "echo": "123"}`
	err = conn.WriteMessage(websocket.TextMessage, []byte(reqJSON))
	require.NoError(t, err)

	// 读取响应
	_, msg, err := conn.ReadMessage()
	require.NoError(t, err)

	var wsResp entity.ActionResponseEnvelope

	err = json.Unmarshal(msg, &wsResp)
	require.NoError(t, err)
	assert.Equal(t, "ws response", wsResp.Message)

	// Echo 应该是一个 json.RawMessage (string/bytes)
	// 在 http_client_test.go 或 websocket_test.go 里是 require.JSONEq(t, `"e1"`, string(resp.Echo))
	// 这里 Echo "123" 应该被原样返回
	require.JSONEq(t, `"123"`, string(wsResp.Echo))
}

func TestUnifiedServer_Routing_UniversalPath_DistinguishByProtocol(t *testing.T) {
	t.Parallel()

	// HTTP 和 WS 都使用根路径 /
	cfg := UnifiedConfig{
		HTTP: UnifiedHTTPConfig{
			APIPathPrefix: "/",
		},
		WS: UnifiedWSConfig{
			PathPrefix: "/",
		},
	}

	httpHandler := dispatcher.ActionRequestHandlerFunc(
		func(_ context.Context, _ *entity.ActionRequest) (*entity.ActionRawResponse, error) {
			return &entity.ActionRawResponse{Message: "http"}, nil
		})

	wsHandler := dispatcher.ActionRequestHandlerFunc(
		func(_ context.Context, _ *entity.ActionRequest) (*entity.ActionRawResponse, error) {
			return &entity.ActionRawResponse{Message: "ws"}, nil
		})

	testServer := newTestUnifiedServer(t, cfg, httpHandler, wsHandler)
	defer testServer.Close()

	// 1. 普通 HTTP GET -> HTTP Handler
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, testServer.URL+"/test_action", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer func() {
		_ = resp.Body.Close()
	}()

	var apiResp entity.ActionRawResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&apiResp))
	assert.Equal(t, "http", apiResp.Message)

	// 2. WebSocket Upgrade -> WS Handler
	wsURL := strings.Replace(testServer.URL, "http", "ws", 1) + "/"
	conn, wsRespObj, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)

	if wsRespObj != nil {
		_ = wsRespObj.Body.Close()
	}

	defer func() {
		_ = conn.Close()
	}()

	require.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"action": "test", "echo": "echo"}`)))
	_, msg, err := conn.ReadMessage()
	require.NoError(t, err)

	var wsResp entity.ActionResponseEnvelope
	require.NoError(t, json.Unmarshal(msg, &wsResp))
	assert.Equal(t, "ws", wsResp.Message)
}

func TestUnifiedServer_StartAndShutdown(t *testing.T) {
	t.Parallel()

	// 使用随机端口
	cfg := UnifiedConfig{
		ServerConfig: ServerConfig{
			Addr: "127.0.0.1:0",
		},
	}
	server := NewUnifiedServer(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		errCh <- server.Start(ctx)
	}()

	// 稍微等待启动
	time.Sleep(100 * time.Millisecond)

	// 关闭
	cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Start returned unexpected error: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Start did not return after context cancel")
	}
}
