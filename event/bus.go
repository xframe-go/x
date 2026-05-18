package event

import (
	"context"
	"sync"

	"github.com/xframe-go/x/contracts"
)

type Bus struct {
	driver contracts.EventDriver
	subs   map[string][]func(any)
	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

func NewBus(driver contracts.EventDriver) *Bus {
	ctx, cancel := context.WithCancel(context.Background())

	bus := &Bus{
		driver: driver,
		subs:   make(map[string][]func(any)),
		ctx:    ctx,
		cancel: cancel,
	}

	return bus
}

func (b *Bus) Subscribe(topic string, handler func(any)) error {
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

func (b *Bus) Publish(topic string, data any) error {
	return b.driver.Publish(topic, data)
}

func (b *Bus) Close() error {
	b.cancel()
	return b.driver.Close()
}

func (b *Bus) consume(topic string, ch <-chan interface{}) {
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
				handler(data)
			}
		}
	}
}
