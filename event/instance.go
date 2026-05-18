package event

import (
	"github.com/xframe-go/x/contracts"
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

func Get() *Bus {
	if instance == nil || instance.driver == nil {
		panic("event bus not initialized")
	}

	return NewBus(instance.driver)
}
