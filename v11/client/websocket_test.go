package client

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/q1bksuu/onebot-go-sdk/v11/entity"
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

		env, ok := resp.(*entity.ActionResponseEnvelope)
		require.True(t, ok)
		require.Equal(t, entity.StatusOK, env.Status)
		require.Equal(t, entity.ActionResponseRetcode(0), env.Retcode)
		require.JSONEq(t, `{"result":"ok"}`, string(env.Data))
		require.JSONEq(t, `"123"`, string(env.Echo))
	})

	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()

		req := `{"action":`
		resp := client.handleActionMessage(ctx, []byte(req))

		env, ok := resp.(*entity.ActionResponseEnvelope)
		require.True(t, ok)
		require.Equal(t, entity.StatusFailed, env.Status)
		require.Equal(t, entity.ActionResponseRetcode(1400), env.Retcode)
	})
}
