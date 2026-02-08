package server

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
	"github.com/q1bksuu/onebot-go-sdk/v11/dispatcher"
	"github.com/q1bksuu/onebot-go-sdk/v11/entity"
	"github.com/q1bksuu/onebot-go-sdk/v11/internal/util"
	"github.com/stretchr/testify/require"
)

// stubHandler is a simple dispatcher.ActionRequestHandler for testing.
type stubHandler struct {
	resp *entity.ActionRawResponse
	err  error
}

func (h *stubHandler) HandleActionRequest(
	_ context.Context, _ *entity.ActionRequest,
) (*entity.ActionRawResponse, error) {
	return h.resp, h.err
}

// helper to quickly build WS server backed by WebSocketServer mux.
func newTestWebSocketServer(
	t *testing.T, cfg WSConfig, handler dispatcher.ActionRequestHandler,
) (*WebSocketServer, *httptest.Server) {
	t.Helper()

	wsServer := NewWebSocketServer(
		WithWSConfig(cfg),
		WithWSActionHandler(handler),
	)
	testServer := httptest.NewServer(wsServer.Srv.Handler)

	return wsServer, testServer
}

func wsURL(ts *httptest.Server, path string) string {
	return strings.Replace(ts.URL, "http", "ws", 1) + path
}

// helper to dial websocket with headers.
func mustDialWS(t *testing.T, rawURL string, header http.Header) *websocket.Conn {
	t.Helper()

	d := websocket.Dialer{}
	conn, resp, err := d.Dial(rawURL, header)
	require.NoError(t, err)

	if resp != nil {
		_ = resp.Body.Close()
	}

	return conn
}

func readJSON[T any](t *testing.T, conn *websocket.Conn) T {
	t.Helper()

	_, data, err := conn.ReadMessage()
	require.NoError(t, err)

	var rtn T
	require.NoError(t, json.Unmarshal(data, &rtn))

	return rtn
}

func TestNormalizeAndMatchPath(t *testing.T) {
	t.Parallel()

	t.Run("normalizePathPrefix", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			in   string
			want string
		}{
			{"", ""},
			{"/", ""},
			{"api", "/api"},
			{"//api//v1", "/api//v1"},
		}
		for _, tc := range cases {
			t.Run(tc.in, func(t *testing.T) {
				t.Parallel()

				got := util.NormalizePath(tc.in)
				require.Equal(t, tc.want, got)
			})
		}
	})

	t.Run("matchPath", func(t *testing.T) {
		t.Parallel()

		s := &WebSocketServer{cfg: WSConfig{PathPrefix: "/api/v1"}}
		require.True(t, s.matchPath("/api/v1/api", "/api"))
		require.True(t, s.matchPath("/api/v1/api/", "/api"))
		require.False(t, s.matchPath("/other/api", "/api"))
	})

	t.Run("matchUniversalPath", func(t *testing.T) {
		t.Parallel()

		s1 := &WebSocketServer{cfg: WSConfig{}}
		require.True(t, s1.matchUniversalPath("/"))
		require.False(t, s1.matchUniversalPath("/x"))

		s2 := &WebSocketServer{cfg: WSConfig{PathPrefix: "/api"}}
		require.True(t, s2.matchUniversalPath("/api"))
		require.True(t, s2.matchUniversalPath("/api/"))
		require.False(t, s2.matchUniversalPath("/"))
	})
}

