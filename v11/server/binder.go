package server

import (
	"context"
	"fmt"

	"github.com/q1bksuu/onebot-go-sdk/v11/entity"
	"github.com/q1bksuu/onebot-go-sdk/v11/internal/util"
)

// Binder 将 params 绑定到具体请求结构，并调用业务处理.
type Binder[Req any, Resp any] struct {
	action string
	fn     func(ctx context.Context, req *Req) (*entity.ActionResponse[Resp], error)
}

// NewBinder 创建一个与 action 绑定的业务处理器.
func NewBinder[Req any, Resp any](
	action string,
	fn func(context.Context, *Req) (*entity.ActionResponse[Resp], error),
) *Binder[Req, Resp] {
	return &Binder[Req, Resp]{action: action, fn: fn}
}

// Action 返回绑定的 action 名.
func (b *Binder[Req, Resp]) Action() string { return b.action }

// Handler 返回 ActionHandler，完成绑定和调用.
func (b *Binder[Req, Resp]) Handler() ActionHandler {
	return func(ctx context.Context, params map[string]any) (*entity.ActionRawResponse, error) {
		var req Req

		err := util.JsonTagMapping(params, &req)
		if err != nil {
			return nil, fmt.Errorf("bind params to request failed: %w", err)
		}

		resp, err := b.fn(ctx, &req)
		if err != nil {
			return nil, err
		}

		return resp.ToActionRawResponse()
	}
}
