package util

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestJsonUnmarshalToMapAndStructBasic 测试基本的反序列化功能.
func TestJsonUnmarshalToMapAndStructBasic(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	jsonData := []byte(`{"name": "Alice", "age": 30}`)

	var (
		person    Person
		originMap map[string]any
	)

	err := JsonUnmarshalToMapAndStruct(jsonData, &person, &originMap)
	require.NoError(t, err)

	require.Equal(t, "Alice", person.Name)
	require.Equal(t, 30, person.Age)
	require.Equal(t, "Alice", originMap["name"])
	require.InEpsilon(t, float64(30), originMap["age"], 0)
}

// TestUnmarshalToMapAndStructNested 测试嵌套结构的反序列化.
func TestUnmarshalToMapAndStructNested(t *testing.T) {
	t.Parallel()

	type Address struct {
		City string `json:"city"`
	}

	type Person struct {
		Name    string  `json:"name"`
		Address Address `json:"address"`
	}

	jsonData := []byte(`{
		"name": "Bob",
		"address": {
			"city": "New York"
		}
	}`)

	var (
		person    Person
		originMap map[string]any
	)

	err := JsonUnmarshalToMapAndStruct(jsonData, &person, &originMap)
	require.NoError(t, err)

	require.Equal(t, "Bob", person.Name)
	require.Equal(t, "New York", person.Address.City)

	addressMap, ok := originMap["address"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "New York", addressMap["city"])
}

// TestUnmarshalToMapAndStructExtraFields 测试包含额外字段的 JSON.
func TestUnmarshalToMapAndStructExtraFields(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name string `json:"name"`
	}

	jsonData := []byte(`{
		"name": "Charlie",
		"email": "charlie@example.com",
		"phone": "123-456-7890"
	}`)

	var (
		person    Person
		originMap map[string]any
	)

	err := JsonUnmarshalToMapAndStruct(jsonData, &person, &originMap)
	require.NoError(t, err)

	require.Equal(t, "Charlie", person.Name)
	require.Equal(t, "charlie@example.com", originMap["email"])
	require.Equal(t, "123-456-7890", originMap["phone"])
}

// TestUnmarshalToMapAndStructArray 测试数组类型.
func TestUnmarshalToMapAndStructArray(t *testing.T) {
	t.Parallel()

	type Data struct {
		Items []string `json:"items"`
		Count int      `json:"count"`
	}

	jsonData := []byte(`{
		"items": ["apple", "banana", "cherry"],
		"count": 3
	}`)

	var (
		data      Data
		originMap map[string]any
	)

	err := JsonUnmarshalToMapAndStruct(jsonData, &data, &originMap)
	require.NoError(t, err)

	require.Len(t, data.Items, 3)
	require.Equal(t, "apple", data.Items[0])

	itemsArray, ok := originMap["items"].([]any)
	require.True(t, ok)
	require.Len(t, itemsArray, 3)
}

// TestUnmarshalToMapAndStructEmptyJSON 测试空 JSON 对象.
func TestUnmarshalToMapAndStructEmptyJSON(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	jsonData := []byte(`{}`)

	var (
		person    Person
		originMap map[string]any
	)

	err := JsonUnmarshalToMapAndStruct(jsonData, &person, &originMap)
	require.NoError(t, err)

	require.Empty(t, person.Name)
	require.Zero(t, person.Age)
	require.Empty(t, originMap)
}

// TestUnmarshalToMapAndStructNullValues 测试包含 null 值的 JSON.
func TestUnmarshalToMapAndStructNullValues(t *testing.T) {
	t.Parallel()

	type Data struct {
		Name     *string `json:"name"`
		Optional *string `json:"optional"`
	}

	jsonData := []byte(`{
		"name": "David",
		"optional": null
	}`)

	var (
		data      Data
		originMap map[string]any
	)

	err := JsonUnmarshalToMapAndStruct(jsonData, &data, &originMap)
	require.NoError(t, err)

	require.NotNil(t, data.Name)
	require.Equal(t, "David", *data.Name)
	require.Nil(t, data.Optional)
	require.Nil(t, originMap["optional"])
}

