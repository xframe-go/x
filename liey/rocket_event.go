package liey

import (
	"cnb.cool/liey/liey-go/event"
)

func RegisterEvent(fn func() event.Config) {
	cfg := fn()

	if cfg.Driver == nil {
		return
	}

	if rocket.bus == nil {
		rocket.bus = event.NewBus[any](cfg.Driver)
	}

	for _, handler := range cfg.Handlers {
		topic := handler.Event()
		if topic == "" {
			continue
		}
		err := rocket.bus.Subscribe(topic, func(data any) {
			handler.Handle(data)
		})
		if err != nil {
			panic(err)
		}
	}

	plugin := event.NewPlugin(rocket.bus, event.GormPluginConfig{
		PublishCreated: true,
		PublishUpdated: true,
		PublishDeleted: true,
		Prefix:         "wanderlust",
	})

	if err := DB().Use(plugin); err != nil {
		panic(err)
	}
}

func Event[T any]() *event.Bus[T] {
	if rocket.bus == nil {
		panic("event bus not initialized")
	}

	var bus any = rocket.bus

	return bus.(*event.Bus[T])
}
