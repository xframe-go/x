package event

import (
	"cnb.cool/liey/liey-go/contracts"
)

var (
	instance *EventBusInstance
)

type EventBusInstance struct {
	driver contracts.EventDriver
}

func NewInstance() *EventBusInstance {
	return &EventBusInstance{}
}

func Register(driver contracts.EventDriver) {
	if instance == nil {
		instance = NewInstance()
	}

	instance.driver = driver
}

func Get[T any]() *Bus[T] {
	if instance == nil || instance.driver == nil {
		panic("event bus not initialized")
	}

	return NewBus[T](instance.driver)
}
