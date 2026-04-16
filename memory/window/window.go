package window

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/rajveer43/goagentflow/runtime"
)

// Memory is a token-aware sliding window memory.
// Pattern: Repository - maintains a memory window with approximate token counting
// Keeps messages until the total token count exceeds maxTokens, then evicts oldest.
// Uses word-count heuristic: ~1.3 tokens per word (approximate for English).
type Memory struct {
	mu        sync.RWMutex
	messages  []runtime.Message
	maxTokens int
	values    map[string]any
}

// New creates a new token-aware window memory.
// maxTokens: maximum tokens to keep in the window (approximate)
func New(maxTokens int) *Memory {
	if maxTokens <= 0 {
		maxTokens = 4096 // default: typical GPT context
	}
	return &Memory{
		messages:  make([]runtime.Message, 0),
		maxTokens: maxTokens,
		values:    make(map[string]any),
	}
}

// AddMessage adds a message and evicts old ones if over token budget.
func (m *Memory) AddMessage(ctx context.Context, msg runtime.Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.messages = append(m.messages, msg)

	// Check total tokens and evict oldest messages if necessary
	for totalTokens(m.messages) > m.maxTokens && len(m.messages) > 1 {
		m.messages = m.messages[1:] // evict oldest
	}

	return nil
}

// GetMessages returns all messages currently in the window.
func (m *Memory) GetMessages(ctx context.Context) ([]runtime.Message, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]runtime.Message, len(m.messages))
	copy(result, m.messages)
	return result, nil
}

// Set stores a key-value pair.
func (m *Memory) Set(ctx context.Context, key string, value any) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.values[key] = value
	return nil
}

// Get retrieves a value by key.
func (m *Memory) Get(ctx context.Context, key string) (any, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, ok := m.values[key]
	if !ok {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	return value, nil
}

// EstimateTokens returns the approximate token count of current messages.
func (m *Memory) EstimateTokens(ctx context.Context) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return totalTokens(m.messages), nil
}

// SetMaxTokens updates the max token budget.
func (m *Memory) SetMaxTokens(maxTokens int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.maxTokens = maxTokens
}

// totalTokens estimates the total token count of messages using word-count heuristic.
// Heuristic: approximately 1.3 tokens per word on average.
func totalTokens(messages []runtime.Message) int {
	total := 0
	for _, msg := range messages {
		// Count words (split by whitespace)
		words := strings.Fields(msg.Content)
		// Estimate tokens: 1 token per word * 1.3 (conservative estimate)
		tokens := int(float64(len(words)) * 1.3)
		// Ensure at least 1 token per message (for overhead)
		if tokens < 1 {
			tokens = 1
		}
		total += tokens
	}
	return total
}
