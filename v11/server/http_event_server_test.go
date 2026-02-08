package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/q1bksuu/onebot-go-sdk/v11/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testEventService struct {
	UnimplementedOneBotEventService

	called atomic.Bool
	reply  string
}

func (s *testEventService) HandlePrivateMessage(
	_ context.Context,
	_ *entity.PrivateMessageEvent,
) (map[string]any, error) {
	s.called.Store(true)

	return map[string]any{"reply": s.reply}, nil
}

func sendEventRequest(
	t *testing.T,
	server *HTTPServer,
	path string,
	payload string,
) *httptest.ResponseRecorder {
	t.Helper()

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	server.Handler().ServeHTTP(recorder, req)

	return recorder
}

func TestNewHTTPEventServer_EventPathAndHandler(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		eventPath string
		wantPath  string
	}{
		{
			name:     "default_event_path",
			wantPath: "/event",
		},
		{
			name:      "custom_event_path",
			eventPath: "/hooks/onebot",
			wantPath:  "/hooks/onebot",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			cfg := HTTPConfig{
				Addr:      ":0",
				EventPath: testCase.eventPath,
			}

			service := &testEventService{reply: "ok"}
			server := NewHTTPEventServer(cfg, service)

			assert.Equal(t, testCase.wantPath, server.cfg.EventPath)

			const payload = `{
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

			recorder := sendEventRequest(t, server, testCase.wantPath, payload)
			assert.Equal(t, http.StatusOK, recorder.Code)

			var quickOp map[string]any

			err := json.NewDecoder(recorder.Body).Decode(&quickOp)
			require.NoError(t, err)
			assert.Equal(t, "ok", quickOp["reply"])
			assert.True(t, service.called.Load())
		})
	}
}

func TestNewHTTPEventServer_ActionRequestReturnsNotFound(t *testing.T) {
	t.Parallel()

	cfg := HTTPConfig{Addr: ":0"}
	server := NewHTTPEventServer(cfg, &testEventService{reply: "ok"})

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/send_private_msg", nil)
	server.Handler().ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
}
