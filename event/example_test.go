package event

import (
	"fmt"
	"time"
)

type UserCreatedEvent struct {
	UserID   uint64
	Username string
	Email    string
}

type PostPublishedEvent struct {
	PostID    uint64
	Title     string
	AuthorID  uint64
	Published bool
}

func ExampleChannelDriver() {
	driver := NewChannelDriver()
	defer driver.Close()

	ch, _ := driver.Subscribe("user-events")
	go func() {
		for data := range ch {
			fmt.Printf("Received event: %v\n", data)
		}
	}()

	driver.Publish("user-events", UserCreatedEvent{
		UserID:   1,
		Username: "testuser",
		Email:    "test@example.com",
	})

	time.Sleep(100 * time.Millisecond)
}

func ExampleNewBus() {
	driver := NewChannelDriver()
	defer driver.Close()

	bus := NewBus[UserCreatedEvent](driver)

	bus.Subscribe("user.created", func(event UserCreatedEvent) {
		fmt.Printf("User created: %s (%s)\n", event.Username, event.Email)
	})

	bus.Publish("user.created", UserCreatedEvent{
		UserID:   1,
		Username: "testuser",
		Email:    "test@example.com",
	})

	time.Sleep(100 * time.Millisecond)
}

func ExampleNewBus_multipleTopics() {
	driver := NewChannelDriver()
	defer driver.Close()

	bus := NewBus[interface{}](driver)

	bus.Subscribe("user.created", func(event interface{}) {
		if userEvent, ok := event.(UserCreatedEvent); ok {
			fmt.Printf("User created: %s\n", userEvent.Username)
		}
	})

	bus.Subscribe("post.published", func(event interface{}) {
		if postEvent, ok := event.(PostPublishedEvent); ok {
			fmt.Printf("Post published: %s\n", postEvent.Title)
		}
	})

	bus.Publish("user.created", UserCreatedEvent{
		UserID:   1,
		Username: "testuser",
		Email:    "test@example.com",
	})

	bus.Publish("post.published", PostPublishedEvent{
		PostID:    1,
		Title:     "My First Post",
		AuthorID:  1,
		Published: true,
	})

	time.Sleep(100 * time.Millisecond)
}

func ExampleNewBus_concurrent() {
	driver := NewChannelDriver()
	defer driver.Close()

	bus := NewBus[UserCreatedEvent](driver)

	count := 0
	for i := 0; i < 10; i++ {
		bus.Subscribe("user.created", func(event UserCreatedEvent) {
			fmt.Printf("Handler %d received user: %s\n", count, event.Username)
			count++
		})
	}

	for i := 0; i < 5; i++ {
		go func(id int) {
			bus.Publish("user.created", UserCreatedEvent{
				UserID:   uint64(id),
				Username: fmt.Sprintf("user%d", id),
				Email:    fmt.Sprintf("user%d@example.com", id),
			})
		}(i)
	}

	time.Sleep(200 * time.Millisecond)
	fmt.Printf("Total events received: %d\n", count)
}
