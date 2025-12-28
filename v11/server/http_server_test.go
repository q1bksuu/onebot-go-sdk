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
		Addr:          ":0",
		APIPathPrefix: "api", // 不带斜杠
	}
	server := NewHTTPServer(cfg, WithActionHandler(ActionRequestHandlerFunc(
		func(_ context.Context, req *entity.ActionRequest) (*entity.ActionRawResponse, error) {
			// 简单回显 action，便于断言
			require.Equal(t, "test_action", req.Action)

			return &entity.ActionRawResponse{
				Status:  entity.StatusOK,
				Retcode: 0,
				Message: "ok",
			}, nil
		},
	)))

	// APIPathPrefix 应被归一化为 "/api/"
	assert.Equal(t, "/api/", server.cfg.APIPathPrefix)

	// Handler() 应该是可用的 http.Handler
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/test_action", nil)
	server.Handler().ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)

	var resp entity.ActionRawResponse

	err := json.NewDecoder(recorder.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, entity.StatusOK, resp.Status)
	assert.Equal(t, entity.RetcodeSuccess, resp.Retcode)
}

func newTestServer(cfg HTTPConfig, handler ActionRequestHandlerFunc) *HTTPServer {
	if cfg.APIPathPrefix == "" {
		cfg.APIPathPrefix = "/"
	}

	return NewHTTPServer(cfg, WithActionHandler(handler))
}

func TestHTTPServer_HandleRoot_PathAndNotFound(t *testing.T) {
	t.Parallel()

	server := newTestServer(HTTPConfig{APIPathPrefix: "/onebot"},
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
		APIPathPrefix: "/onebot",
		AccessToken:   "secret",
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
		APIPathPrefix: "/onebot",
		AccessToken:   "secret",
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
	server := newTestServer(HTTPConfig{APIPathPrefix: "/onebot"}, handler)

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
	server := newTestServer(HTTPConfig{APIPathPrefix: "/onebot"}, handler)

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
	server := newTestServer(HTTPConfig{APIPathPrefix: "/onebot"}, handler)

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
			server := newTestServer(HTTPConfig{APIPathPrefix: "/onebot"}, handler)

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
		//nolint:nilnil // 测试代码中返回 nil, nil 用于测试默认行为
		return nil, nil
	}
	server := newTestServer(HTTPConfig{APIPathPrefix: "/onebot"}, handler)

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
		Addr:          "127.0.0.1:0",
		APIPathPrefix: "/onebot",
	}
	server := NewHTTPServer(cfg, WithActionHandler(ActionRequestHandlerFunc(
		func(_ context.Context, _ *entity.ActionRequest) (*entity.ActionRawResponse, error) {
			return &entity.ActionRawResponse{Status: entity.StatusOK, Retcode: 0}, nil
		},
	)))

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

// 事件处理测试.
func TestHTTPServer_EventPath_Registration(t *testing.T) {
	t.Parallel()

	cfg := HTTPConfig{
		Addr:      ":0",
		EventPath: "/event",
	}
	eventHandler := EventRequestHandlerFunc(func(_ context.Context, _ entity.Event) (map[string]any, error) {
		return map[string]any{"reply": "test"}, nil
	})

	server := NewHTTPServer(cfg, WithActionHandler(ActionRequestHandlerFunc(
		func(_ context.Context, _ *entity.ActionRequest) (*entity.ActionRawResponse, error) {
			return &entity.ActionRawResponse{Status: entity.StatusOK, Retcode: 0}, nil
		})), WithEventHandler(eventHandler))

	// EventPath 应该被正确注册

	// 测试事件路由是否注册
	recorder := httptest.NewRecorder()
	eventJSON := `{
		"time": 1515204254,
		"self_id": 10001000,
		"post_type": "message",
		"message_type": "private",
		"sub_type": "friend",
		"message_id": 12,
		"user_id": 12345678,
		"message": "Hello~",
		"raw_message": "Hello~",
		"font": 456,
		"sender": {
			"user_id": 12345678,
			"nickname": "A User",
			"sex": "male",
			"age": 18
		}
	}`
	req := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBufferString(eventJSON))
	req.Header.Set("Content-Type", "application/json")
	server.Handler().ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)

	var quickOp map[string]any

	err := json.NewDecoder(recorder.Body).Decode(&quickOp)
	require.NoError(t, err)
	assert.Equal(t, "test", quickOp["reply"])
}

