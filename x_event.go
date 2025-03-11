package x

import (
	"github.com/xframe-go/x/contracts"
	"github.com/xframe-go/x/xevent"
)

func SubscribeEvent[P any](evt contracts.Event[P], listeners ...contracts.Listener[P]) {
	items := make([]contracts.XListener, 0, len(listeners))
	for i := range listeners {
		listener := listeners[i]

		item := xevent.FuncListener(func(payload any) error {
			// todo 断言在symbol冲突的情况下可能出错，需要友好提示
			return listener.Handle(payload.(P))
		})

		if async, ok := listener.(contracts.AsyncListener); ok {
			items = append(items, xevent.NewAsyncFuncListener(item, async.AsyncBuffer()))
		} else {
			items = append(items, item)
		}
	}

	xApp.Get().Event().Listen(evt.Symbol(), items...)
}
