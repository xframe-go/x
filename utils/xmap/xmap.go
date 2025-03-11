package xmap

import "github.com/xframe-go/x/internal/rwmutex"

// Package xmap 提供了一个支持并发安全的泛型 Map 实现
// 该实现支持基本的 Map 操作，并可选择是否启用并发安全特性

// Map 是一个支持并发安全的泛型 Map 实现
// K 必须是可比较的类型
// V 可以是任意类型
// 示例:
//
//	m := NewMap[string, int]()      // 非并发安全
//	m := NewMap[string, int](true)  // 并发安全
type Map[K comparable, V any] struct {
	data map[K]V
	mu   rwmutex.RWMutex
}

// NewMap 创建一个新的 Map 实例
// safe 参数用于控制是否启用并发安全，默认为 false
func NewMap[K comparable, V any](safe ...bool) *Map[K, V] {
	return &Map[K, V]{
		data: make(map[K]V),
		mu:   rwmutex.Create(safe...),
	}
}

// Set 设置键值对
// 如果键已存在，则更新对应的值
// 该操作在并发安全模式下是线程安全的
func (m *Map[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

// Get 获取指定键的值
// 返回值:
//   - value: 键对应的值
//   - ok: 如果键存在则为 true，否则为 false
func (m *Map[K, V]) Get(key K) (value V, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok = m.data[key]
	return
}

// Delete 删除指定键的值
// 如果键不存在，该操作不会产生任何效果
func (m *Map[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

// Len 返回 Map 中键值对的数量
func (m *Map[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}

// Clear 清空 Map 中的所有键值对
func (m *Map[K, V]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = make(map[K]V)
}

// Keys 返回 Map 中所有的键
func (m *Map[K, V]) Keys() []K {
	m.mu.RLock()
	defer m.mu.RUnlock()
	keys := make([]K, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys
}

// Values 返回 Map 中所有的值
func (m *Map[K, V]) Values() []V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	values := make([]V, 0, len(m.data))
	for _, v := range m.data {
		values = append(values, v)
	}
	return values
}

// Range 遍历 Map 中的所有键值对
// 参数 fn 是一个回调函数，接收 key 和 value 作为参数
// 如果 fn 返回 false，则停止遍历
// 注意: 在遍历过程中，Map 是被读锁定的，请避免在回调函数中执行耗时操作
func (m *Map[K, V]) Range(fn func(key K, value V) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.data {
		if !fn(k, v) {
			break
		}
	}
}

// GetOrSet 获取键对应的值，如果键不存在则通过 valueFn 生成默认值并设置
// 返回值:
//   - value: 键对应的值（已存在的值或新设置的默认值）
//   - loaded: 如果键已存在则为 true，否则为 false
func (m *Map[K, V]) GetOrSet(key K, valueFn func() V) (value V, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	value, loaded = m.data[key]
	if !loaded {
		value = valueFn()
		m.data[key] = value
	}
	return
}
