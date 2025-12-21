package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/q1bksuu/onebot-go-sdk/v11/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errUnexpected = errors.New("unexpected")
	errOther      = errors.New("other")
)

// 基础行为.
func TestNewHTTPServer_PathPrefixNormalizeAndHandler(t *testing.T) {
	t.Parallel()

	cfg := HTTPConfig{
		Addr:       ":0",
		PathPrefix: "api", // 不带斜杠
	}
	server := NewHTTPServer(cfg, ActionRequestHandlerFunc(
		func(_ context.Context, req *entity.ActionRequest) (*entity.ActionRawResponse, error) {
			// 简单回显 action，便于断言
			require.Equal(t, "test_action", req.Action)

			return &entity.ActionRawResponse{
				Status:  entity.StatusOK,
				Retcode: 0,
				Message: "ok",
			}, nil
		},
	))

	// PathPrefix 应被归一化为 "/api/"
	assert.Equal(t, "/api/", server.cfg.PathPrefix)

	// Handler() 应该是可用的 http.Handler
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/test_action", nil)
	server.Handler().ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)

	var resp entity.ActionRawResponse

	err := json.NewDecoder(recorder.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, entity.StatusOK, resp.Status)
	assert.Equal(t, entity.ActionResponseRetcode(0), resp.Retcode)
}

func newTestServer(cfg HTTPConfig, handler ActionRequestHandlerFunc) *HTTPServer {
	if cfg.PathPrefix == "" {
		cfg.PathPrefix = "/"
	}

	return NewHTTPServer(cfg, handler)
}

