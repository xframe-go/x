package event

import (
	"context"
	"sync"

	"cnb.cool/liey/liey-go/contracts"
)

type Bus[T any] struct {
	driver contracts.EventDriver
	subs   map[string][]func(T)
	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

func NewBus[T any](driver contracts.EventDriver) *Bus[T] {
	ctx, cancel := context.WithCancel(context.Background())

	bus := &Bus[T]{
		driver: driver,
		subs:   make(map[string][]func(T)),
		ctx:    ctx,
		cancel: cancel,
	}

	return bus
}

func (b *Bus[T]) Subscribe(topic string, handler func(T)) error {
	b.mu.Lock()
	b.subs[topic] = append(b.subs[topic], handler)
	b.mu.Unlock()

	ch, err := b.driver.Subscribe(topic)
	if err != nil {
		return err
	}

	go b.consume(topic, ch)
	return nil
}

func (b *Bus[T]) Publish(topic string, data T) error {
	return b.driver.Publish(topic, data)
}

func (b *Bus[T]) Close() error {
	b.cancel()
	return b.driver.Close()
}

func (b *Bus[T]) consume(topic string, ch <-chan interface{}) {
	for {
		select {
		case <-b.ctx.Done():
			return
		case data, ok := <-ch:
			if !ok {
				return
			}

			b.mu.RLock()
			handlers := b.subs[topic]
			b.mu.RUnlock()

			for _, handler := range handlers {
				if event, ok := data.(T); ok {
					handler(event)
				}
			}
		}
	}
}
