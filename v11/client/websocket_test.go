package client

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

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
