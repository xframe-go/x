package x

import (
	"reflect"
	"sync"

	"github.com/xframe-go/x/event"
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
}

func extraTopicName(model interface{}) string {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	typeName := "liey." + modelType.Name()
	return typeName
}

var Events = sync.OnceValue(func() *event.Bus[any] {
	return event.NewBus[any](event.NewChannelDriver())
})

func Subscribe[M any](topic string, fn func(M)) error {
	return Events().Subscribe(topic, func(a any) {
		fn(a.(M))
	})
}

func Publish[M any](topic string, data M) error {
	return Events().Publish(topic, data)
}

func Observer[M any](m M, o event.ModelEventObserver[M]) error {
	topic := extraTopicName(m)

	err := Subscribe(topic+".created", func(a M) {
		if err := o.Created(a); err != nil {
			Logger().Error(err)
		}
	})
	if err != nil {
		return err
	}

	err = Subscribe(topic+".updated", func(a M) {
		if err := o.Updated(a); err != nil {
			Logger().Error(err)
		}
	})
	if err != nil {
		return err
	}

	err = Subscribe(topic+".deleted", func(a M) {
		if err := o.Deleted(a); err != nil {
			Logger().Error(err)
		}
	})
	if err != nil {
		return err
	}
	return nil
}
