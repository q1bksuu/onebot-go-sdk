package dispatcher

import "errors"

var (
	// ErrActionNotFound 表示 action 未注册 / 不存在，应映射为 404.
	ErrActionNotFound = errors.New("action not found")
)
