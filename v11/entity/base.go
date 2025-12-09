//go:generate go run ../cmd/entity-gen
package entity

import (
	"fmt"

	"github.com/q1bksuu/onebot-go-sdk/v11/internal/util"
)

type StatusMeta struct {
	// 当前 QQ 在线，`null` 表示无法查询到在线状态
	Online bool `json:"online"`
	// 状态符合预期，意味着各模块正常运行、功能正常，且 QQ 在线
	Good bool `json:"good"`
	// 原始字段的值，只能通过 GetOrigin / SetOrigin 方法访问
	origin map[string]any
}

// GetOrigin 获取指定 key 的原始字段值.
func (r *StatusMeta) GetOrigin(key string) any {
	if r == nil || r.origin == nil {
		return nil
	}

	return r.origin[key]
}

// SetOrigin 设置指定 key 的原始字段值.
func (r *StatusMeta) SetOrigin(key string, value any) *StatusMeta {
	if r == nil {
		return r
	}

	if r.origin == nil {
		r.origin = make(map[string]any)
	}

	r.origin[key] = value

	return r
}

// UnmarshalJSON 自定义反序列化，同时捕获所有原始字段数据.
func (r *StatusMeta) UnmarshalJSON(data []byte) error {
	err := util.JsonUnmarshalToMapAndStruct(data, r, &r.origin)
	if err != nil {
		return fmt.Errorf("failed to unmarshal StatusMeta: %w", err)
	}

	return nil
}

type GroupAnonymousUser struct {
	// 匿名用户 ID
	Id int64 `json:"id"`
	// 匿名用户名称
	Name string `json:"name"`
	// 匿名用户 flag，在调用禁言 API 时需要传入
	Flag string `json:"flag"`
}
