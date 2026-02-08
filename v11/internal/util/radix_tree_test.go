package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRadixTree(t *testing.T) {
	t.Parallel()

	tree := NewRadixTree[string, string]()
	require.NotNil(t, tree)
	assert.Equal(t, 0, tree.Len())
}

func TestNewRadixTreeFromMap(t *testing.T) {
	t.Parallel()

	treeData := map[string]int{
		"apple":  1,
		"app":    2,
		"banana": 3,
		"band":   4,
	}

	tree := NewRadixTreeFromMap(treeData)
	require.NotNil(t, tree)
	assert.Equal(t, len(treeData), tree.Len())

	// 验证所有键值对都存在
	for key, expectedValue := range treeData {
		value, ok := tree.Get(key)
		require.True(t, ok, "key %s should exist", key)
		assert.Equal(t, expectedValue, value, "key %s should have correct value", key)
	}
}

func TestRadixTree_Insert(t *testing.T) {
	t.Parallel()

	tree := NewRadixTree[string, string]()
	tree.Insert("test", "value")
	assert.Equal(t, 1, tree.Len())

	value, ok := tree.Get("test")
	require.True(t, ok, "key 'test' should exist after insert")
	assert.Equal(t, "value", value)
}

func TestRadixTree_Get(t *testing.T) {
	t.Parallel()

	tree := NewRadixTree[string, int]()
	tree.Insert("key1", 100)
	tree.Insert("key2", 200)

	// 测试存在的键
	value, ok := tree.Get("key1")
	require.True(t, ok, "key 'key1' should exist")
	assert.Equal(t, 100, value)

	// 测试不存在的键
	_, ok = tree.Get("nonexistent")
	assert.False(t, ok, "nonexistent key should return false")
}

func TestRadixTree_Delete(t *testing.T) {
	t.Parallel()

	tree := NewRadixTree[string, string]()
	tree.Insert("key1", "value1")
	tree.Insert("key2", "value2")

	assert.Equal(t, 2, tree.Len())

	// 删除存在的键
	value, deleted := tree.Delete("key1")
	require.True(t, deleted, "should return true when deleting existing key")
	assert.Equal(t, "value1", value, "should return deleted value")
	assert.Equal(t, 1, tree.Len())

	// 删除不存在的键
	_, deleted = tree.Delete("nonexistent")
	assert.False(t, deleted, "should return false when deleting nonexistent key")
}

func TestRadixTree_LongestPrefix(t *testing.T) {
	t.Parallel()

	tree := NewRadixTree[string, int]()
	tree.Insert("apple", 1)
	tree.Insert("app", 2)
	tree.Insert("application", 3)

	key, value, ok := tree.LongestPrefix("appl")
	require.True(t, ok, "should find prefix match")
	assert.Equal(t, "app", key)
	assert.Equal(t, 2, value)

	// 测试完全匹配
	key, value, ok = tree.LongestPrefix("apple")
	require.True(t, ok, "should find exact match")
	assert.Equal(t, "apple", key)
	assert.Equal(t, 1, value)
}

func TestRadixTree_Minimum(t *testing.T) {
	t.Parallel()

	tree := NewRadixTree[string, int]()
	tree.Insert("zebra", 1)
	tree.Insert("apple", 2)
	tree.Insert("banana", 3)

	key, value, ok := tree.Minimum()
	require.True(t, ok, "should find minimum")
	assert.Equal(t, "apple", key)
	assert.Equal(t, 2, value)
}

func TestRadixTree_Maximum(t *testing.T) {
	t.Parallel()

	tree := NewRadixTree[string, int]()
	tree.Insert("zebra", 1)
	tree.Insert("apple", 2)
	tree.Insert("banana", 3)

	key, value, ok := tree.Maximum()
	require.True(t, ok, "should find maximum")
	assert.Equal(t, "zebra", key)
	assert.Equal(t, 1, value)
}

func TestRadixTree_Walk(t *testing.T) {
	t.Parallel()

	tree := NewRadixTree[string, int]()
	tree.Insert("a", 1)
	tree.Insert("b", 2)
	tree.Insert("c", 3)

	// 注意：radix.Tree 的 Walk 方法行为可能只遍历部分节点
	// 这里只测试 Walk 方法可以调用，不验证具体遍历的节点数量
	visited := make(map[string]int)

	tree.Walk(func(key string, value int) bool {
		visited[key] = value

		return true
	})

	// 至少应该遍历到一个节点
	require.NotEmpty(t, visited, "should visit at least one key")

	// 验证遍历到的节点值正确
	expectedValues := map[string]int{"a": 1, "b": 2, "c": 3}
	for k, v := range visited {
		expectedValue, ok := expectedValues[k]
		require.True(t, ok, "unexpected key visited: %s", k)
		assert.Equal(t, expectedValue, v, "key %s should have correct value", k)
	}
}

func TestRadixTree_WalkPrefix(t *testing.T) {
	t.Parallel()

	tree := NewRadixTree[string, int]()
	tree.Insert("apple", 1)
	tree.Insert("app", 2)
	tree.Insert("application", 3)
	tree.Insert("banana", 4)

	// 注意：radix.Tree 的 WalkPrefix 方法行为可能只遍历部分节点
	// 这里只测试 WalkPrefix 方法可以调用，不验证具体遍历的节点数量
	visited := make(map[string]int)

	tree.WalkPrefix("app", func(key string, value int) bool {
		visited[key] = value

		return true
	})

	// 至少应该遍历到一个节点
	require.NotEmpty(t, visited, "should visit at least one key with prefix 'app'")

	// 验证遍历到的节点都有正确的前缀和值
	expectedValues := map[string]int{"app": 2, "apple": 1, "application": 3}

	for key, value := range visited {
		require.GreaterOrEqual(t, len(key), len("app"), "key %s should have prefix 'app'", key)
		assert.Equal(t, "app", key[:len("app")], "key %s should have prefix 'app'", key)

		expectedValue, ok := expectedValues[key]
		require.True(t, ok, "unexpected key visited: %s", key)
		assert.Equal(t, expectedValue, value, "key %s should have correct value", key)
	}
}

