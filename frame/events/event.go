package events

import "github.com/xframe-go/x"

type Event[P any] struct {
	symbol any
}

func NewEvent[P any](symbol any) *Event[P] {
	return &Event[P]{
		symbol: symbol,
	}
}

func (e Event[P]) Emit(payload P) error {
	return x.App().Event().Emit(e.symbol, payload)
}

func (e Event[P]) Symbol() any {
	return e.symbol
}
