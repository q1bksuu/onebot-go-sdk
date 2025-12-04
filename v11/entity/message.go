package entity

import (
	"bytes"
	"encoding/json"
)

// MessageValue 表示 OneBot message 字段的值
// 可以是纯文本字符串或消息段数组
type MessageValue struct {
	// 如果 Type 为 "string"，则使用 StringValue
	// 如果 Type 为 "array"，则使用 ArrayValue
	Type        MessageValueType `json:"-"`
	StringValue string           `json:"-"`
	ArrayValue  []*Segment       `json:"-"`
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
// 用于在反序列化时自动选择正确的类型
func (m *MessageValue) UnmarshalJSON(data []byte) error {
	// 首先尝试作为字符串解析
	var str string
	if bytes.HasPrefix(data, []byte{'"'}) {
		if err := json.Unmarshal(data, &str); err == nil {
			m.Type = MessageValueTypeString
			m.StringValue = str
			return nil
		}
	}

	// 然后尝试作为数组解析
	var arr []*Segment
	if bytes.HasPrefix(data, []byte{'['}) {
		if err := json.Unmarshal(data, &arr); err == nil {
			m.Type = MessageValueTypeArray
			m.ArrayValue = arr
			return nil
		}
	}

	return nil
}

// MarshalJSON 实现 json.Marshaler 接口
// 用于在序列化时正确地输出值
func (m *MessageValue) MarshalJSON() ([]byte, error) {
	if m.Type == MessageValueTypeString {
		return json.Marshal(m.StringValue)
	}
	if m.Type == MessageValueTypeArray {
		return json.Marshal(m.ArrayValue)
	}
	return []byte{'n', 'u', 'l', 'l'}, nil
}
