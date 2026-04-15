package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/rajveer43/goagentflow/runtime"
)

// Memory implements runtime.Memory using Redis as persistent backend.
// Uses Redis LIST for messages (ordered) and HASH for key-value storage.
// Pattern: Repository + Factory
// DSA: Redis intrinsic B-tree indexes, O(1) HGET/HSET, O(n) LRANGE for messages
type Memory struct {
	client    *redis.Client
	sessionID string
	mu        sync.RWMutex
}

// New creates a new Redis-backed memory instance.
// addr: Redis address (e.g., "localhost:6379")
// sessionID: unique identifier for this agent's memory scope
func New(addr, sessionID string) (*Memory, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5000000000) // 5s
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return &Memory{
		client:    client,
		sessionID: sessionID,
	}, nil
}

// AddMessage adds a message to the message history (append-only list).
// Uses Redis LPUSH to prepend, maintaining insertion order via negative indices on retrieval.
func (m *Memory) AddMessage(ctx context.Context, msg runtime.Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("messages:%s", m.sessionID)
	return m.client.LPush(ctx, key, string(data)).Err()
}

// GetMessages retrieves all messages in order (oldest to newest).
// Uses Redis LRANGE with negative indices to reverse order.
func (m *Memory) GetMessages(ctx context.Context) ([]runtime.Message, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key := fmt.Sprintf("messages:%s", m.sessionID)
	results, err := m.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	messages := make([]runtime.Message, 0, len(results))
	// Reverse iteration since LPUSH added them in reverse order
	for i := len(results) - 1; i >= 0; i-- {
		var msg runtime.Message
		if err := json.Unmarshal([]byte(results[i]), &msg); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// Set stores a key-value pair using Redis HASH.
// Values are stored as JSON for arbitrary types (matches inmemory impl).
func (m *Memory) Set(ctx context.Context, key string, value any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	hashKey := fmt.Sprintf("kv:%s", m.sessionID)
	return m.client.HSet(ctx, hashKey, key, string(data)).Err()
}

// Get retrieves a value from the key-value store.
func (m *Memory) Get(ctx context.Context, key string) (any, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	hashKey := fmt.Sprintf("kv:%s", m.sessionID)
	val, err := m.client.HGet(ctx, hashKey, key).Result()
	if err == redis.Nil {
		return nil, nil // key not found (matches inmemory behavior)
	}
	if err != nil {
		return nil, err
	}

	var result any
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Close closes the Redis connection.
func (m *Memory) Close() error {
	return m.client.Close()
}
