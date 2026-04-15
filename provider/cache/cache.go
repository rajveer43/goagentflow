package cache

import (
	"context"
	"crypto/md5"
	"fmt"
	"sync"

	"github.com/rajveer43/goagentflow/runtime"
)

// LRU Cache (Least Recently Used) using doubly-linked list + HashMap
// DSA: Doubly-linked list for O(1) eviction, hashmap for O(1) lookup
// Caches only Complete() calls, not Stream() calls

// cacheEntry represents a cached LLM response
type cacheEntry struct {
	key    string
	value  string
	node   *listNode
	ttl    int64 // not used yet, but ready for expiry
}

// listNode is part of doubly-linked list for LRU tracking
type listNode struct {
	prev  *listNode
	next  *listNode
	entry *cacheEntry
}

// CachedLLM wraps any runtime.LLM and caches Complete() responses.
// Pattern: Decorator - transparent caching layer over any LLM provider
type CachedLLM struct {
	underlying runtime.LLM
	capacity   int

	mu       sync.RWMutex
	cache    map[string]*cacheEntry // key (hash of prompt+temp) -> entry
	head     *listNode               // most recently used
	tail     *listNode               // least recently used
	size     int                      // current size
}

// New creates a new cached LLM wrapper.
// underlying: any runtime.LLM implementation
// capacity: max number of entries to cache
func New(underlying runtime.LLM, capacity int) *CachedLLM {
	return &CachedLLM{
		underlying: underlying,
		capacity:   capacity,
		cache:      make(map[string]*cacheEntry),
		head:       nil,
		tail:       nil,
		size:       0,
	}
}

// cacheKey generates a deterministic key from prompt and temperature.
// Uses MD5 hash for compact key representation.
func cacheKey(prompt string, temp float64) string {
	hash := md5.Sum([]byte(fmt.Sprintf("%s:%.2f", prompt, temp)))
	return fmt.Sprintf("%x", hash)
}

// Complete returns a cached response if available, otherwise calls the underlying LLM.
func (c *CachedLLM) Complete(ctx context.Context, prompt string, opts ...runtime.LLMOption) (string, error) {
	// Extract temperature from opts to include in cache key
	cfg := &runtime.LLMConfig{Temperature: 0.7}
	for _, opt := range opts {
		opt(cfg)
	}

	key := cacheKey(prompt, cfg.Temperature)

	// Try cache hit (read lock)
	c.mu.RLock()
	if entry, found := c.cache[key]; found {
		// Move to front (most recently used)
		c.mu.RUnlock()
		c.moveToFront(entry.node)
		return entry.value, nil
	}
	c.mu.RUnlock()

	// Cache miss: call underlying LLM
	response, err := c.underlying.Complete(ctx, prompt, opts...)
	if err != nil {
		return "", err
	}

	// Store in cache
	c.put(key, response)
	return response, nil
}

// Stream is not cached — streams are passed through directly.
// Caching streaming responses would require buffering entire responses anyway.
func (c *CachedLLM) Stream(ctx context.Context, prompt string, opts ...runtime.LLMOption) (<-chan string, <-chan error) {
	return c.underlying.Stream(ctx, prompt, opts...)
}

// put adds or updates a cache entry.
// DSA: Evicts LRU entry (tail) if capacity is exceeded.
func (c *CachedLLM) put(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if key already exists
	if entry, found := c.cache[key]; found {
		entry.value = value
		c.moveToFrontLocked(entry.node)
		return
	}

	// Create new entry
	node := &listNode{
		entry: &cacheEntry{
			key:   key,
			value: value,
		},
	}
	node.entry.node = node

	// Add to front
	if c.head == nil {
		c.head = node
		c.tail = node
	} else {
		node.next = c.head
		c.head.prev = node
		c.head = node
	}

	c.cache[key] = node.entry
	c.size++

	// Evict LRU if over capacity
	if c.size > c.capacity {
		c.evictLocked()
	}
}

// moveToFront moves a node to the front (most recently used).
func (c *CachedLLM) moveToFront(node *listNode) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.moveToFrontLocked(node)
}

// moveToFrontLocked moves a node to front (assumes lock is held).
func (c *CachedLLM) moveToFrontLocked(node *listNode) {
	if node == c.head {
		return // already at front
	}

	// Remove from current position
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		c.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		c.tail = node.prev
	}

	// Add to front
	node.prev = nil
	node.next = c.head
	if c.head != nil {
		c.head.prev = node
	}
	c.head = node

	if c.tail == nil {
		c.tail = node
	}
}

// evictLocked removes the LRU entry (tail).
// Assumes lock is held.
func (c *CachedLLM) evictLocked() {
	if c.tail == nil {
		return
	}

	key := c.tail.entry.key
	delete(c.cache, key)
	c.size--

	if c.tail.prev != nil {
		c.tail.prev.next = nil
		c.tail = c.tail.prev
	} else {
		c.head = nil
		c.tail = nil
	}
}

// Stats returns cache statistics (for monitoring).
type Stats struct {
	Size     int
	Capacity int
	Entries  map[string]bool
}

// GetStats returns current cache statistics.
func (c *CachedLLM) GetStats() Stats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entries := make(map[string]bool)
	for key := range c.cache {
		entries[key] = true
	}

	return Stats{
		Size:     c.size,
		Capacity: c.capacity,
		Entries:  entries,
	}
}
