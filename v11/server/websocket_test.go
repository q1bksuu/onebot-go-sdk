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
	"github.com/q1bksuu/onebot-go-sdk/v11/entity"
	"github.com/q1bksuu/onebot-go-sdk/v11/internal/util"
	"github.com/stretchr/testify/require"
)

// stubHandler is a simple ActionRequestHandler for testing.
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
	t *testing.T, cfg WSConfig, handler ActionRequestHandler,
) (*WebSocketServer, *httptest.Server) {
	t.Helper()

	wsServer := NewWebSocketServer(cfg, handler)
	testServer := httptest.NewServer(wsServer.srv.Handler)

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

		wsServer := &WebSocketServer{handler: &stubHandler{err: ErrActionNotFound}}
		resp := wsServer.handleActionMessage(ctx, []byte(`{"action":"x","params":{}}`))
		require.Equal(t, entity.ActionResponseRetcode(1404), resp.Retcode)
		require.Equal(t, entity.StatusFailed, resp.Status)
		require.Equal(t, ErrActionNotFound.Error(), resp.Message)
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

func TestWriteHandshakeError(t *testing.T) {
	t.Parallel()

	wsServer := &WebSocketServer{}

	t.Run("401", func(t *testing.T) {
		t.Parallel()

		rr := httptest.NewRecorder()
		wsServer.writeHandshakeError(rr, &actionResponseEnvelope{Retcode: 1401, Status: entity.StatusFailed})
		require.Equal(t, http.StatusUnauthorized, rr.Code)
		require.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	})

	t.Run("403", func(t *testing.T) {
		t.Parallel()

		rr := httptest.NewRecorder()
		wsServer.writeHandshakeError(rr, &actionResponseEnvelope{Retcode: 1403, Status: entity.StatusFailed})
		require.Equal(t, http.StatusForbidden, rr.Code)
	})
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
		msg := readJSON[*actionResponseEnvelope](t, conn)
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
		msg := readJSON[*actionResponseEnvelope](t, conn)
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
