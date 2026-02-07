package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/q1bksuu/onebot-go-sdk/v11/dispatcher"
	"github.com/q1bksuu/onebot-go-sdk/v11/entity"
	"github.com/stretchr/testify/require"
)

type testError string

func (e testError) Error() string {
	return string(e)
}

const (
	ErrBadRequest testError = "bad request"
	ErrOther      testError = "other error"
)

func TestHandleActionMessageSuccess(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	handler := dispatcher.ActionRequestHandlerFunc(
		func(_ context.Context, req *entity.ActionRequest) (*entity.ActionRawResponse, error) {
			require.Equal(t, "ping", req.Action)
			require.Equal(t, "pong", req.Params["msg"])

			return &entity.ActionRawResponse{
				Status:  entity.StatusOK,
				Retcode: 0,
				Data:    json.RawMessage(`{"ok":true}`),
				Message: "ok",
			}, nil
		},
	)

	payload := []byte(`{"action":"ping","params":{"msg":"pong"},"echo":"e1"}`)
	resp := HandleActionMessage(ctx, payload, handler, ErrBadRequest)

	require.Equal(t, entity.StatusOK, resp.Status)
	require.Equal(t, entity.ActionResponseRetcode(0), resp.Retcode)
	require.Equal(t, "ok", resp.Message)
	require.JSONEq(t, `{"ok":true}`, string(resp.Data))
	require.JSONEq(t, `"e1"`, string(resp.Echo))
}

func TestHandleActionMessageInvalidJSON(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	handler := dispatcher.ActionRequestHandlerFunc(
		func(_ context.Context, _ *entity.ActionRequest) (*entity.ActionRawResponse, error) {
			return &entity.ActionRawResponse{Status: entity.StatusOK}, nil
		},
	)

	resp := HandleActionMessage(ctx, []byte(`{"action":`), handler, ErrBadRequest)
	require.Equal(t, entity.StatusFailed, resp.Status)
	require.Equal(t, entity.ActionResponseRetcode(1400), resp.Retcode)
	require.Equal(t, "invalid json", resp.Message)
}

func TestHandleActionMessageHandlerNotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	handler := dispatcher.ActionRequestHandlerFunc(
		func(_ context.Context, _ *entity.ActionRequest) (*entity.ActionRawResponse, error) {
			return nil, dispatcher.ErrActionNotFound
		},
	)

	payload := []byte(`{"action":"missing","params":{},"echo":"e2"}`)
	resp := HandleActionMessage(ctx, payload, handler, ErrBadRequest)

	require.Equal(t, entity.StatusFailed, resp.Status)
	require.Equal(t, entity.ActionResponseRetcode(1404), resp.Retcode)
	require.Equal(t, dispatcher.ErrActionNotFound.Error(), resp.Message)
	require.JSONEq(t, `"e2"`, string(resp.Echo))
}

func TestHandleActionMessageHandlerBadRequest(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	handler := dispatcher.ActionRequestHandlerFunc(
		func(_ context.Context, _ *entity.ActionRequest) (*entity.ActionRawResponse, error) {
			return nil, fmt.Errorf("wrap: %w", ErrBadRequest)
		},
	)

	resp := HandleActionMessage(ctx, []byte(`{"action":"bad","params":{}}`), handler, ErrBadRequest)

	require.Equal(t, entity.StatusFailed, resp.Status)
	require.Equal(t, entity.ActionResponseRetcode(1400), resp.Retcode)
	require.Equal(t, "wrap: "+ErrBadRequest.Error(), resp.Message)
}

func TestHandleActionMessageHandlerOtherError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	handler := dispatcher.ActionRequestHandlerFunc(
		func(_ context.Context, _ *entity.ActionRequest) (*entity.ActionRawResponse, error) {
			return nil, ErrOther
		},
	)

	resp := HandleActionMessage(ctx, []byte(`{"action":"boom","params":{}}`), handler, ErrBadRequest)

	require.Equal(t, entity.StatusFailed, resp.Status)
	require.Equal(t, entity.ActionResponseRetcode(1500), resp.Retcode)
	require.Equal(t, ErrOther.Error(), resp.Message)
}

func TestHandleActionMessageNilResponse(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	handler := dispatcher.ActionRequestHandlerFunc(
		func(_ context.Context, _ *entity.ActionRequest) (*entity.ActionRawResponse, error) {
			return nil, nil //nolint:nilnil // intentional to cover nil response path
		},
	)

	payload := []byte(`{"action":"nil","params":{},"echo":"e3"}`)
	resp := HandleActionMessage(ctx, payload, handler, ErrBadRequest)

	require.Equal(t, entity.StatusFailed, resp.Status)
	require.Equal(t, entity.ActionResponseRetcode(-1), resp.Retcode)
	require.Equal(t, "empty response", resp.Message)
	require.JSONEq(t, `"e3"`, string(resp.Echo))
}
