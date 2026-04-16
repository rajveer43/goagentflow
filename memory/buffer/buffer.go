package buffer

import (
	"context"
	"fmt"
	"sync"

	"github.com/rajveer43/goagentflow/runtime"
)

// Memory is an in-memory FIFO buffer that keeps only the last N messages.
// Pattern: Repository - sliding window memory management
// Older messages are evicted when the buffer reaches maxMessages capacity.
// Thread-safe with RWMutex.
type Memory struct {
	mu          sync.RWMutex
	messages    []runtime.Message
	maxMessages int
	values      map[string]any // key-value store
}

// New creates a new buffer memory with maximum message capacity.
// maxMessages: maximum number of messages to keep (older messages are evicted)
func New(maxMessages int) *Memory {
	if maxMessages <= 0 {
		maxMessages = 100 // default
	}
	return &Memory{
		messages:    make([]runtime.Message, 0, maxMessages),
		maxMessages: maxMessages,
		values:      make(map[string]any),
	}
}

// AddMessage appends a message to the buffer.
// When the buffer exceeds maxMessages, the oldest message is evicted.
func (m *Memory) AddMessage(ctx context.Context, msg runtime.Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.messages = append(m.messages, msg)

	// Evict oldest messages if over capacity
	if len(m.messages) > m.maxMessages {
		m.messages = m.messages[len(m.messages)-m.maxMessages:]
	}

	return nil
}

// GetMessages returns all messages currently in the buffer (in order).
func (m *Memory) GetMessages(ctx context.Context) ([]runtime.Message, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy to prevent external modification
	result := make([]runtime.Message, len(m.messages))
	copy(result, m.messages)
	return result, nil
}

// Set stores a key-value pair in the buffer.
func (m *Memory) Set(ctx context.Context, key string, value any) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.values[key] = value
	return nil
}

// Get retrieves a value by key, returning (nil, false) if not found.
func (m *Memory) Get(ctx context.Context, key string) (any, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, ok := m.values[key]
	if !ok {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	return value, nil
}

// Clear removes all messages and values from the buffer.
func (m *Memory) Clear(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.messages = m.messages[:0]
	m.values = make(map[string]any)
	return nil
}

// Size returns the current number of messages in the buffer.
func (m *Memory) Size(ctx context.Context) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.messages), nil
}
