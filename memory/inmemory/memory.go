package inmemory

import (
	"context"
	"sync"

	"goagentflow/runtime"
)

type Memory struct {
	mu       sync.RWMutex
	values   map[string]any
	messages []runtime.Message
}

func New() *Memory {
	return &Memory{values: make(map[string]any)}
}

func (m *Memory) AddMessage(_ context.Context, msg runtime.Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = append(m.messages, msg)
	return nil
}

func (m *Memory) GetMessages(_ context.Context) ([]runtime.Message, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]runtime.Message, len(m.messages))
	copy(out, m.messages)
	return out, nil
}

func (m *Memory) Set(_ context.Context, key string, value any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.values[key] = value
	return nil
}

func (m *Memory) Get(_ context.Context, key string) (any, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.values[key], nil
}
