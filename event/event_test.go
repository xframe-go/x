package event_test

import (
	"sync"
	"testing"
	"time"

	"cnb.cool/liey/liey-go/event"
)

type TestEvent struct {
	Message string
}

type MockDriver struct {
	channels map[string]chan interface{}
}

func NewMockDriver() *MockDriver {
	return &MockDriver{
		channels: make(map[string]chan interface{}),
	}
}

func (d *MockDriver) Publish(topic string, data interface{}) error {
	if ch, exists := d.channels[topic]; exists {
		ch <- data
	}
	return nil
}

func (d *MockDriver) Subscribe(topic string) (<-chan interface{}, error) {
	ch := make(chan interface{}, 100)
	d.channels[topic] = ch
	return ch, nil
}

func (d *MockDriver) Close() error {
	for _, ch := range d.channels {
		close(ch)
	}
	return nil
}

func TestChannelDriver_PublishSubscribe(t *testing.T) {
	driver := event.NewChannelDriver()
	defer driver.Close()

	ch, err := driver.Subscribe("test-topic")
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	testData := TestEvent{Message: "hello"}
	go func() {
		if err := driver.Publish("test-topic", testData); err != nil {
			t.Errorf("Publish failed: %v", err)
		}
	}()

	select {
	case data := <-ch:
		if event, ok := data.(TestEvent); ok {
			if event.Message != "hello" {
				t.Errorf("Expected message 'hello', got '%s'", event.Message)
			}
		} else {
			t.Error("Data is not TestEvent type")
		}
	case <-time.After(1 * time.Second):
		t.Error("Timeout waiting for event")
	}
}

func TestChannelDriver_MultipleSubscribers(t *testing.T) {
	driver := event.NewChannelDriver()
	defer driver.Close()

	count := 0
	var mu sync.Mutex

	for i := 0; i < 3; i++ {
		ch, err := driver.Subscribe("test-topic")
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		go func(ch <-chan interface{}) {
			data := <-ch
			mu.Lock()
			count++
			mu.Unlock()

			if event, ok := data.(TestEvent); ok {
				if event.Message != "hello" {
					t.Errorf("Expected message 'hello', got '%s'", event.Message)
				}
			}
		}(ch)
	}

	driver.Publish("test-topic", TestEvent{Message: "hello"})

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	if count != 3 {
		t.Errorf("Expected 3 subscribers to receive event, got %d", count)
	}
	mu.Unlock()
}

func TestEventBus_SubscribePublish(t *testing.T) {
	driver := NewMockDriver()
	defer driver.Close()

	bus := event.NewBus[TestEvent](driver)

	received := false
	err := bus.Subscribe("test-topic", func(event TestEvent) {
		if event.Message == "hello" {
			received = true
		}
	})
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	bus.Publish("test-topic", TestEvent{Message: "hello"})

	time.Sleep(100 * time.Millisecond)

	if !received {
		t.Error("Event not received")
	}
}

func TestEventBus_MultipleHandlers(t *testing.T) {
	driver := NewMockDriver()
	defer driver.Close()

	bus := event.NewBus[TestEvent](driver)

	count := 0

	for i := 0; i < 3; i++ {
		if err := bus.Subscribe("test-topic", func(event TestEvent) {
			if event.Message == "hello" {
				count++
			}
		}); err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}
	}

	bus.Publish("test-topic", TestEvent{Message: "hello"})

	time.Sleep(100 * time.Millisecond)

	if count != 3 {
		t.Errorf("Expected 3 handlers to be called, got %d", count)
	}
}

func TestEventBus_DifferentTopics(t *testing.T) {
	driver := NewMockDriver()
	defer driver.Close()

	bus := event.NewBus[TestEvent](driver)

	topic1Received := false
	topic2Received := false

	bus.Subscribe("topic1", func(event TestEvent) {
		if event.Message == "topic1" {
			topic1Received = true
		}
	})

	bus.Subscribe("topic2", func(event TestEvent) {
		if event.Message == "topic2" {
			topic2Received = true
		}
	})

	bus.Publish("topic1", TestEvent{Message: "topic1"})
	bus.Publish("topic2", TestEvent{Message: "topic2"})

	time.Sleep(100 * time.Millisecond)

	if !topic1Received {
		t.Error("Topic1 event not received")
	}
	if !topic2Received {
		t.Error("Topic2 event not received")
	}
}