func TestCheckAccess(t *testing.T) {
	t.Parallel()

	wsServer := &WebSocketServer{cfg: WSConfig{AccessToken: "secret"}}

	t.Run("missing", func(t *testing.T) {
		t.Parallel()

		r := httptest.NewRequest(http.MethodGet, "https://example.com/api", nil)
		resp := wsServer.checkAccess(r)
		require.NotNil(t, resp)
		require.Equal(t, entity.StatusFailed, resp.Status)
		require.Equal(t, entity.ActionResponseRetcode(1401), resp.Retcode)
	})

	t.Run("wrong", func(t *testing.T) {
		t.Parallel()

		r := httptest.NewRequest(http.MethodGet, "https://example.com/api?access_token=bad", nil)
		resp := wsServer.checkAccess(r)
		require.NotNil(t, resp)
		require.Equal(t, entity.ActionResponseRetcode(1403), resp.Retcode)
	})

	t.Run("header", func(t *testing.T) {
		t.Parallel()

		r := httptest.NewRequest(http.MethodGet, "https://example.com/api", nil)
		r.Header.Set("Authorization", "Bearer secret")
		resp := wsServer.checkAccess(r)
		require.Nil(t, resp)
	})

	t.Run("query", func(t *testing.T) {
		t.Parallel()

		r := httptest.NewRequest(http.MethodGet, "https://example.com/api?access_token=secret", nil)
		resp := wsServer.checkAccess(r)
		require.Nil(t, resp)
	})

	t.Run("no_access_token_configured", func(t *testing.T) {
		t.Parallel()

		noAuthServer := &WebSocketServer{cfg: WSConfig{}}
		r := httptest.NewRequest(http.MethodGet, "https://example.com/api", nil)
		resp := noAuthServer.checkAccess(r)
		require.Nil(t, resp)
	})

	t.Run("header_not_bearer_does_not_fallback_to_query", func(t *testing.T) {
		t.Parallel()

		r := httptest.NewRequest(http.MethodGet, "https://example.com/api?access_token=secret", nil)
		r.Header.Set("Authorization", "Token secret")
		resp := wsServer.checkAccess(r)
		require.NotNil(t, resp)
		require.Equal(t, entity.ActionResponseRetcode(1403), resp.Retcode)
	})
}

func TestHandleActionMessageMapping(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()

		wsServer := &WebSocketServer{handler: &stubHandler{}}
		resp := wsServer.handleActionMessage(ctx, []byte("{"))
		require.Equal(t, entity.StatusFailed, resp.Status)
		require.Equal(t, entity.ActionResponseRetcode(1400), resp.Retcode)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		wsServer := &WebSocketServer{handler: &stubHandler{err: dispatcher.ErrActionNotFound}}
		resp := wsServer.handleActionMessage(ctx, []byte(`{"action":"x","params":{}}`))
		require.Equal(t, entity.ActionResponseRetcode(1404), resp.Retcode)
		require.Equal(t, entity.StatusFailed, resp.Status)
		require.Equal(t, dispatcher.ErrActionNotFound.Error(), resp.Message)
	})

	t.Run("bad request", func(t *testing.T) {
		t.Parallel()

		wsServer := &WebSocketServer{handler: &stubHandler{err: ErrBadRequest}}
		resp := wsServer.handleActionMessage(ctx, []byte(`{"action":"x","params":{}}`))
		require.Equal(t, entity.ActionResponseRetcode(1400), resp.Retcode)
		require.Equal(t, entity.StatusFailed, resp.Status)
		require.Equal(t, ErrBadRequest.Error(), resp.Message)
	})

	t.Run("internal", func(t *testing.T) {
		t.Parallel()
		// 测试内部错误，使用动态错误是合理的
		//nolint:err113 // 测试代码中使用动态错误是合理的
		boomErr := errors.New("boom")
		wsServer := &WebSocketServer{handler: &stubHandler{err: boomErr}}
		resp := wsServer.handleActionMessage(ctx, []byte(`{"action":"x","params":{}}`))
		require.Equal(t, entity.ActionResponseRetcode(1500), resp.Retcode)
		require.Equal(t, entity.StatusFailed, resp.Status)
		require.Equal(t, "boom", resp.Message)
	})

	t.Run("success with echo", func(t *testing.T) {
		t.Parallel()

		wsServer := &WebSocketServer{
			handler: &stubHandler{
				resp: &entity.ActionRawResponse{
					Status:  entity.StatusOK,
					Retcode: 0,
					Data:    json.RawMessage(`{"ok":true}`),
					Message: "ok",
				},
			},
		}
		resp := wsServer.handleActionMessage(ctx, []byte(`{"action":"x","params":{},"echo":"e1"}`))
		require.Equal(t, entity.StatusOK, resp.Status)
		require.Equal(t, entity.RetcodeSuccess, resp.Retcode)
		require.JSONEq(t, `"e1"`, string(resp.Echo))
		require.Equal(t, "ok", resp.Message)
	})

	t.Run("nil response default failed", func(t *testing.T) {
		t.Parallel()

		wsServer := &WebSocketServer{handler: &stubHandler{resp: nil}}
		resp := wsServer.handleActionMessage(ctx, []byte(`{"action":"x","params":{}}`))
		require.Equal(t, entity.StatusFailed, resp.Status)
		require.Equal(t, entity.ActionResponseRetcode(-1), resp.Retcode)
		require.Equal(t, "empty response", resp.Message)
	})
}

