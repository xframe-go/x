package event

import (
	"sync"
)

type ChannelDriver struct {
	channels map[string][]chan interface{}
	mu       sync.RWMutex
}

func NewChannelDriver() *ChannelDriver {
	return &ChannelDriver{
		channels: make(map[string][]chan interface{}),
	}
}

func (d *ChannelDriver) Publish(topic string, data interface{}) error {
	d.mu.RLock()
	subscribers, exists := d.channels[topic]
	d.mu.RUnlock()

	if !exists {
		return nil
	}

	for _, ch := range subscribers {
		select {
		case ch <- data:
		default:
			continue
		}
	}

	return nil
}

func (d *ChannelDriver) Subscribe(topic string) (<-chan interface{}, error) {
	ch := make(chan interface{}, 100)

	d.mu.Lock()
	d.channels[topic] = append(d.channels[topic], ch)
	d.mu.Unlock()

	return ch, nil
}

func (d *ChannelDriver) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	for topic, channels := range d.channels {
		for _, ch := range channels {
			close(ch)
		}
		delete(d.channels, topic)
	}

	return nil
}
