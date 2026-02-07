package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/q1bksuu/onebot-go-sdk/v11/entity"
	"github.com/q1bksuu/onebot-go-sdk/v11/server"
	"github.com/stretchr/testify/require"
)

type mockActionHandler struct {
	handleFn func(ctx context.Context, req *entity.ActionRequest) (*entity.ActionRawResponse, error)
}

func (m *mockActionHandler) HandleActionRequest(
	ctx context.Context, req *entity.ActionRequest,
) (*entity.ActionRawResponse, error) {
	return m.handleFn(ctx, req)
}

func newActionConnServer(
	t *testing.T,
	responseCh chan<- *entity.ActionResponseEnvelope,
	serverErrCh chan<- error,
	holdConnCh <-chan struct{},
) *httptest.Server {
	t.Helper()

	reportErr := func(err error) {
		if err == nil {
			return
		}

		select {
		case serverErrCh <- err:
		default:
		}
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(_ *http.Request) bool { return true },
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			reportErr(err)

			return
		}

		defer func() {
			closeErr := conn.Close()
			reportErr(closeErr)
		}()

		requestPayload := `{"action":"ping","params":{"foo":"bar"},"echo":"e1"}`

		err = conn.WriteMessage(websocket.TextMessage, []byte(requestPayload))
		if err != nil {
			reportErr(err)

			return
		}

		_, data, err := conn.ReadMessage()
		if err != nil {
			reportErr(err)

			return
		}

		var resp entity.ActionResponseEnvelope

		err = json.Unmarshal(data, &resp)
		if err != nil {
			reportErr(err)

			return
		}

		responseCh <- &resp

		<-holdConnCh
	}))
}

func newPingActionHandler() *mockActionHandler {
	return &mockActionHandler{
		handleFn: func(_ context.Context, req *entity.ActionRequest) (*entity.ActionRawResponse, error) {
			if req.Action != "ping" {
				return &entity.ActionRawResponse{
					Status:  entity.StatusFailed,
					Retcode: 1404,
					Message: "unexpected action",
				}, nil
			}

			return &entity.ActionRawResponse{
				Status:  entity.StatusOK,
				Retcode: 0,
				Data:    json.RawMessage(`{"result":"ok"}`),
			}, nil
		},
	}
}

func TestWebSocketClient_HandleActionMessage(t *testing.T) {
	t.Parallel()

	errMock := errors.New("mock error") //nolint:err113 // mock error for testing
	handler := &mockActionHandler{
		handleFn: func(_ context.Context, req *entity.ActionRequest) (*entity.ActionRawResponse, error) {
			if req.Action == "test" {
				return &entity.ActionRawResponse{
					Status:  entity.StatusOK,
					Retcode: 0,
					Data:    json.RawMessage(`{"result":"ok"}`),
				}, nil
			}

			return nil, errMock
		},
	}

	client := &WebSocketClient{
		actionHandler: handler,
	}

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		req := `{"action":"test","params":{"foo":"bar"},"echo":"123"}`
		resp := client.handleActionMessage(ctx, []byte(req))

		require.Equal(t, entity.StatusOK, resp.Status)
		require.Equal(t, entity.ActionResponseRetcode(0), resp.Retcode)
		require.JSONEq(t, `{"result":"ok"}`, string(resp.Data))
		require.JSONEq(t, `"123"`, string(resp.Echo))
	})

	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()

		req := `{"action":`
		resp := client.handleActionMessage(ctx, []byte(req))

		require.Equal(t, entity.StatusFailed, resp.Status)
		require.Equal(t, entity.ActionResponseRetcode(1400), resp.Retcode)
	})
}

func TestNewWebSocketClient_WithWSConfig(t *testing.T) {
	t.Parallel()

	cfg := WSClientConfig{
		URL:               "ws://example",
		ReconnectInterval: 2 * time.Second,
		SelfID:            123,
		AccessToken:       "token",
		ReadTimeout:       3 * time.Second,
		WriteTimeout:      4 * time.Second,
	}

	client := NewWebSocketClient(WithWSConfig(cfg))

	require.Equal(t, cfg, client.cfg)
}

func TestNewWebSocketClient_OptionsOverride(t *testing.T) {
	t.Parallel()

	cfg := WSClientConfig{
		URL:               "ws://example",
		ReconnectInterval: 2 * time.Second,
		SelfID:            123,
		AccessToken:       "token",
		ReadTimeout:       3 * time.Second,
		WriteTimeout:      4 * time.Second,
	}

	client := NewWebSocketClient(
		WithWSConfig(cfg),
		WithWSURL("ws://override"),
		WithWSSelfID(456),
		WithWSAccessToken("override-token"),
	)

	require.Equal(t, "ws://override", client.cfg.URL)
	require.Equal(t, int64(456), client.cfg.SelfID)
	require.Equal(t, "override-token", client.cfg.AccessToken)
	require.Equal(t, cfg.ReconnectInterval, client.cfg.ReconnectInterval)
	require.Equal(t, cfg.ReadTimeout, client.cfg.ReadTimeout)
	require.Equal(t, cfg.WriteTimeout, client.cfg.WriteTimeout)
}