// TestUnmarshalToMapAndStructInvalidJSON 测试无效的 JSON.
func TestUnmarshalToMapAndStructInvalidJSON(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name string `json:"name"`
	}

	jsonData := []byte(`{invalid json}`)

	var (
		person    Person
		originMap map[string]any
	)

	err := JsonUnmarshalToMapAndStruct(jsonData, &person, &originMap)
	require.Error(t, err)
}

// TestJsonUnmarshalToMapAndStruct_mapping_error 测试映射失败错误.
func TestJsonUnmarshalToMapAndStruct_mapping_error(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name string `json:"name"`
	}

	jsonData := []byte(`{"name": "Alice"}`)

	var (
		person    Person
		originMap map[string]any
	)

	err := JsonUnmarshalToMapAndStruct(jsonData, person, &originMap)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to map JSON data to struct")
}

// TestJsonTagMapping_invalid_result 测试非法目标类型.
func TestJsonTagMapping_invalid_result(t *testing.T) {
	t.Parallel()

	source := map[string]any{"name": "Alice"}

	type Person struct {
		Name string `json:"name"`
	}

	var person Person

	err := JsonTagMapping(source, person)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create decoder")
}

// TestJsonTagMapping_decode_error 测试解码失败.
func TestJsonTagMapping_decode_error(t *testing.T) {
	t.Parallel()

	source := []any{"not", "a", "map"}

	type Person struct {
		Name string `json:"name"`
	}

	var person Person

	err := JsonTagMapping(source, &person)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to decode data to struct")
}

// TestUnmarshalToMapAndStructTypeConversion 测试类型转换.
func TestUnmarshalToMapAndStructTypeConversion(t *testing.T) {
	t.Parallel()

	type Data struct {
		Count  int     `json:"count"`
		Price  float64 `json:"price"`
		Active bool    `json:"active"`
	}

	jsonData := []byte(`{
		"count": 42,
		"price": 19.99,
		"active": true
	}`)

	var (
		data      Data
		originMap map[string]any
	)

	err := JsonUnmarshalToMapAndStruct(jsonData, &data, &originMap)
	require.NoError(t, err)

	require.Equal(t, 42, data.Count)
	require.InEpsilon(t, 19.99, data.Price, 0)
	require.True(t, data.Active)
}

func TestNormalizePath(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty", in: "", want: ""},
		{name: "root", in: "/", want: ""},
		{name: "no_slash", in: "api", want: "/api"},
		{name: "trim_both", in: "/api/v1/", want: "/api/v1"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tc.want, NormalizePath(tc.in))
		})
	}
}

// TestUnmarshalToMapAndStructUnicodeCharacters 测试 Unicode 字符.
func TestUnmarshalToMapAndStructUnicodeCharacters(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name string `json:"name"`
	}

	jsonData := []byte(`{"name": "Foo"}`)

	var (
		person    Person
		originMap map[string]any
	)

	err := JsonUnmarshalToMapAndStruct(jsonData, &person, &originMap)
	require.NoError(t, err)

	expectedName := "Foo"
	require.Equal(t, expectedName, person.Name)
	require.Equal(t, expectedName, originMap["name"])
}

// TestUnmarshalToMapAndStructLargeNumbers 测试大数字.
func TestUnmarshalToMapAndStructLargeNumbers(t *testing.T) {
	t.Parallel()

	type Data struct {
		BigInt   int64   `json:"big_int"`
		BigFloat float64 `json:"big_float"`
	}

	jsonData := []byte(`{
		"big_int": 9007199254740991,
		"big_float": 1.7976931348623157e+308
	}`)

	var (
		data      Data
		originMap map[string]any
	)

	err := JsonUnmarshalToMapAndStruct(jsonData, &data, &originMap)
	require.NoError(t, err)

	require.Equal(t, int64(9007199254740991), data.BigInt)
	require.InEpsilon(t, 1.7976931348623157e+308, data.BigFloat, 0)

	bigIntVal, ok := originMap["big_int"].(float64)
	require.True(t, ok)
	require.InEpsilon(t, 9007199254740991.0, bigIntVal, 0)

	bigFloatVal, ok := originMap["big_float"].(float64)
	require.True(t, ok)
	require.InEpsilon(t, 1.7976931348623157e+308, bigFloatVal, 0)
}