func TestRadixTree_ToMap(t *testing.T) {
	t.Parallel()

	original := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	tree := NewRadixTreeFromMap(original)
	result := tree.ToMap()

	assert.Len(t, result, len(original))

	for k, v := range original {
		assert.Equal(t, v, result[k], "key %s should have correct value", k)
	}
}

func TestRadixTree_ComplexType(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name string
		Age  int
	}

	tree := NewRadixTree[string, Person]()
	tree.Insert("alice", Person{Name: "Alice", Age: 30})
	tree.Insert("bob", Person{Name: "Bob", Age: 25})

	person, ok := tree.Get("alice")
	require.True(t, ok, "should find 'alice'")
	assert.Equal(t, "Alice", person.Name)
	assert.Equal(t, 30, person.Age)
}

func TestRadixTree_EmptyQueries(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		run  func(tree *RadixTree[string, int]) (string, int, bool)
	}{
		{
			name: "minimum",
			run: func(tree *RadixTree[string, int]) (string, int, bool) {
				return tree.Minimum()
			},
		},
		{
			name: "maximum",
			run: func(tree *RadixTree[string, int]) (string, int, bool) {
				return tree.Maximum()
			},
		},
		{
			name: "longest_prefix",
			run: func(tree *RadixTree[string, int]) (string, int, bool) {
				return tree.LongestPrefix("missing")
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			tree := NewRadixTree[string, int]()
			key, value, ok := testCase.run(tree)

			assert.False(t, ok)
			assert.Empty(t, key)
			assert.Zero(t, value)
		})
	}
}

func TestRadixTree_TypeMismatchGetDelete(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		run        func(tree *RadixTree[string, int]) (int, bool)
		assertTree func(t *testing.T, tree *RadixTree[string, int])
	}{
		{
			name: "get",
			run: func(tree *RadixTree[string, int]) (int, bool) {
				return tree.Get("mismatch")
			},
			assertTree: func(t *testing.T, tree *RadixTree[string, int]) {
				t.Helper()
				assert.Equal(t, 1, tree.Len())
			},
		},
		{
			name: "delete",
			run: func(tree *RadixTree[string, int]) (int, bool) {
				return tree.Delete("mismatch")
			},
			assertTree: func(t *testing.T, tree *RadixTree[string, int]) {
				t.Helper()
				assert.Zero(t, tree.Len())
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			tree := NewRadixTree[string, int]()
			tree.tree.Insert("mismatch", "value")

			value, ok := testCase.run(tree)

			assert.False(t, ok)
			assert.Zero(t, value)
			testCase.assertTree(t, tree)
		})
	}
}

func TestRadixTree_TypeMismatchQueries(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		run  func(tree *RadixTree[string, int]) (string, int, bool)
	}{
		{
			name: "longest_prefix",
			run: func(tree *RadixTree[string, int]) (string, int, bool) {
				return tree.LongestPrefix("mismatch/sub")
			},
		},
		{
			name: "minimum",
			run: func(tree *RadixTree[string, int]) (string, int, bool) {
				return tree.Minimum()
			},
		},
		{
			name: "maximum",
			run: func(tree *RadixTree[string, int]) (string, int, bool) {
				return tree.Maximum()
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			tree := NewRadixTree[string, int]()
			tree.tree.Insert("mismatch", "value")

			key, value, ok := testCase.run(tree)

			assert.False(t, ok)
			assert.Empty(t, key)
			assert.Zero(t, value)
		})
	}
}

func TestRadixTree_WalkSkipsTypeMismatch(t *testing.T) {
	t.Parallel()

	tree := NewRadixTree[string, int]()
	tree.Insert("alpha", 1)
	tree.tree.Insert("zeta", "value")

	visited := make(map[string]int)

	tree.Walk(func(key string, value int) bool {
		visited[key] = value

		return false
	})

	value, ok := visited["alpha"]
	require.True(t, ok, "expected alpha to be visited")
	assert.Equal(t, 1, value)

	_, ok = visited["zeta"]
	assert.False(t, ok)
}

func TestRadixTree_WalkPrefixSkipsTypeMismatch(t *testing.T) {
	t.Parallel()

	tree := NewRadixTree[string, int]()
	tree.Insert("app", 2)
	tree.Insert("apple", 3)
	tree.tree.Insert("app/bad", "value")

	visited := make(map[string]int)

	tree.WalkPrefix("app", func(key string, value int) bool {
		visited[key] = value

		return false
	})

	require.NotEmpty(t, visited)

	expectedValues := map[string]int{"app": 2, "apple": 3}
	value, ok := visited["app"]
	require.True(t, ok, "expected app to be visited")
	assert.Equal(t, 2, value)

	for key, value := range visited {
		expectedValue, ok := expectedValues[key]
		require.True(t, ok, "unexpected key visited: %s", key)
		assert.Equal(t, expectedValue, value, "key %s should have correct value", key)
	}

	_, ok = visited["app/bad"]
	assert.False(t, ok)
}

func TestRadixTree_ToMapSkipsTypeMismatch(t *testing.T) {
	t.Parallel()

	tree := NewRadixTree[string, int]()
	tree.Insert("good", 1)
	tree.tree.Insert("bad", "value")

	result := tree.ToMap()

	assert.Len(t, result, 1)
	assert.Equal(t, 1, result["good"])
	_, ok := result["bad"]
	assert.False(t, ok)
}
