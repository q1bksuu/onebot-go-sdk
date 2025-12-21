package server

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/q1bksuu/onebot-go-sdk/v11/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errBizError = errors.New("biz error")

func TestBinder_ActionAndHandler_Success(t *testing.T) {
	t.Parallel()

	type testReq struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	type testResp struct {
		OK bool `json:"ok"`
	}

	handler1 := func(_ context.Context, req *testReq) (*entity.ActionResponse[testResp], error) {
		// 验证参数绑定
		assert.Equal(t, "alice", req.Name)
		assert.Equal(t, 18, req.Age)

		return &entity.ActionResponse[testResp]{
			Status:  entity.StatusOK,
			Retcode: 0,
			Data:    &testResp{OK: true},
			Message: "ok",
		}, nil
	}
	b := NewBinder[testReq, testResp]("test_action", handler1)

	assert.Equal(t, "test_action", b.Action())

	handler := b.Handler()
	params := map[string]any{
		"name": "alice",
		"age":  18,
	}

	raw, err := handler(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, raw)

	assert.Equal(t, entity.StatusOK, raw.Status)
	assert.Equal(t, entity.ActionResponseRetcode(0), raw.Retcode)

	var decoded testResp
	require.NoError(t, json.Unmarshal(raw.Data, &decoded))
	assert.True(t, decoded.OK)
}

func TestBinder_Handler_PropagatesFuncError(t *testing.T) {
	t.Parallel()

	type testReq struct {
		ID int64 `json:"id"`
	}

	type testResp struct{}

	handler2 := func(_ context.Context, _ *testReq) (*entity.ActionResponse[testResp], error) {
		return nil, errBizError
	}
	b := NewBinder[testReq, testResp]("do_something", handler2)

	handler := b.Handler()
	raw, err := handler(context.Background(), map[string]any{"id": 1})

	assert.Nil(t, raw)
	require.Error(t, err)
	assert.ErrorIs(t, err, errBizError)
}
