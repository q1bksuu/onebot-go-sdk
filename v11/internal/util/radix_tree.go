package util

import (
	"github.com/armon/go-radix"
)

// RadixTree 是对 radix.Tree 的泛型封装，提供类型安全的操作接口
// K 是键类型，约束为底层类型为 string 的类型
// T 是值类型.
type RadixTree[K ~string, T any] struct {
	tree *radix.Tree
}

type RadixTreeStrKey[T any] = RadixTree[string, T]

// NewRadixTree 创建一个新的空 RadixTree.
func NewRadixTree[K ~string, T any]() *RadixTree[K, T] {
	return &RadixTree[K, T]{
		tree: radix.New(),
	}
}

// NewRadixTreeFromMap 从 map 初始化 RadixTree
// 将 map 中的键值对插入到 RadixTree 中.
func NewRadixTreeFromMap[K ~string, T any](m map[K]T) *RadixTree[K, T] {
	tree := radix.New()
	for key, value := range m {
		tree.Insert(string(key), value)
	}

	return &RadixTree[K, T]{
		tree: tree,
	}
}

// Insert 插入或更新键值对.
func (rt *RadixTree[K, T]) Insert(key K, value T) {
	rt.tree.Insert(string(key), value)
}

// Get 根据键获取值，返回值和是否存在.
func (rt *RadixTree[K, T]) Get(key K) (T, bool) {
	value, ok := rt.tree.Get(string(key))
	if !ok {
		var zero T

		return zero, false
	}

	typedValue, ok := value.(T)
	if !ok {
		var zero T

		return zero, false
	}

	return typedValue, true
}

// Delete 删除指定键，返回被删除的值和是否删除成功.
func (rt *RadixTree[K, T]) Delete(key K) (T, bool) {
	value, ok := rt.tree.Delete(string(key))
	if !ok {
		var zero T

		return zero, false
	}

	typedValue, ok := value.(T)
	if !ok {
		var zero T

		return zero, false
	}

	return typedValue, true
}

// LongestPrefix 查找最长前缀匹配的键值对.
func (rt *RadixTree[K, T]) LongestPrefix(prefix K) (K, T, bool) {
	key, value, ok := rt.tree.LongestPrefix(string(prefix))
	if !ok {
		var zero T

		return "", zero, false
	}

	typedValue, ok := value.(T)
	if !ok {
		var zero T

		return "", zero, false
	}

	return K(key), typedValue, true
}

// Minimum 返回树中的最小键值对.
func (rt *RadixTree[K, T]) Minimum() (K, T, bool) {
	key, value, ok := rt.tree.Minimum()
	if !ok {
		var zero T

		return "", zero, false
	}

	typedValue, ok := value.(T)
	if !ok {
		var zero T

		return "", zero, false
	}

	return K(key), typedValue, true
}

// Maximum 返回树中的最大键值对.
func (rt *RadixTree[K, T]) Maximum() (K, T, bool) {
	key, value, ok := rt.tree.Maximum()
	if !ok {
		var zero T

		return "", zero, false
	}

	typedValue, ok := value.(T)
	if !ok {
		var zero T

		return "", zero, false
	}

	return K(key), typedValue, true
}

// Walk 遍历树，对每个键值对执行回调函数
// 如果回调函数返回 false，则停止遍历.
func (rt *RadixTree[K, T]) Walk(fn func(key K, value T) bool) {
	rt.tree.Walk(func(key string, value any) bool {
		typedValue, ok := value.(T)
		if !ok {
			return true
		}

		return fn(K(key), typedValue)
	})
}

// WalkPrefix 遍历具有指定前缀的所有键值对.
func (rt *RadixTree[K, T]) WalkPrefix(prefix K, fn func(key K, value T) bool) {
	rt.tree.WalkPrefix(string(prefix), func(key string, value any) bool {
		typedValue, ok := value.(T)
		if !ok {
			return true
		}

		return fn(K(key), typedValue)
	})
}

// Len 返回树中键值对的数量.
func (rt *RadixTree[K, T]) Len() int {
	return rt.tree.Len()
}

// ToMap 将 RadixTree 转换为 map.
func (rt *RadixTree[K, T]) ToMap() map[K]T {
	m := rt.tree.ToMap()

	result := make(map[K]T, len(m))
	for key, value := range m {
		typedValue, ok := value.(T)
		if !ok {
			continue
		}

		result[K(key)] = typedValue
	}

	return result
}