// TestUnmarshalToMapAndStructMapPointerNil 测试 map 指针为 nil.
func TestUnmarshalToMapAndStructMapPointerNil(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name string `json:"name"`
	}

	jsonData := []byte(`{"name": "Eve"}`)

	var (
		person Person
		nilMap *map[string]any
	)

	err := JsonUnmarshalToMapAndStruct(jsonData, &person, nilMap)
	require.Error(t, err)
	require.Equal(t, "map pointer cannot be nil", err.Error())
}

// TestUnmarshalToMapAndStructComplexStructure 测试复杂的嵌套结构.
func TestUnmarshalToMapAndStructComplexStructure(t *testing.T) {
	t.Parallel()

	type Contact struct {
		Email string `json:"email"`
		Phone string `json:"phone"`
	}

	type Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
	}

	type Person struct {
		Name    string  `json:"name"`
		Contact Contact `json:"contact"`
		Address Address `json:"address"`
	}

	jsonData := []byte(`{
		"name": "Frank",
		"contact": {
			"email": "frank@example.com",
			"phone": "555-1234"
		},
		"address": {
			"street": "123 Main St",
			"city": "Springfield"
		}
	}`)

	var (
		person    Person
		originMap map[string]any
	)

	err := JsonUnmarshalToMapAndStruct(jsonData, &person, &originMap)
	require.NoError(t, err)

	require.Equal(t, "Frank", person.Name)
	require.Equal(t, "frank@example.com", person.Contact.Email)
	require.Equal(t, "123 Main St", person.Address.Street)

	require.Equal(t, "Frank", originMap["name"])

	contactMap, ok := originMap["contact"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "frank@example.com", contactMap["email"])
	require.Equal(t, "555-1234", contactMap["phone"])

	addressMap, ok := originMap["address"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "123 Main St", addressMap["street"])
	require.Equal(t, "Springfield", addressMap["city"])
}

// TestUnmarshalToMapAndStructWithSpecialCharacters 测试包含特殊字符的字符串.
func TestUnmarshalToMapAndStructWithSpecialCharacters(t *testing.T) {
	t.Parallel()

	type Data struct {
		Text string `json:"text"`
	}

	jsonData := []byte(`{"text": "Line1\nLine2\tTabbed"}`)

	var (
		data      Data
		originMap map[string]any
	)

	err := JsonUnmarshalToMapAndStruct(jsonData, &data, &originMap)
	require.NoError(t, err)

	expected := "Line1\nLine2\tTabbed"
	require.Equal(t, expected, data.Text)
	require.Equal(t, expected, originMap["text"])
}

// TestUnmarshalToMapAndStructBooleanValues 测试布尔值.
func TestUnmarshalToMapAndStructBooleanValues(t *testing.T) {
	t.Parallel()

	type Flags struct {
		Enabled bool `json:"enabled"`
		Active  bool `json:"active"`
	}

	jsonData := []byte(`{
		"enabled": true,
		"active": false
	}`)

	var (
		flags     Flags
		originMap map[string]any
	)

	err := JsonUnmarshalToMapAndStruct(jsonData, &flags, &originMap)
	require.NoError(t, err)

	require.True(t, flags.Enabled)
	require.False(t, flags.Active)

	enabledVal, ok := originMap["enabled"].(bool)
	require.True(t, ok)
	require.True(t, enabledVal)

	activeVal, ok := originMap["active"].(bool)
	require.True(t, ok)
	require.False(t, activeVal)
}

