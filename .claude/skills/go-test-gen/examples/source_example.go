// examples/source_example.go
// 基于 onebot-go-sdk 项目的实际代码示例

package entity

import (
	"encoding/json"
	"errors"
)

// StatusMeta 状态元信息（简化版示例）.
type StatusMeta struct {
	Online bool           `json:"online"`
	Good   bool           `json:"good"`
	origin map[string]any // 原始数据
}

// GetOrigin 获取原始字段值.
func (m *StatusMeta) GetOrigin(key string) any {
	if m.origin == nil {
		return nil
	}

	return m.origin[key]
}

// SetOrigin 设置原始字段值（支持链式调用）.
func (m *StatusMeta) SetOrigin(key string, value any) *StatusMeta {
	if m.origin == nil {
		m.origin = make(map[string]any)
	}
	m.origin[key] = value

	return m
}

// UnmarshalJSON 自定义 JSON 反序列化.
func (m *StatusMeta) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return errors.New("empty data")
	}

	type Alias StatusMeta
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// 保存原始数据
	if err := json.Unmarshal(data, &m.origin); err != nil {
		return err
	}

	return nil
}

// Message 消息段（简化版示例）.
type Message struct {
	Type string         `json:"type"`
	Data map[string]any `json:"data"`
}

// NewTextMessage 创建文本消息.
func NewTextMessage(text string) *Message {
	if text == "" {
		return nil
	}

	return &Message{
		Type: "text",
		Data: map[string]any{"text": text},
	}
}

// ValidateMessage 验证消息段.
func ValidateMessage(msg *Message) error {
	if msg == nil {
		return errors.New("message is nil")
	}
	if msg.Type == "" {
		return errors.New("message type is required")
	}
	if msg.Data == nil {
		return errors.New("message data is required")
	}

	return nil
}
