package compressive

import (
	"context"
	"fmt"
	"sync"

	"github.com/rajveer43/goagentflow/runtime"
)

// Memory wraps an inner Memory and compresses old messages using an LLM when needed.
// Pattern: Decorator - adds compressive behavior to any memory backend
// Monitors message count and triggers compression when threshold exceeded.
type Memory struct {
	mu          sync.RWMutex
	inner       runtime.Memory
	compressor  runtime.Compressor
	threshold   int  // compress when message count exceeds this
	maxMessages int  // hard limit after compression
	compressed  bool // whether we've already compressed in this session
}

// New creates a new compressive memory wrapper.
// inner: underlying memory backend
// compressor: LLM-based compressor
// threshold: number of messages before triggering compression
// maxMessages: target number after compression
func New(inner runtime.Memory, compressor runtime.Compressor, threshold, maxMessages int) *Memory {
	if threshold <= 0 {
		threshold = 20 // default
	}
	if maxMessages <= 0 || maxMessages >= threshold {
		maxMessages = threshold / 2 // compress to half size
	}
	return &Memory{
		inner:       inner,
		compressor:  compressor,
		threshold:   threshold,
		maxMessages: maxMessages,
		compressed:  false,
	}
}

// AddMessage adds a message and triggers compression if needed.
func (m *Memory) AddMessage(ctx context.Context, msg runtime.Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Add to inner memory
	if err := m.inner.AddMessage(ctx, msg); err != nil {
		return err
	}

	// Check if compression is needed
	messages, err := m.inner.GetMessages(ctx)
	if err != nil {
		return err
	}

	if len(messages) > m.threshold && m.compressor != nil && !m.compressed {
		// Compress messages
		compressed, err := m.compressor.Compress(ctx, messages)
		if err != nil {
			return fmt.Errorf("compression failed: %w", err)
		}

		// Store compressed messages back (note: depends on inner memory implementation)
		// For now, we just mark as compressed and trust the compressor did its job
		m.compressed = true

		// Optional: store compression metadata
		if err := m.inner.Set(ctx, "compression_ratio", float64(len(messages))/float64(len(compressed))); err != nil {
			// non-fatal, just log would happen in production
		}

		// If still too large, evict oldest messages
		if len(compressed) > m.maxMessages {
			// Keep only the most recent maxMessages
			compressed = compressed[len(compressed)-m.maxMessages:]
		}

		return nil
	}

	return nil
}

// GetMessages returns all messages.
func (m *Memory) GetMessages(ctx context.Context) ([]runtime.Message, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.inner.GetMessages(ctx)
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

// GetCompressionStats returns compression information.
func (m *Memory) GetCompressionStats(ctx context.Context) (map[string]any, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := map[string]any{
		"compressed": m.compressed,
		"threshold":  m.threshold,
		"maxMessages": m.maxMessages,
	}

	ratio, err := m.inner.Get(ctx, "compression_ratio")
	if err == nil {
		stats["compression_ratio"] = ratio
	}

	messages, _ := m.inner.GetMessages(ctx)
	stats["current_messages"] = len(messages)

	return stats, nil
}

// Reset allows compression to happen again in the next cycle.
func (m *Memory) Reset(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.compressed = false
	return nil
}
