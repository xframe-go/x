package xevent

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"github.com/xframe-go/x/contracts"
)

type (
	Subscriber struct {
		listeners []contracts.XListener
		queue     chan *AsyncSubscriber
	}

	AsyncSubscriber struct {
		events   contracts.XEvent
		listener contracts.XListener
	}
)

func NewSubscriber() *Subscriber {
	return &Subscriber{
		listeners: make([]contracts.XListener, 0),
	}
}

func (s *Subscriber) Add(listeners ...contracts.XListener) {
	for i := range listeners {
		listener := listeners[i]
		s.listeners = append(s.listeners, WrapListener(listener))
	}
}

func (s *Subscriber) All() []contracts.XListener {
	return s.listeners
}

type SyncListener struct {
	listener contracts.XListener
}

func NewSyncListener(listener contracts.XListener) *SyncListener {
	return &SyncListener{
		listener: listener,
	}
}

func (s *SyncListener) Handle(payload any) error {
	return s.listener.Handle(payload)
}

func (s *SyncListener) Stop(ctx context.Context) error {
	return nil
}

type AsyncListener struct {
	listener  contracts.XListener
	queue     chan any
	done      chan struct{}
	errors    chan error
	wg        sync.WaitGroup
	closeOnce sync.Once
	isClosed  atomic.Bool
	mu        sync.RWMutex
}

func WrapListener(listener contracts.XListener) contracts.XListener {
	async, ok := listener.(contracts.AsyncListener)
	if !ok {
		return NewSyncListener(listener)
	}

	bufferSize := async.AsyncBuffer()
	h := &AsyncListener{
		listener: listener,
		queue:    make(chan any, bufferSize),
		done:     make(chan struct{}),
		errors:   make(chan error, bufferSize),
		isClosed: atomic.Bool{},
	}

	for i := 0; i < bufferSize; i++ {
		h.wg.Add(1)
		go h.worker()
	}

	return h
}

func (a *AsyncListener) Handle(payload any) error {
	if a.isClosed.Load() {
		return ErrListenerClosed
	}
	a.queue <- payload
	return nil
}

func (a *AsyncListener) worker() {
	defer a.wg.Done()

	for {
		select {
		case <-a.done:
			return
		case evt := <-a.queue:
			if err := a.listener.Handle(evt); err != nil {
				select {
				case a.errors <- err:
				default:
				}
			}
		}
	}
}

func (a *AsyncListener) Stop(ctx context.Context) error {
	var err error
	a.closeOnce.Do(func() {
		a.isClosed.Store(true)

		close(a.done)

		done := make(chan struct{})
		go func() {
			a.wg.Wait()
			close(done)
		}()

		select {
		case <-ctx.Done():
			err = ctx.Err()
		case <-done:
			close(a.queue)
			close(a.errors)
		}
	})

	return err
}

func (a *AsyncListener) Errors() <-chan error {
	return a.errors
}

func (a *AsyncListener) QueueSize() int {
	return len(a.queue)
}

func (a *AsyncListener) IsClosed() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.isClosed.Load()
}

var (
	ErrListenerClosed = errors.New("listener is closed")
	ErrQueueFull      = errors.New("xevent queue is full")
)
