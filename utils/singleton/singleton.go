package singleton

import "sync"

type Singleton[V any] struct {
	value   V
	once    sync.Once
	creator func() V
}

func New[V any](creator func() V) *Singleton[V] {
	return &Singleton[V]{
		once:    sync.Once{},
		creator: creator,
	}
}

func (s *Singleton[V]) Get() V {
	s.once.Do(func() {
		s.value = s.creator()
	})
	return s.value
}

func (s *Singleton[V]) GetOrSet(creator func() V) V {
	s.once.Do(func() {
		s.creator = creator
		s.value = creator()
	})
	return s.value
}
