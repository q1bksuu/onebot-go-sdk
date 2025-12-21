package server

import (
	"context"

	"github.com/q1bksuu/onebot-go-sdk/v11/entity"
)

// Dispatcher 根据 action 路由到对应 handler.
type Dispatcher struct {
	handlers map[string]ActionHandler
}

var _ ActionRequestHandler = (*Dispatcher)(nil)

// NewDispatcher 创建分发器.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{handlers: make(map[string]ActionHandler)}
}

// Register 注册 action handler.
func (d *Dispatcher) Register(action string, h ActionHandler) {
	d.handlers[action] = h
}

// HandleActionRequest 调用对应 action handler.
func (d *Dispatcher) HandleActionRequest(
	ctx context.Context,
	req *entity.ActionRequest,
) (*entity.ActionRawResponse, error) {
	h, ok := d.handlers[req.Action]
	if !ok {
		return nil, ErrActionNotFound
	}

	return h(ctx, req.Params)
}
