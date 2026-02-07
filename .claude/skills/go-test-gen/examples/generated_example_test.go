// examples/generated_test_example.go
// go-test-gen 生成的测试示例（基于 source_example.go）

package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStatusMeta_SetOrigin 表驱动测试示例.
func TestStatusMeta_SetOrigin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		initial   map[string]any
		key       string
		value     any
		wantValue any
	}{
		{
			name:      "nil_origin_creates_map",
			initial:   nil,
			key:       "custom",
			value:     "value",
			wantValue: "value",
		},
		{
			name:      "existing_origin_adds_new_key",
			initial:   map[string]any{"old": "data"},
			key:       "new",
			value:     123,
			wantValue: 123,
		},
		{
			name:      "overwrites_existing_key",
			initial:   map[string]any{"key": "old"},
			key:       "key",
			value:     "new",
			wantValue: "new",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			meta := &StatusMeta{origin: tt.initial}

			// Act
			result := meta.SetOrigin(tt.key, tt.value)

			// Assert
			assert.Equal(t, meta, result, "should return self for chaining")
			assert.Equal(t, tt.wantValue, meta.GetOrigin(tt.key))
		})
	}
}

// TestStatusMeta_GetOrigin 简单测试示例.
func TestStatusMeta_GetOrigin_NilOrigin(t *testing.T) {
	t.Parallel()

	// Arrange
	meta := &StatusMeta{}

	// Act
	value := meta.GetOrigin("non_existent")

	// Assert
	assert.Nil(t, value)
}

// TestValidateMessage 错误路径测试示例.
func TestValidateMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		msg     *Message
		wantErr string
	}{
		{
			name:    "nil_message_returns_error",
			msg:     nil,
			wantErr: "message is nil",
		},
		{
			name:    "empty_type_returns_error",
			msg:     &Message{Type: "", Data: map[string]any{}},
			wantErr: "message type is required",
		},
		{
			name:    "nil_data_returns_error",
			msg:     &Message{Type: "text", Data: nil},
			wantErr: "message data is required",
		},
		{
			name:    "valid_message_passes",
			msg:     &Message{Type: "text", Data: map[string]any{"text": "hello"}},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Act
			err := ValidateMessage(tt.msg)

			// Assert
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}

			assert.NoError(t, err)
		})
	}
}

// TestNewTextMessage 边界条件测试示例.
func TestNewTextMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		text string
		want *Message
	}{
		{
			name: "empty_text_returns_nil",
			text: "",
			want: nil,
		},
		{
			name: "valid_text_creates_message",
			text: "hello",
			want: &Message{
				Type: "text",
				Data: map[string]any{"text": "hello"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Act
			got := NewTextMessage(tt.text)

			// Assert
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestStatusMeta_UnmarshalJSON JSON 反序列化测试示例.
func TestStatusMeta_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    string
		wantErr bool
		check   func(t *testing.T, meta *StatusMeta)
	}{
		{
			name:    "empty_data_returns_error",
			data:    "",
			wantErr: true,
		},
		{
			name:    "invalid_json_returns_error",
			data:    "{invalid}",
			wantErr: true,
		},
		{
			name:    "valid_json_unmarshals_successfully",
			data:    `{"online":true,"good":false,"custom":"value"}`,
			wantErr: false,
			check: func(t *testing.T, meta *StatusMeta) {
				assert.True(t, meta.Online)
				assert.False(t, meta.Good)
				assert.Equal(t, "value", meta.GetOrigin("custom"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			meta := &StatusMeta{}

			// Act
			err := meta.UnmarshalJSON([]byte(tt.data))

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.check != nil {
				tt.check(t, meta)
			}
		})
	}
}

// BenchmarkStatusMeta_SetOrigin 基准测试示例.
func BenchmarkStatusMeta_SetOrigin(b *testing.B) {
	meta := &StatusMeta{}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		meta.SetOrigin("key", "value")
	}
}

// BenchmarkValidateMessage 基准测试示例.
func BenchmarkValidateMessage(b *testing.B) {
	msg := &Message{
		Type: "text",
		Data: map[string]any{"text": "hello"},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ValidateMessage(msg)
	}
}
