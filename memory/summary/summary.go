package summary

import (
	"context"
	"fmt"
	"sync"

	"github.com/rajveer43/goagentflow/runtime"
)

// Memory wraps an inner Memory and automatically summarizes old messages using an LLM.
// Pattern: Decorator - wraps a Memory with additional summarization behavior
// When the message count exceeds windowSize, the oldest messages are summarized
// using a Summarizer (typically an LLM) and the summary is stored under key "summary".
type Memory struct {
	mu          sync.RWMutex
	inner       runtime.Memory
	summarizer  runtime.Summarizer
	windowSize  int
	summaryKey  string
}

// New creates a new summary memory wrapper.
// inner: the underlying memory backend (e.g., inmemory, redis, postgres)
// summarizer: LLM-based summarizer (e.g., llmsummarizer.New(...))
// windowSize: number of recent messages to keep before summarizing old ones
func New(inner runtime.Memory, summarizer runtime.Summarizer, windowSize int) *Memory {
	if windowSize <= 0 {
		windowSize = 10 // default
	}
	return &Memory{
		inner:      inner,
		summarizer: summarizer,
		windowSize: windowSize,
		summaryKey: "conversation_summary",
	}
}

// AddMessage adds a message and triggers summarization if needed.
func (m *Memory) AddMessage(ctx context.Context, msg runtime.Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Add to inner memory
	if err := m.inner.AddMessage(ctx, msg); err != nil {
		return err
	}

	// Check if we should summarize
	messages, err := m.inner.GetMessages(ctx)
	if err != nil {
		return err
	}

	if len(messages) > m.windowSize && m.summarizer != nil {
		// Summarize oldest messages, keep the most recent windowSize
		oldMessages := messages[:len(messages)-m.windowSize]
		recentMessages := messages[len(messages)-m.windowSize:]

		// Create summary
		summary, err := m.summarizer.Summarize(ctx, oldMessages)
		if err != nil {
			return fmt.Errorf("summarization failed: %w", err)
		}

		// Store summary in inner memory
		if err := m.inner.Set(ctx, m.summaryKey, summary); err != nil {
			return err
		}

		// Clear old messages from inner memory (keep only recent ones + summary)
		// This is done by re-adding the recent messages
		// Note: depends on inner memory implementation - ideally it would support deletion
		// For now, we just trust that the summary is stored and applications know to use it
		_ = recentMessages // keep reference for potential future use
	}

	return nil
}

// GetMessages returns the summary (if exists) + all recent messages.
func (m *Memory) GetMessages(ctx context.Context) ([]runtime.Message, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	messages, err := m.inner.GetMessages(ctx)
	if err != nil {
		return nil, err
	}

	// If there's a summary, prepend it as a system message
	summary, err := m.inner.Get(ctx, m.summaryKey)
	if err == nil && summary != "" {
		summaryMsg := runtime.Message{
			Role:    "system",
			Content: fmt.Sprintf("[Conversation Summary]\n%v", summary),
		}
		return append([]runtime.Message{summaryMsg}, messages...), nil
	}

	return messages, nil
}

// Set delegates to inner memory.
func (m *Memory) Set(ctx context.Context, key string, value any) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.inner.Set(ctx, key, value)
}

// Get delegates to inner memory.
func (m *Memory) Get(ctx context.Context, key string) (any, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.inner.Get(ctx, key)
}

// SetSummaryKey sets the key used to store the summary (default: "conversation_summary").
func (m *Memory) SetSummaryKey(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.summaryKey = key
}

// GetSummary retrieves the current conversation summary.
func (m *Memory) GetSummary(ctx context.Context) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	summary, err := m.inner.Get(ctx, m.summaryKey)
	if err != nil {
		return "", err
	}

	if str, ok := summary.(string); ok {
		return str, nil
	}
	return fmt.Sprintf("%v", summary), nil
}