func TestHTTPServer_EventPath_NoHandler_Returns204(t *testing.T) {
	t.Parallel()

	cfg := HTTPConfig{
		Addr:      ":0",
		EventPath: "/event",
	}
	server := NewHTTPServer(cfg, WithActionHandler(ActionRequestHandlerFunc(
		func(_ context.Context, _ *entity.ActionRequest) (*entity.ActionRawResponse, error) {
			return &entity.ActionRawResponse{Status: entity.StatusOK, Retcode: 0}, nil
		})))

	// 没有事件处理器时应该返回 204
	recorder := httptest.NewRecorder()
	eventJSON := `{
		"time": 1515204254,
		"self_id": 10001000,
		"post_type": "message",
		"message_type": "private",
		"sub_type": "friend",
		"message_id": 12,
		"user_id": 12345678,
		"message": "Hello~",
		"raw_message": "Hello~",
		"font": 456
	}`
	req := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBufferString(eventJSON))
	req.Header.Set("Content-Type", "application/json")
	server.Handler().ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestHTTPServer_EventPath_OnlyAcceptsPOST(t *testing.T) {
	t.Parallel()

	cfg := HTTPConfig{
		Addr:      ":0",
		EventPath: "/event",
	}
	eventHandler := EventRequestHandlerFunc(func(_ context.Context, _ entity.Event) (map[string]any, error) {
		//nolint:nilnil // 测试代码中返回 nil, nil 表示没有快速操作，这是预期的行为
		return nil, nil
	})

	server := NewHTTPServer(cfg, WithActionHandler(ActionRequestHandlerFunc(
		func(_ context.Context, _ *entity.ActionRequest) (*entity.ActionRawResponse, error) {
			//nolint:nilnil // 测试代码中返回 nil, nil 用于测试默认行为
			return nil, nil
		})), WithEventHandler(eventHandler))

	// GET 请求应该返回 405
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/event", nil)
	server.Handler().ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
}

func TestHTTPServer_EventPath_InvalidJSON_Returns400(t *testing.T) {
	t.Parallel()

	cfg := HTTPConfig{
		Addr:      ":0",
		EventPath: "/event",
	}
	eventHandler := EventRequestHandlerFunc(func(_ context.Context, _ entity.Event) (map[string]any, error) {
		//nolint:nilnil // 测试代码中返回 nil, nil 表示没有快速操作，这是预期的行为
		return nil, nil
	})

	server := NewHTTPServer(cfg, WithActionHandler(ActionRequestHandlerFunc(
		func(_ context.Context, _ *entity.ActionRequest) (*entity.ActionRawResponse, error) {
			//nolint:nilnil // 测试代码中返回 nil, nil 用于测试默认行为
			return nil, nil
		})), WithEventHandler(eventHandler))

	// 无效的 JSON 应该返回 400
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	server.Handler().ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestHTTPServer_EventPath_EmptyQuickOp_Returns204(t *testing.T) {
	t.Parallel()

	cfg := HTTPConfig{
		Addr:      ":0",
		EventPath: "/event",
	}
	eventHandler := EventRequestHandlerFunc(func(_ context.Context, _ entity.Event) (map[string]any, error) {
		//nolint:nilnil // 返回 nil，应该返回 204
		return nil, nil
	})

	server := NewHTTPServer(cfg, WithActionHandler(ActionRequestHandlerFunc(
		func(_ context.Context, _ *entity.ActionRequest) (*entity.ActionRawResponse, error) {
			//nolint:nilnil // 测试代码中返回 nil, nil 用于测试默认行为
			return nil, nil
		})), WithEventHandler(eventHandler))

	recorder := httptest.NewRecorder()
	eventJSON := `{
		"time": 1515204254,
		"self_id": 10001000,
		"post_type": "message",
		"message_type": "private",
		"sub_type": "friend",
		"message_id": 12,
		"user_id": 12345678,
		"message": "Hello~",
		"raw_message": "Hello~",
		"font": 456
	}`
	req := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBufferString(eventJSON))
	req.Header.Set("Content-Type", "application/json")
	server.Handler().ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
}