func TestNewWebSocketServer_OptionsApplied(t *testing.T) {
	t.Parallel()

	checkOrigin := func(_ *http.Request) bool { return false }

	server := NewWebSocketServer(
		WithWSAddr(":6700"),
		WithWSPathPrefix("/api"),
		WithWSAccessToken("token"),
		WithWSCheckOrigin(checkOrigin),
		WithWSReadTimeout(1*time.Second),
		WithWSWriteTimeout(2*time.Second),
		WithWSIdleTimeout(3*time.Second),
	)

	require.Equal(t, ":6700", server.cfg.Addr)
	require.Equal(t, "/api", server.cfg.PathPrefix)
	require.Equal(t, "token", server.cfg.AccessToken)
	require.Equal(t, 1*time.Second, server.cfg.ReadTimeout)
	require.Equal(t, 2*time.Second, server.cfg.WriteTimeout)
	require.Equal(t, 3*time.Second, server.cfg.IdleTimeout)
	require.False(t, server.upgrader.CheckOrigin(httptest.NewRequest(http.MethodGet, "http://example.com", nil)))
}

func TestWebSocketServer_Handler_NotFoundRoutes(t *testing.T) {
	t.Parallel()

	t.Run("api_and_event_mismatch", func(t *testing.T) {
		t.Parallel()

		server := NewWebSocketServer(WithWSPathPrefix("/ws"))
		handler := server.Handler()
		require.Same(t, server.Srv.Handler, handler)

		cases := []string{
			"/ws/api/extra",
			"/ws/event/extra",
		}
		for _, path := range cases {
			t.Run(path, func(t *testing.T) {
				t.Parallel()

				recorder := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "http://example.com"+path, nil)
				handler.ServeHTTP(recorder, req)
				require.Equal(t, http.StatusNotFound, recorder.Code)
			})
		}
	})

	t.Run("universal_mismatch", func(t *testing.T) {
		t.Parallel()

		server := NewWebSocketServer()
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "http://example.com/other", nil)
		server.Handler().ServeHTTP(recorder, req)
		require.Equal(t, http.StatusNotFound, recorder.Code)
	})
}

func TestWebSocketServer_Start_ShutdownOnContextCancel(t *testing.T) {
	t.Parallel()

	server := NewWebSocketServer(WithWSAddr("127.0.0.1:0"))

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)

	go func() {
		errCh <- server.Start(ctx)
	}()

	time.Sleep(20 * time.Millisecond)
	cancel()

	select {
	case err := <-errCh:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for server to stop")
	}
}

func TestWriteHandshakeError(t *testing.T) {
	t.Parallel()

	wsServer := &WebSocketServer{}

	t.Run("401", func(t *testing.T) {
		t.Parallel()

		recorder := httptest.NewRecorder()
		wsServer.writeHandshakeError(recorder, &entity.ActionResponseEnvelope{
			ActionRawResponse: entity.ActionRawResponse{
				Retcode: 1401,
				Status:  entity.StatusFailed,
			},
		})
		require.Equal(t, http.StatusUnauthorized, recorder.Code)
		require.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	})

	t.Run("403", func(t *testing.T) {
		t.Parallel()

		recorder := httptest.NewRecorder()
		wsServer.writeHandshakeError(recorder, &entity.ActionResponseEnvelope{
			ActionRawResponse: entity.ActionRawResponse{
				Retcode: 1403,
				Status:  entity.StatusFailed,
			},
		})
		require.Equal(t, http.StatusForbidden, recorder.Code)
	})
}

func TestWebSocketServer_CloseAllConns(t *testing.T) {
	t.Parallel()

	cfg := WSConfig{
		ReadTimeout:  50 * time.Millisecond,
		WriteTimeout: 50 * time.Millisecond,
		IdleTimeout:  50 * time.Millisecond,
	}
	handler := &stubHandler{resp: &entity.ActionRawResponse{Status: entity.StatusOK, Retcode: 0}}

	wsServer, testServer := newTestWebSocketServer(t, cfg, handler)
	t.Cleanup(func() {
		testServer.Close()
	})

	eventConn := mustDialWS(t, wsURL(testServer, "/event"), nil)
	universalConn := mustDialWS(t, wsURL(testServer, "/"), nil)

	t.Cleanup(func() {
		_ = eventConn.Close()
		_ = universalConn.Close()
	})

	require.Eventually(t, func() bool {
		wsServer.mu.Lock()
		defer wsServer.mu.Unlock()

		return len(wsServer.eventConns) == 1 && len(wsServer.universalConn) == 1
	}, 500*time.Millisecond, 10*time.Millisecond)

	wsServer.closeAllConns()

	require.Eventually(t, func() bool {
		wsServer.mu.Lock()
		defer wsServer.mu.Unlock()

		return len(wsServer.eventConns) == 0 && len(wsServer.universalConn) == 0
	}, 500*time.Millisecond, 10*time.Millisecond)

	_ = eventConn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	_, _, err := eventConn.ReadMessage()
	require.Error(t, err)

	_ = universalConn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	_, _, err = universalConn.ReadMessage()
	require.Error(t, err)
}

