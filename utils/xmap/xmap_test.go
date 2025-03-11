package xmap

import (
	"fmt"
	"sync"
	"testing"
)

func TestMap(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		m := NewMap[string, int]()

		// 测试 Set 和 Get
		m.Set("one", 1)
		if val, ok := m.Get("one"); !ok || val != 1 {
			t.Errorf("expected Get(\"one\") = (1, true), got (%v, %v)", val, ok)
		}

		// 测试不存在的键
		if _, ok := m.Get("nonexistent"); ok {
			t.Error("expected Get(\"nonexistent\") to return false")
		}

		// 测试 Len
		if l := m.Len(); l != 1 {
			t.Errorf("expected Len() = 1, got %d", l)
		}

		// 测试 Delete
		m.Delete("one")
		if _, ok := m.Get("one"); ok {
			t.Error("expected key \"one\" to be deleted")
		}
	})

	t.Run("concurrent safety", func(t *testing.T) {
		m := NewMap[int, int](true) // 启用并发安全
		var wg sync.WaitGroup
		n := 1000

		// 并发写入
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(val int) {
				defer wg.Done()
				m.Set(val, val*2)
			}(i)
		}
		wg.Wait()

		// 验证结果
		if m.Len() != n {
			t.Errorf("expected length %d, got %d", n, m.Len())
		}

		// 测试 Range
		count := 0
		m.Range(func(key, value int) bool {
			if value != key*2 {
				t.Errorf("expected value %d for key %d, got %d", key*2, key, value)
			}
			count++
			return true
		})
		if count != n {
			t.Errorf("Range visited %d elements, expected %d", count, n)
		}
	})

	t.Run("keys and values", func(t *testing.T) {
		m := NewMap[string, int]()
		testData := map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		}

		// 添加测试数据
		for k, v := range testData {
			m.Set(k, v)
		}

		// 测试 Keys
		keys := m.Keys()
		if len(keys) != len(testData) {
			t.Errorf("expected %d keys, got %d", len(testData), len(keys))
		}
		for _, k := range keys {
			if _, ok := testData[k]; !ok {
				t.Errorf("unexpected key: %v", k)
			}
		}

		// 测试 Values
		values := m.Values()
		if len(values) != len(testData) {
			t.Errorf("expected %d values, got %d", len(testData), len(values))
		}

		// 测试 Clear
		m.Clear()
		if m.Len() != 0 {
			t.Errorf("expected empty map after Clear(), got length %d", m.Len())
		}
	})

	t.Run("concurrent operations", func(t *testing.T) {
		m := NewMap[int, int](true)
		var wg sync.WaitGroup
		n := 1000

		// 并发写入和读取
		for i := 0; i < n; i++ {
			wg.Add(2)
			// 写入 goroutine
			go func(val int) {
				defer wg.Done()
				m.Set(val, val*2)
			}(i)
			// 读取 goroutine
			go func(val int) {
				defer wg.Done()
				// 反复读取直到写入完成
				for {
					if v, ok := m.Get(val); ok && v == val*2 {
						break
					}
				}
			}(i)
		}
		wg.Wait()

		if m.Len() != n {
			t.Errorf("expected length %d, got %d", n, m.Len())
		}
	})

	t.Run("concurrent mixed operations", func(t *testing.T) {
		m := NewMap[int, string](true)
		var wg sync.WaitGroup
		n := 100

		// 第一阶段：并发设置和获取
		for i := 0; i < n; i++ {
			wg.Add(2)
			// Set
			go func(val int) {
				defer wg.Done()
				m.Set(val, fmt.Sprintf("value-%d", val))
			}(i)
			// Get
			go func(val int) {
				defer wg.Done()
				_, _ = m.Get(val)
			}(i)
		}
		wg.Wait()

		// 第二阶段：删除偶数键
		for i := 0; i < n; i++ {
			if i%2 == 0 {
				m.Delete(i)
			}
		}

		// 第三阶段：并发 Range 操作
		var wg2 sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg2.Add(1)
			go func() {
				defer wg2.Done()
				m.Range(func(key int, value string) bool {
					// 验证所有遍历到的键都是奇数
					if key%2 == 0 {
						t.Errorf("found unexpected even key: %d", key)
					}
					return true
				})
			}()
		}
		wg2.Wait()

		// 最终验证
		for i := 0; i < n; i++ {
			val, ok := m.Get(i)
			if i%2 == 0 {
				if ok {
					t.Errorf("expected key %d to be deleted, but got value: %s", i, val)
				}
			} else {
				if !ok {
					t.Errorf("expected key %d to exist", i)
				}
				expectedVal := fmt.Sprintf("value-%d", i)
				if val != expectedVal {
					t.Errorf("expected value %q for key %d, got %q", expectedVal, i, val)
				}
			}
		}
	})

	t.Run("edge cases", func(t *testing.T) {
		m := NewMap[string, interface{}]()

		// 测试 nil 值
		m.Set("nil-value", nil)
		if val, ok := m.Get("nil-value"); !ok || val != nil {
			t.Error("failed to handle nil value correctly")
		}

		// 测试空字符串键
		m.Set("", "empty-key")
		if val, ok := m.Get(""); !ok || val != "empty-key" {
			t.Error("failed to handle empty string key")
		}

		// 测试大量数据
		for i := 0; i < 10000; i++ {
			m.Set(fmt.Sprintf("key-%d", i), i)
		}
		if m.Len() != 10002 { // 10000 + 之前的2个键值对
			t.Errorf("unexpected length after bulk insert: %d", m.Len())
		}

		// 测试 Clear 后的操作
		m.Clear()
		m.Set("new-key", "new-value")
		if m.Len() != 1 {
			t.Error("map should work correctly after Clear()")
		}
	})

	t.Run("type constraints", func(t *testing.T) {
		// 测试不同类型的 Map
		type customKey struct {
			id int
		}
		type customValue struct {
			data string
		}

		// 这些代码应该能够编译
		_ = NewMap[int, int]()
		_ = NewMap[string, interface{}]()
		_ = NewMap[customKey, customValue]() // 自定义类型
		_ = NewMap[int, []string]()          // 切片作为值
		_ = NewMap[string, map[string]int]() // map 作为值
	})

	t.Run("GetOrSet operations", func(t *testing.T) {
		m := NewMap[string, int]()

		// 测试不存在的键
		val, loaded := m.GetOrSet("key1", func() int { return 100 })
		if loaded {
			t.Error("expected key1 to not exist")
		}
		if val != 100 {
			t.Errorf("expected value 100, got %d", val)
		}

		// 测试已存在的键（确保不会调用 valueFn）
		called := false
		val, loaded = m.GetOrSet("key1", func() int {
			called = true
			return 200
		})
		if !loaded {
			t.Error("expected key1 to exist")
		}
		if val != 100 {
			t.Errorf("expected value 100, got %d", val)
		}
		if called {
			t.Error("valueFn should not have been called for existing key")
		}

		// 测试并发 GetOrSet
		m = NewMap[string, int](true)
		var wg sync.WaitGroup
		n := 100
		callCount := 0
		var mu sync.Mutex

		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				val, _ := m.GetOrSet("concurrent-key", func() int {
					mu.Lock()
					callCount++
					mu.Unlock()
					return 1
				})
				if val != 1 {
					t.Errorf("expected value 1, got %d", val)
				}
			}()
		}
		wg.Wait()

		// 验证最终只调用了一次 valueFn
		if callCount != 1 {
			t.Errorf("expected valueFn to be called exactly once, got %d calls", callCount)
		}

		// 验证值被正确设置
		val, loaded = m.Get("concurrent-key")
		if !loaded {
			t.Error("expected concurrent-key to exist")
		}
		if val != 1 {
			t.Errorf("expected value 1, got %d", val)
		}
	})
}
 
// BenchmarkMap 提供性能基准测试
func BenchmarkMap(b *testing.B) {
	b.Run("concurrent set", func(b *testing.B) {
		m := NewMap[int, int](true)
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				m.Set(i, i)
				i++
			}
		})
	})

	b.Run("concurrent get", func(b *testing.B) {
		m := NewMap[int, int](true)
		for i := 0; i < b.N; i++ {
			m.Set(i, i)
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				_, _ = m.Get(i % b.N)
				i++
			}
		})
	})
}