func TestNewWebSocketClient_WithWSActionHandler(t *testing.T) {
	t.Parallel()

	handler := &mockActionHandler{
		handleFn: func(
			ctx context.Context,
			req *entity.ActionRequest,
		) (*entity.ActionRawResponse, error) {
			_ = ctx
			_ = req

			return &entity.ActionRawResponse{Status: entity.StatusOK}, nil
		},
	}

	client := NewWebSocketClient(WithWSActionHandler(handler))

	require.Same(t, handler, client.actionHandler)
}

func TestWebSocketClient_Start_URLEmpty(t *testing.T) {
	t.Parallel()

	client := NewWebSocketClient()

	err := client.Start(context.Background())

	require.ErrorIs(t, err, server.ErrUniversalClientURLEmpty)
}

func TestWebSocketClient_AdditionalOptions(t *testing.T) {
	t.Parallel()

	client := NewWebSocketClient(
		WithWSReconnectInterval(2*time.Second),
		WithWSReadTimeout(3*time.Second),
		WithWSWriteTimeout(4*time.Second),
	)

	require.Equal(t, 2*time.Second, client.cfg.ReconnectInterval)
	require.Equal(t, 3*time.Second, client.cfg.ReadTimeout)
	require.Equal(t, 4*time.Second, client.cfg.WriteTimeout)
}

func TestWebSocketClient_BuildHeaders(t *testing.T) {
	t.Parallel()

	client := NewWebSocketClient(
		WithWSSelfID(12345),
		WithWSAccessToken("token"),
	)

	headers := client.buildHeaders("Universal")
	require.Equal(t, "12345", headers.Get("X-Self-Id"))
	require.Equal(t, "Universal", headers.Get("X-Client-Role"))
	require.Equal(t, "Bearer token", headers.Get("Authorization"))
}

func TestWebSocketClient_DialWithReconnect_ContextCanceled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := NewWebSocketClient(WithWSReconnectInterval(10 * time.Millisecond))

	_, err := client.dialWithReconnect(ctx, "ws://example.com/ws", http.Header{})
	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)
	require.ErrorContains(t, err, "dial context canceled")
}

func TestWebSocketClient_Shutdown_NoConn(t *testing.T) {
	t.Parallel()

	client := NewWebSocketClient()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := client.Shutdown(ctx)
	require.NoError(t, err)
}

func TestWebSocketClient_BroadcastEvent_NoConn(t *testing.T) {
	t.Parallel()

	client := NewWebSocketClient()

	err := client.BroadcastEvent(&entity.PrivateMessageEvent{})
	require.NoError(t, err)
}

func TestWebSocketClient_ClearConn(t *testing.T) {
	t.Parallel()

	client := NewWebSocketClient()
	conn := &websocket.Conn{}
	client.conn = conn

	client.clearConn(conn)
	require.Nil(t, client.conn)
}

func TestWebSocketClient_Start_RunServeActionConn(t *testing.T) {
	t.Parallel()

	responseCh := make(chan *entity.ActionResponseEnvelope, 1)
	serverErrCh := make(chan error, 1)
	holdConnCh := make(chan struct{})

	serverForTest := newActionConnServer(t, responseCh, serverErrCh, holdConnCh)
	t.Cleanup(serverForTest.Close)
	t.Cleanup(func() { close(holdConnCh) })

	wsURL := "ws" + strings.TrimPrefix(serverForTest.URL, "http")

	handler := newPingActionHandler()

	client := NewWebSocketClient(
		WithWSURL(wsURL),
		WithWSActionHandler(handler),
		WithWSReadTimeout(200*time.Millisecond),
		WithWSWriteTimeout(200*time.Millisecond),
		WithWSReconnectInterval(10*time.Millisecond),
	)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	startErrCh := make(chan error, 1)

	go func() {
		startErrCh <- client.Start(ctx)
	}()

	select {
	case err := <-serverErrCh:
		require.NoError(t, err)
	case resp := <-responseCh:
		require.Equal(t, entity.StatusOK, resp.Status)
		require.Equal(t, entity.ActionResponseRetcode(0), resp.Retcode)
		require.JSONEq(t, `{"result":"ok"}`, string(resp.Data))
		require.JSONEq(t, `"e1"`, string(resp.Echo))
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for websocket response")
	}

	cancel()

	require.NoError(t, <-startErrCh)

	select {
	case err := <-serverErrCh:
		require.NoError(t, err)
	default:
	}
}
