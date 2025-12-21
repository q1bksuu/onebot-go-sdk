package server

import "errors"

var (
	// ErrActionNotFound 表示 action 未注册 / 不存在，应映射为 404.
	ErrActionNotFound = errors.New("action not found")
	// ErrBadRequest 表示参数解析/校验失败，应映射为 400.
	ErrBadRequest = errors.New("bad request")
)
