package x

import (
	"reflect"
	"strings"

	"github.com/xframe-go/x/event"
)

func RegisterEvent(fn func() event.Config) {
	cfg := fn()

	if cfg.Driver == nil {
		return
	}

	if rocket.bus == nil {
		rocket.bus = event.NewBus(cfg.Driver)
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
	return strings.ToLower(typeName)
}

func Subscribe[M any](topic string, fn func(M)) error {
	if rocket.bus == nil {
		panic("rocket bus is nil")
	}
	return rocket.bus.Subscribe(topic, func(a any) {
		fn(a.(M))
	})
}

func Publish[M any](topic string, data M) error {
	if rocket.bus == nil {
		panic("rocket bus is nil")
	}
	return rocket.bus.Publish(topic, data)
}

func Observer[M any](m M, o event.ModelEventObserver[M]) error {
	topic := extraTopicName(m)

	err := Subscribe[event.ModelCreatedEvent](topic+".created", func(a event.ModelCreatedEvent) {
		if err := o.Created(&event.Observer[M]{
			TX:    a.TX,
			Model: a.Model.(*M),
		}); err != nil {
			Logger().Error(err)
		}
	})
	if err != nil {
		return err
	}

	err = Subscribe[event.ModelUpdatedEvent](topic+".updated", func(a event.ModelUpdatedEvent) {
		origin := a.Origin.(M)
		if err := o.Updated(&event.Observer[M]{
			Model:  a.Model.(*M),
			Origin: &origin,
		}); err != nil {
			Logger().Error(err)
		}
	})
	if err != nil {
		return err
	}

	err = Subscribe(topic+".deleted", func(a M) {
		if err := o.Deleted(&event.Observer[M]{
			Model: &a,
		}); err != nil {
			Logger().Error(err)
		}
	})
	if err != nil {
		return err
	}
	return nil
}
