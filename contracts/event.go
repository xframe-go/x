package contracts

type EventBus interface {
	Listen(symbol any, listeners ...XListener)
	Emit(symbol any, payload any) error
}

type Listener[P any] interface {
	Handle(evt P) error
}

type AsyncListener interface {
	AsyncBuffer() int
}

type Event[P any] interface {
	Emit(payload P) error
	Symbol() any
}

type XEvent interface {
	Symbol() any
	Payload() any
}

type XListener interface {
	Handle(payload any) error
}