// BenchmarkUnmarshalToMapAndStructSmallPayload 小负载性能测试.
func BenchmarkUnmarshalToMapAndStructSmallPayload(b *testing.B) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	jsonData := []byte(`{"name": "Alice", "age": 30}`)

	b.ResetTimer()

	for range b.N {
		var (
			person    Person
			originMap map[string]any
		)

		err := JsonUnmarshalToMapAndStruct(jsonData, &person, &originMap)
		require.NoError(b, err)
		require.Equal(b, "Alice", person.Name)
		require.Equal(b, 30, person.Age)
		require.Equal(b, "Alice", originMap["name"])
		require.InEpsilon(b, float64(30), originMap["age"], 0)
	}
}

// BenchmarkUnmarshalToMapAndStructMediumPayload 中等负载性能测试.
func BenchmarkUnmarshalToMapAndStructMediumPayload(b *testing.B) {
	type Person struct {
		Name    string
		Age     int
		Email   string
		Phone   string
		Address string
		City    string
		Country string
	}

	jsonData := []byte(`{
		"name": "Alice",
		"age": 30,
		"email": "alice@example.com",
		"phone": "123-456-7890",
		"address": "123 Main St",
		"city": "New York",
		"country": "USA"
	}`)

	b.ResetTimer()

	for range b.N {
		var (
			person    Person
			originMap map[string]any
		)

		err := JsonUnmarshalToMapAndStruct(jsonData, &person, &originMap)
		require.NoError(b, err)
		require.Equal(b, "Alice", person.Name)
		require.Equal(b, 30, person.Age)
		require.Equal(b, "alice@example.com", person.Email)
		require.Equal(b, "Alice", originMap["name"])
		require.InEpsilon(b, float64(30), originMap["age"], 0)
		require.Equal(b, "alice@example.com", originMap["email"])
	}
}

// BenchmarkUnmarshalToMapAndStructLargePayload 大负载性能测试.
func BenchmarkUnmarshalToMapAndStructLargePayload(b *testing.B) {
	type Item struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	type Data struct {
		Items []Item `json:"items"`
	}

	// 创建包含 100 个项的大 JSON 负载
	items := make([]Item, 100)
	for i := range 100 {
		items[i] = Item{ID: i, Name: "item", Value: "value"}
	}

	data := Data{Items: items}

	jsonData, err := json.Marshal(data)
	if err != nil {
		b.Fatalf("failed to marshal data: %v", err)
	}

	b.ResetTimer()

	for range b.N {
		var (
			result    Data
			originMap map[string]any
		)

		err := JsonUnmarshalToMapAndStruct(jsonData, &result, &originMap)
		require.NoError(b, err)
		require.Len(b, result.Items, 100)
		require.Equal(b, 0, result.Items[0].ID)
		require.Equal(b, "item", result.Items[0].Name)

		itemsArray, ok := originMap["items"].([]any)
		require.True(b, ok)
		require.Len(b, itemsArray, 100)
	}
}

// BenchmarkUnmarshalToMapAndStructDeepNesting 深层嵌套性能测试.
func BenchmarkUnmarshalToMapAndStructDeepNesting(b *testing.B) {
	type Level4 struct {
		Value string `json:"value"`
	}

	type Level3 struct {
		Level4 Level4 `json:"level4"`
	}

	type Level2 struct {
		Level3 Level3 `json:"level3"`
	}

	type Level1 struct {
		Level2 Level2 `json:"level2"`
	}

	jsonData := []byte(`{
		"level2": {
			"level3": {
				"level4": {
					"value": "deep"
				}
			}
		}
	}`)

	b.ResetTimer()

	for range b.N {
		var (
			level1    Level1
			originMap map[string]any
		)

		err := JsonUnmarshalToMapAndStruct(jsonData, &level1, &originMap)
		require.NoError(b, err)
		require.Equal(b, "deep", level1.Level2.Level3.Level4.Value)

		level2Map, ok := originMap["level2"].(map[string]any)
		require.True(b, ok)
		level3Map, ok := level2Map["level3"].(map[string]any)
		require.True(b, ok)
		level4Map, ok := level3Map["level4"].(map[string]any)
		require.True(b, ok)
		require.Equal(b, "deep", level4Map["value"])
	}
}
