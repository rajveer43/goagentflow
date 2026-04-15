package goagentflow

import "sync"

type State struct {
	mu       sync.RWMutex
	Input    any
	Output   any
	Step     int
	Values   map[string]any
	Messages []Message
}

func NewState(input any) *State {
	return &State{
		Input:  input,
		Values: make(map[string]any),
	}
}

func (s *State) Set(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Values[key] = value
}

func (s *State) Get(key string) (any, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.Values[key]
	return value, ok
}

func (s *State) AddMessage(msg Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Messages = append(s.Messages, msg)
}
