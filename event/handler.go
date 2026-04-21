package event

type Handler interface {
	Event() string
	Handle(data any)
}