func TestHandleAPIAndUniversalFlow(t *testing.T) {
	t.Parallel()

	cfg := WSConfig{
		Addr:         "",
		PathPrefix:   "",
		AccessToken:  "tok",
		ReadTimeout:  50 * time.Millisecond,
		WriteTimeout: 50 * time.Millisecond,
		IdleTimeout:  50 * time.Millisecond,
	}

	handler := &stubHandler{resp: &entity.ActionRawResponse{Status: entity.StatusOK, Retcode: 0, Message: "ok"}}

	_, testServer := newTestWebSocketServer(t, cfg, handler)

	t.Cleanup(func() {
		testServer.Close()
	})

	t.Run("unauthorized", func(t *testing.T) {
		t.Parallel()

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, testServer.URL+"/api", nil)
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		_ = resp.Body.Close()
	})

	t.Run("authorized api action", func(t *testing.T) {
		t.Parallel()

		header := http.Header{"Authorization": []string{"Bearer tok"}}

		conn := mustDialWS(t, wsURL(testServer, "/api"), header)

		defer func() {
			_ = conn.Close()
		}()

		_ = conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))

		require.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"action":"ping","params":{},"echo":"e"}`)))
		msg := readJSON[*entity.ActionResponseEnvelope](t, conn)
		require.Equal(t, entity.RetcodeSuccess, msg.Retcode)
		require.JSONEq(t, `"e"`, string(msg.Echo))
	})

	t.Run("authorized universal action", func(t *testing.T) {
		t.Parallel()

		header := http.Header{"Authorization": []string{"Bearer tok"}}

		conn := mustDialWS(t, wsURL(testServer, "/"), header)

		defer func() {
			_ = conn.Close()
		}()

		_ = conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))

		require.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"action":"ping","params":{},"echo":"u"}`)))
		msg := readJSON[*entity.ActionResponseEnvelope](t, conn)
		require.Equal(t, entity.RetcodeSuccess, msg.Retcode)
		require.JSONEq(t, `"u"`, string(msg.Echo))
	})
}

func TestBroadcastEventIntegration(t *testing.T) {
	t.Parallel()

	cfg := WSConfig{
		AccessToken:  "tok",
		ReadTimeout:  50 * time.Millisecond,
		WriteTimeout: 50 * time.Millisecond,
		IdleTimeout:  50 * time.Millisecond,
	}
	handler := &stubHandler{resp: &entity.ActionRawResponse{Status: entity.StatusOK, Retcode: 0}}

	wsServer, testServer := newTestWebSocketServer(t, cfg, handler)
	defer testServer.Close()

	header := http.Header{"Authorization": []string{"Bearer tok"}}

	eventConn := mustDialWS(t, wsURL(testServer, "/event"), header)

	defer func() {
		_ = eventConn.Close()
	}()

	_ = eventConn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))

	universalConn := mustDialWS(t, wsURL(testServer, "/"), header)

	defer func() {
		_ = universalConn.Close()
	}()

	_ = universalConn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))

	testEvent := &entity.PrivateMessageEvent{
		Time:        time.Now().Unix(),
		SelfId:      1,
		PostType:    entity.EventPostTypeMessage,
		MessageType: entity.EventMessageTypePrivate,
		SubType:     entity.EventPrivateMessageSubTypeFriend,
		MessageId:   2,
		UserId:      3,
		Message:     nil,
		RawMessage:  "",
		Font:        0,
		Sender:      nil,
	}
	// 广播事件
	wsServer.BroadcastEvent(testEvent)

	msg1 := readJSON[entity.PrivateMessageEvent](t, eventConn)
	msg2 := readJSON[entity.PrivateMessageEvent](t, universalConn)

	require.Equal(t, *testEvent, msg1)
	require.Equal(t, *testEvent, msg2)
}
