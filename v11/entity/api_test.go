package entity

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetVersionInfoResponseUnmarshalJSON(t *testing.T) {
	t.Parallel()

	// 测试基本的 JSON 反序列化
	jsonData := []byte(`{
		"app_name": "mirai-native",
		"app_version": "1.2.3",
		"protocol_version": "v11",
		"custom_field": "custom_value",
		"extra_number": 42
	}`)

	var response *GetVersionInfoResponse

	// GetOrigin 应该返回 nil 而不是 panic
	result := response.GetOrigin("any_key")
	require.Nil(t, result)

	err := json.Unmarshal(jsonData, &response)
	require.NoError(t, err)

	// 验证已知字段
	require.Equal(t, "mirai-native", response.AppName)
	require.Equal(t, "1.2.3", response.AppVersion)
	require.Equal(t, "v11", response.ProtocolVersion)

	// 验证原始字段获取
	require.Equal(t, "mirai-native", response.GetOrigin("app_name"))
	require.Equal(t, "1.2.3", response.GetOrigin("app_version"))
	require.Equal(t, "v11", response.GetOrigin("protocol_version"))
	require.Equal(t, "custom_value", response.GetOrigin("custom_field"))
	require.InEpsilon(t, float64(42), response.GetOrigin("extra_number"), 0)

	// 测试 SetOrigin
	response.SetOrigin("new_field", "new_value")
	require.Equal(t, "new_value", response.GetOrigin("new_field"))

	// 验证链式调用
	response.SetOrigin("chain1", "value1").SetOrigin("chain2", "value2")
	require.Equal(t, "value1", response.GetOrigin("chain1"))
	require.Equal(t, "value2", response.GetOrigin("chain2"))
}