func TestHTTPServer_HandleRoot_PathAndNotFound(t *testing.T) {
	t.Parallel()

	server := newTestServer(HTTPConfig{PathPrefix: "/onebot"},
		func(_ context.Context, _ *entity.ActionRequest) (*entity.ActionRawResponse, error) {
			return &entity.ActionRawResponse{Status: entity.StatusOK, Retcode: 0}, nil
		})

	// 不匹配前缀 -> 404
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/other/action", nil)
	server.Handler().ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusNotFound, recorder.Code)

	// 只有前缀，没有 action -> 404
	recorder = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/onebot/", nil)
	server.Handler().ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestHTTPServer_AuthRequired_MissingOrWrongToken(t *testing.T) {
	t.Parallel()

	cfg := HTTPConfig{
		PathPrefix:  "/onebot",
		AccessToken: "secret",
	}
	server := newTestServer(cfg,
		func(_ context.Context, _ *entity.ActionRequest) (*entity.ActionRawResponse, error) {
			require.Fail(t, "handler should not be called when auth fails")

			return nil, errUnexpected
		},
	)

	// 无 token
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/onebot/test", nil)
	server.Handler().ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

	// 错误 token
	recorder = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/onebot/test?access_token=wrong", nil)
	server.Handler().ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestHTTPServer_AuthRequired_WithHeaderAndQuery(t *testing.T) {
	t.Parallel()

	cfg := HTTPConfig{
		PathPrefix:  "/onebot",
		AccessToken: "secret",
	}

	called := false
	server := newTestServer(cfg, func(_ context.Context, _ *entity.ActionRequest) (*entity.ActionRawResponse, error) {
		called = true

		return &entity.ActionRawResponse{Status: entity.StatusOK, Retcode: 0}, nil
	})

	// 使用 Authorization: Bearer
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/onebot/test", nil)
	req.Header.Set("Authorization", "Bearer secret")
	server.Handler().ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.True(t, called, "handler should be called when auth success via header")

	// 使用 query access_token
	called = false
	recorder = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/onebot/test?access_token=secret", nil)
	server.Handler().ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.True(t, called, "handler should be called when auth success via query")
}

// 参数解析 & JSON.
func TestHTTPServer_Params_QueryAndFormAndJSON(t *testing.T) {
	t.Parallel()

	type captured struct {
		req *entity.ActionRequest
	}

	var capturedReq captured

	handler := func(_ context.Context, req *entity.ActionRequest) (*entity.ActionRawResponse, error) {
		capturedReq.req = req

		return &entity.ActionRawResponse{Status: entity.StatusOK, Retcode: 0}, nil
	}
	server := newTestServer(HTTPConfig{PathPrefix: "/onebot"}, handler)

	body := map[string]any{
		"b": "override",
		"c": 3,
	}
	b, err := json.Marshal(body)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/onebot/do_something?a=1&arr=1&arr=2&b=2", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	server.Handler().ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)
	require.NotNil(t, capturedReq.req, "handler was not called")
	assert.Equal(t, "do_something", capturedReq.req.Action)

	params := capturedReq.req.Params
	assert.Equal(t, "1", params["a"])
	gotArr, ok := params["arr"].([]string)
	require.True(t, ok, "arr should be []string")
	require.Len(t, gotArr, 2)
	assert.Equal(t, []string{"1", "2"}, gotArr)

	// JSON 应覆盖同名字段 b，并追加 c
	assert.Equal(t, "override", params["b"])
	_, ok = params["c"]
	assert.True(t, ok, "expected param c from json")
}

func TestHTTPServer_Params_InvalidForm(t *testing.T) {
	t.Parallel()

	handler := func(_ context.Context, _ *entity.ActionRequest) (*entity.ActionRawResponse, error) {
		require.Fail(t, "handler should not be called on invalid form")

		return nil, errUnexpected
	}
	server := newTestServer(HTTPConfig{PathPrefix: "/onebot"}, handler)

	recorder := httptest.NewRecorder()
	// Content-Type 设置为 application/x-www-form-urlencoded，body 随意
	req := httptest.NewRequest(http.MethodPost, "/onebot/test", bytes.NewBufferString("%"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	server.Handler().ServeHTTP(recorder, req)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestHTTPServer_JSON_InvalidOrUnsupportedContentType(t *testing.T) {
	t.Parallel()

	handler := func(_ context.Context, _ *entity.ActionRequest) (*entity.ActionRawResponse, error) {
		require.Fail(t, "handler should not be called when json invalid or content-type unsupported")

		return nil, errUnexpected
	}
	server := newTestServer(HTTPConfig{PathPrefix: "/onebot"}, handler)

	// 无效 JSON
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/onebot/test", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	server.Handler().ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	// 不支持的 Content-Type
	recorder = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/onebot/test", bytes.NewBufferString("hello"))
	req.Header.Set("Content-Type", "text/plain")
	server.Handler().ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

// 错误映射 & 默认响应.
func TestHTTPServer_WriteError_Mapping(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{"not found", ErrActionNotFound, http.StatusNotFound},
		{"bad request", ErrBadRequest, http.StatusBadRequest},
		{"internal", errOther, http.StatusInternalServerError},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			handler := func(_ context.Context, _ *entity.ActionRequest) (*entity.ActionRawResponse, error) {
				return nil, testCase.err
			}
			server := newTestServer(HTTPConfig{PathPrefix: "/onebot"}, handler)

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/onebot/test", nil)
			server.Handler().ServeHTTP(recorder, req)

			assert.Equal(t, testCase.wantStatus, recorder.Code)
		})
	}
}

func TestHTTPServer_NilResponse_DefaultFailed(t *testing.T) {
	t.Parallel()

	handler := func(_ context.Context, _ *entity.ActionRequest) (*entity.ActionRawResponse, error) {
		//nolint:nilnil
		return nil, nil
	}
	server := newTestServer(HTTPConfig{PathPrefix: "/onebot"}, handler)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/onebot/test", nil)
	server.Handler().ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)

	var resp entity.ActionRawResponse

	err := json.NewDecoder(recorder.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, entity.StatusFailed, resp.Status)
	assert.Equal(t, int64(-1), int64(resp.Retcode))
	assert.NotEmpty(t, resp.Message)
}

// Start & Shutdown 简单验证.
func TestHTTPServer_StartAndShutdown_ContextCancel(t *testing.T) {
	t.Parallel()

	// 使用本地随机端口
	cfg := HTTPConfig{
		Addr:       "127.0.0.1:0",
		PathPrefix: "/onebot",
	}
	server := NewHTTPServer(cfg, ActionRequestHandlerFunc(
		func(_ context.Context, _ *entity.ActionRequest) (*entity.ActionRawResponse, error) {
			return &entity.ActionRawResponse{Status: entity.StatusOK, Retcode: 0}, nil
		},
	))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		errCh <- server.Start(ctx)
	}()

	// 等待服务器启动片刻
	time.Sleep(100 * time.Millisecond)

	// 触发关闭
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
