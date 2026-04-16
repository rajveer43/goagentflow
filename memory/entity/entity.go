package entity

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/rajveer43/goagentflow/runtime"
)

// EntityInfo holds information about a tracked entity.
type EntityInfo struct {
	Name       string // the entity name
	Type       string // e.g., "person", "place", "date", "concept"
	Mentions   int    // how many times mentioned
	LastMention string // last context it was mentioned in
	Summary    string // summary of what we know about it
}

// Memory tracks named entities mentioned in conversations.
// Pattern: Repository - maintains entity knowledge base alongside message history
// Uses simple regex patterns to extract entities; no external NLP dependency.
type Memory struct {
	mu        sync.RWMutex
	inner     runtime.Memory
	entities  map[string]EntityInfo
	messages  []runtime.Message
	values    map[string]any
}

// New creates a new entity-tracking memory wrapper.
func New(inner runtime.Memory) *Memory {
	return &Memory{
		inner:    inner,
		entities: make(map[string]EntityInfo),
		messages: make([]runtime.Message, 0),
		values:   make(map[string]any),
	}
}

// AddMessage adds a message and extracts entities from it.
func (m *Memory) AddMessage(ctx context.Context, msg runtime.Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Add to inner memory
	if err := m.inner.AddMessage(ctx, msg); err != nil {
		return err
	}

	// Extract and track entities
	m.extractEntities(msg.Content)

	return nil
}

// GetMessages returns messages with entity context prepended.
func (m *Memory) GetMessages(ctx context.Context) ([]runtime.Message, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	messages, err := m.inner.GetMessages(ctx)
	if err != nil {
		return nil, err
	}

	// Prepend entity summary as system message if any entities exist
	if len(m.entities) > 0 {
		entitySummary := m.buildEntitySummary()
		entityMsg := runtime.Message{
			Role:    "system",
			Content: fmt.Sprintf("[Known Entities]\n%s", entitySummary),
		}
		return append([]runtime.Message{entityMsg}, messages...), nil
	}

	return messages, nil
}

// Set delegates to inner memory.
func (m *Memory) Set(ctx context.Context, key string, value any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.inner.Set(ctx, key, value)
}

// Get delegates to inner memory.
func (m *Memory) Get(ctx context.Context, key string) (any, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.inner.Get(ctx, key)
}

// GetEntities returns all tracked entities.
func (m *Memory) GetEntities(ctx context.Context) (map[string]EntityInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]EntityInfo)
	for k, v := range m.entities {
		result[k] = v
	}
	return result, nil
}

// AddEntity manually adds or updates an entity.
func (m *Memory) AddEntity(ctx context.Context, name, entityType, summary string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if name == "" {
		return fmt.Errorf("entity name cannot be empty")
	}

	info := m.entities[name]
	info.Name = name
	info.Type = entityType
	info.Summary = summary
	info.Mentions++
	m.entities[name] = info

	return nil
}

// extractEntities uses simple regex patterns to find potential entities.
func (m *Memory) extractEntities(text string) {
	// Extract capitalized words (likely proper nouns)
	// Pattern: word starting with capital followed by 1+ capital letters or lowercase letters
	properNounPattern := regexp.MustCompile(`\b[A-Z][a-z]+(?:\s+[A-Z][a-z]+)*\b`)
	nouns := properNounPattern.FindAllString(text, -1)

	for _, noun := range nouns {
		noun = strings.TrimSpace(noun)
		if noun == "" || len(noun) < 2 {
			continue
		}

		info, exists := m.entities[noun]
		if !exists {
			info = EntityInfo{
				Name:    noun,
				Type:    "unknown", // classifier could be improved
				Summary: "",
			}
		}

		info.Mentions++
		info.LastMention = text[:min(len(text), 100)] // store first 100 chars as context
		m.entities[noun] = info
	}

	// Extract dates (simple pattern)
	datePattern := regexp.MustCompile(`\b\d{4}-\d{2}-\d{2}\b|\b(?:Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)[a-z]*\s+\d{1,2}(?:,?\s+\d{4})?\b`)
	dates := datePattern.FindAllString(text, -1)

	for _, date := range dates {
		date = strings.TrimSpace(date)
		if date == "" {
			continue
		}

		info, exists := m.entities[date]
		if !exists {
			info = EntityInfo{
				Name:    date,
				Type:    "date",
				Summary: "",
			}
		}

		info.Mentions++
		info.LastMention = text[:min(len(text), 100)]
		m.entities[date] = info
	}
}

// buildEntitySummary creates a human-readable summary of tracked entities.
func (m *Memory) buildEntitySummary() string {
	if len(m.entities) == 0 {
		return "(No entities tracked)"
	}

	var sb strings.Builder
	for _, entity := range m.entities {
		sb.WriteString(fmt.Sprintf("- %s (%s): mentioned %d times", entity.Name, entity.Type, entity.Mentions))
		if entity.Summary != "" {
			sb.WriteString(fmt.Sprintf(" - %s", entity.Summary))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
