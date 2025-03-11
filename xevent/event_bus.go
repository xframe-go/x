package xevent

import (
	"github.com/xframe-go/x/contracts"
	"github.com/xframe-go/x/utils/xmap"
	"golang.org/x/sync/errgroup"
)

type Bus struct {
	listeners *xmap.Map[any, *Subscriber]
}

func New() contracts.EventBus {
	return &Bus{
		listeners: xmap.NewMap[any, *Subscriber](true),
	}
}

func (m *Bus) Listen(symbol any, listeners ...contracts.XListener) {
	subscriber, _ := m.listeners.GetOrSet(symbol, func() *Subscriber {
		return NewSubscriber()
	})

	subscriber.Add(listeners...)
}

func (m *Bus) Emit(symbol any, payload any) error {
	subscriber, ok := m.listeners.Get(symbol)
	if !ok {
		return nil
	}

	listeners := subscriber.All()

	// 处理同步监听器，有错误直接返回
	errs := errgroup.Group{}
	for i := range listeners {
		listener := listeners[i]
		errs.Go(func() error {
			return listener.Handle(payload)
		})
	}

	return errs.Wait()
}
