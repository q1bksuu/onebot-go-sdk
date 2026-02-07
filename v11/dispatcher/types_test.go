package dispatcher

import (
	"context"
	"testing"

	"github.com/q1bksuu/onebot-go-sdk/v11/entity"
	"github.com/stretchr/testify/require"
)

func TestActionRequestHandlerFunc_HandleActionRequest(t *testing.T) {
	t.Parallel()

	called := false

	handler := ActionRequestHandlerFunc(
		func(_ context.Context, req *entity.ActionRequest) (*entity.ActionRawResponse, error) {
			called = true

			require.Equal(t, "ping", req.Action)

			return &entity.ActionRawResponse{
				Status:  entity.StatusOK,
				Retcode: 0,
				Message: "ok",
			}, nil
		},
	)

	resp, err := handler.HandleActionRequest(context.Background(), &entity.ActionRequest{Action: "ping"})
	require.NoError(t, err)
	require.True(t, called)
	require.Equal(t, entity.StatusOK, resp.Status)
	require.Equal(t, entity.ActionResponseRetcode(0), resp.Retcode)
	require.Equal(t, "ok", resp.Message)
}
