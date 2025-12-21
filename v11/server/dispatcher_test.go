package server

import (
	"context"
	"testing"

	"github.com/q1bksuu/onebot-go-sdk/v11/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDispatcher_RegisterAndHandle_Success(t *testing.T) {
	t.Parallel()

	dispatcher := NewDispatcher()

	called := false

	dispatcher.Register("ping", func(_ context.Context, params map[string]any) (*entity.ActionRawResponse, error) {
		called = true

		assert.Equal(t, "pong", params["msg"])

		return &entity.ActionRawResponse{
			Status:  entity.StatusOK,
			Retcode: 0,
			Message: "ok",
		}, nil
	})

	req := &entity.ActionRequest{
		Action: "ping",
		Params: map[string]any{
			"msg": "pong",
		},
	}

	raw, err := dispatcher.HandleActionRequest(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, raw)

	assert.True(t, called, "registered handler should be called")
	assert.Equal(t, entity.StatusOK, raw.Status)
	assert.Equal(t, entity.ActionResponseRetcode(0), raw.Retcode)
	assert.Equal(t, "ok", raw.Message)
}

func TestDispatcher_HandleActionRequest_NotFound(t *testing.T) {
	t.Parallel()

	dispatcher := NewDispatcher()

	req := &entity.ActionRequest{
		Action: "not_registered",
		Params: map[string]any{},
	}

	raw, err := dispatcher.HandleActionRequest(context.Background(), req)

	assert.Nil(t, raw)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrActionNotFound)
}

func TestDispatcher_Register_OverrideExisting(t *testing.T) {
	t.Parallel()

	dispatcher := NewDispatcher()

	dispatcher.Register("echo", func(_ context.Context, _ map[string]any) (*entity.ActionRawResponse, error) {
		return &entity.ActionRawResponse{
			Status:  entity.StatusOK,
			Retcode: 0,
			Message: "first",
		}, nil
	})

	// 再次注册同名 action，应覆盖旧的 handler
	dispatcher.Register("echo", func(_ context.Context, _ map[string]any) (*entity.ActionRawResponse, error) {
		return &entity.ActionRawResponse{
			Status:  entity.StatusOK,
			Retcode: 0,
			Message: "second",
		}, nil
	})

	req := &entity.ActionRequest{
		Action: "echo",
		Params: map[string]any{},
	}

	raw, err := dispatcher.HandleActionRequest(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, raw)
	assert.Equal(t, "second", raw.Message)
}
