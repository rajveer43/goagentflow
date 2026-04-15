package stream

import "sync"

type Writer interface {
	Write(event any) error
}

type EventStream struct {
	C      chan any
	closed bool
	mu     sync.Mutex
}

func NewEventStream(buffer int) *EventStream {
	if buffer <= 0 {
		buffer = 1
	}
	return &EventStream{C: make(chan any, buffer)}
}

func (s *EventStream) TryEmit(event any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return
	}
	select {
	case s.C <- event:
	default:
	}
}

func (s *EventStream) Write(event any) error {
	s.TryEmit(event)
	return nil
}

func (s *EventStream) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return
	}
	s.closed = true
	close(s.C)
}
