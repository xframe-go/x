package event

import "github.com/xframe-go/x/contracts"

type Config struct {
	Driver    contracts.EventDriver
	Handlers  []Handler
	Observers []IObserver[any]
}

type IObserver[M any] struct {
	Model    M
	Observer ModelEventObserver[M]
}

type ModelEventObserver[M any] interface {
	Created(m *Observer[M]) error
	Updated(m *Observer[M]) error
	Deleted(m *Observer[M]) error
}
