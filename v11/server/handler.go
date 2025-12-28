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

// EventHandler 处理事件，返回快速操作响应（可选）.
// 如果返回 nil，则返回 204 No Content.
type EventHandler func(ctx context.Context, event entity.Event) (map[string]any, error)

// EventRequestHandler 处理事件请求.
type EventRequestHandler interface {
	HandleEvent(ctx context.Context, event entity.Event) (map[string]any, error)
}

// EventRequestHandlerFunc 适配函数.
type EventRequestHandlerFunc func(ctx context.Context, event entity.Event) (map[string]any, error)

func (f EventRequestHandlerFunc) HandleEvent(ctx context.Context, event entity.Event) (map[string]any, error) {
	return f(ctx, event)
}