//nolint:gochecknoglobals
var eventPathTestCases = []struct {
	name       string
	eventJSON  string
	wantStatus int
	wantReply  string
}{
	{
		name: "private message",
		eventJSON: `{
			"time": 1515204254,
			"self_id": 10001000,
			"post_type": "message",
			"message_type": "private",
			"sub_type": "friend",
			"message_id": 12,
			"user_id": 12345678,
			"message": "Hello~",
			"raw_message": "Hello~",
			"font": 456
		}`,
		wantStatus: http.StatusOK,
		wantReply:  "private",
	},
	{
		name: "group upload",
		eventJSON: `{
			"time": 1515204254,
			"self_id": 10001000,
			"post_type": "notice",
			"notice_type": "group_upload",
			"group_id": 123456,
			"user_id": 789012,
			"file": {
				"id": "file1",
				"name": "test.txt",
				"size": 1024,
				"busid": 0
			}
		}`,
		wantStatus: http.StatusOK,
		wantReply:  "upload",
	},
	{
		name: "friend request",
		eventJSON: `{
			"time": 1515204254,
			"self_id": 10001000,
			"post_type": "request",
			"request_type": "friend",
			"user_id": 12345678,
			"comment": "test",
			"flag": "flag123"
		}`,
		wantStatus: http.StatusOK,
	},
	{
		name: "lifecycle event",
		eventJSON: `{
			"time": 1515204254,
			"self_id": 10001000,
			"post_type": "meta_event",
			"meta_event_type": "lifecycle",
			"sub_type": "enable"
		}`,
		wantStatus: http.StatusNoContent,
	},
}

func TestHTTPServer_EventPath_AllEventTypes(t *testing.T) {
	t.Parallel()

	cfg := HTTPConfig{
		Addr:      ":0",
		EventPath: "/event",
	}

	dispatcher := NewEventDispatcher()
	dispatcher.Register("message/private", func(_ context.Context, _ entity.Event) (map[string]any, error) {
		return map[string]any{"reply": "private"}, nil
	})
	dispatcher.Register("notice/group_upload", func(_ context.Context, _ entity.Event) (map[string]any, error) {
		return map[string]any{"reply": "upload"}, nil
	})
	dispatcher.Register("request/friend", func(_ context.Context, _ entity.Event) (map[string]any, error) {
		return map[string]any{"approve": true}, nil
	})
	dispatcher.Register("meta_event/lifecycle", func(_ context.Context, _ entity.Event) (map[string]any, error) {
		return map[string]any{}, nil
	})

	server := NewHTTPServer(cfg, WithActionHandler(ActionRequestHandlerFunc(
		func(_ context.Context, _ *entity.ActionRequest) (*entity.ActionRawResponse, error) {
			return &entity.ActionRawResponse{Status: entity.StatusOK, Retcode: 0}, nil
		})), WithEventHandler(dispatcher))

	for _, testCase := range eventPathTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBufferString(testCase.eventJSON))
			req.Header.Set("Content-Type", "application/json")
			server.Handler().ServeHTTP(recorder, req)

			assert.Equal(t, testCase.wantStatus, recorder.Code)

			if testCase.wantReply != "" {
				var quickOp map[string]any

				err := json.NewDecoder(recorder.Body).Decode(&quickOp)
				require.NoError(t, err)
				assert.Equal(t, testCase.wantReply, quickOp["reply"])
			}
		})
	}
}

func TestEventDispatcher_Routing(t *testing.T) {
	t.Parallel()

	dispatcher := NewEventDispatcher()

	// 注册不同级别的处理器
	var calledKey string

	dispatcher.Register("message", func(_ context.Context, _ entity.Event) (map[string]any, error) {
		calledKey = "message"

		//nolint:nilnil // 测试代码中返回 nil, nil 表示没有快速操作
		return nil, nil
	})
	dispatcher.Register("message/private", func(_ context.Context, _ entity.Event) (map[string]any, error) {
		calledKey = "message/private"

		//nolint:nilnil // 测试代码中返回 nil, nil 表示没有快速操作
		return nil, nil
	})
	dispatcher.Register("message/private/friend", func(_ context.Context, _ entity.Event) (map[string]any, error) {
		calledKey = "message/private/friend"

		//nolint:nilnil // 测试代码中返回 nil, nil 表示没有快速操作
		return nil, nil
	})

	// 测试最具体的处理器被调用
	event := &entity.PrivateMessageEvent{
		Time:        1515204254,
		SelfId:      10001000,
		PostType:    entity.EventPostTypeMessage,
		MessageType: entity.EventMessageTypePrivate,
		SubType:     entity.EventPrivateMessageSubTypeFriend,
		MessageId:   12,
		UserId:      12345678,
	}

	_, _ = dispatcher.HandleEvent(context.Background(), event)

	assert.Equal(t, "message/private/friend", calledKey)
}
