package server

import (
	"context"
	"fmt"

	"github.com/q1bksuu/onebot-go-sdk/v11/entity"
	"github.com/q1bksuu/onebot-go-sdk/v11/internal/util"
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

func APIFuncToActionHandler[Req any, Resp any](
	fn func(ctx context.Context, req *Req) (*entity.ActionResponse[Resp], error),
) ActionHandler {
	return func(ctx context.Context, params map[string]any) (*entity.ActionRawResponse, error) {
		var req Req

		err := util.JsonTagMapping(params, &req)
		if err != nil {
			return nil, fmt.Errorf("bind params to request failed: %w", err)
		}

		resp, err := fn(ctx, &req)
		if err != nil {
			return nil, err
		}

		return resp.ToActionRawResponse()
	}
}
