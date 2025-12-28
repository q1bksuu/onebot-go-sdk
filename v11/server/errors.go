package server

import "errors"

var (
	// ErrActionNotFound 表示 action 未注册 / 不存在，应映射为 404.
	ErrActionNotFound = errors.New("action not found")
	// ErrBadRequest 表示参数解析/校验失败，应映射为 400.
	ErrBadRequest = errors.New("bad request")
	// ErrUniversalClientURLEmpty 表示 universal client URL 为空.
	ErrUniversalClientURLEmpty = errors.New("universal client URL is empty")
	// ErrMissingTypeField 表示缺少类型字段.
	ErrMissingTypeField = errors.New("missing type field")
	// ErrUnknownEventType 表示未知的事件类型.
	ErrUnknownEventType = errors.New("unknown event type")
	// ErrInvalidEventTreeStructure 表示事件树结构无效.
	ErrInvalidEventTreeStructure = errors.New("invalid event tree structure")
	// ErrMissingOrInvalidPostType 表示缺少或无效的 post_type 字段.
	ErrMissingOrInvalidPostType = errors.New("missing or invalid post_type field")
	// ErrUnknownPostType 表示未知的 post_type.
	ErrUnknownPostType = errors.New("unknown post_type")
	// ErrNoEventHandler 表示没有匹配的事件处理器.
	ErrNoEventHandler = errors.New("no event handler")
)
