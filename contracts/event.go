package contracts

type EventDriver interface {
	Publish(topic string, data interface{}) error
	Subscribe(topic string) (<-chan interface{}, error)
	Close() error
}
