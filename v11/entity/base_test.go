package entity

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStatusMetaUnmarshalJSON(t *testing.T) {
	t.Parallel()

	// 测试基本的 JSON 反序列化
	jsonData := []byte(`{
		"online": true,
		"good": false,
		"custom_field": "custom_value",
		"extra_number": 42
	}`)

	var status *StatusMeta

	// GetOrigin 应该返回 nil 而不是 panic
	result := status.GetOrigin("any_key")
	require.Nil(t, result)

	err := json.Unmarshal(jsonData, &status)
	require.NoError(t, err)

	// 验证已知字段
	require.True(t, status.Online)
	require.False(t, status.Good)

	// 验证原始字段获取
	require.Equal(t, true, status.GetOrigin("online"))
	require.Equal(t, false, status.GetOrigin("good"))
	require.Equal(t, "custom_value", status.GetOrigin("custom_field"))
	require.InEpsilon(t, float64(42), status.GetOrigin("extra_number"), 0)

	// 测试 SetOrigin
	status.SetOrigin("new_field", "new_value")
	require.Equal(t, "new_value", status.GetOrigin("new_field"))

	// 验证链式调用
	status.SetOrigin("chain1", "value1").SetOrigin("chain2", "value2")
	require.Equal(t, "value1", status.GetOrigin("chain1"))
	require.Equal(t, "value2", status.GetOrigin("chain2"))
}
