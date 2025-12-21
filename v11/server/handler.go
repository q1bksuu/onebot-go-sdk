package server

import (
	"context"

	"github.com/q1bksuu/onebot-go-sdk/v11/entity"
)

// ActionRequestHandler 处理动作请求.
type ActionRequestHandler interface {
	HandleActionRequest(ctx context.Context, req *entity.ActionRequest) (*entity.ActionRawResponse, error)
}

// ActionRequestHandlerFunc 适配函数.
type ActionRequestHandlerFunc func(ctx context.Context, req *entity.ActionRequest) (*entity.ActionRawResponse, error)

func (f ActionRequestHandlerFunc) HandleActionRequest(
	ctx context.Context,
	req *entity.ActionRequest,
) (*entity.ActionRawResponse, error) {
	return f(ctx, req)
}

// ActionHandler 处理具体 action.
type ActionHandler func(ctx context.Context, params map[string]any) (*entity.ActionRawResponse, error)
