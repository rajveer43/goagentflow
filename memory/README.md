# Memory Backends

All memory implementations follow the `runtime.Memory` interface, making them interchangeable and composable.

---

## 📋 Available Memory Types

| Type | Use Case | Memory | Speed | Token Aware |
|------|----------|--------|-------|-------------|
| **InMemory** | Testing, short sessions | O(n) | ⚡ Fast | ❌ No |
| **Buffer** | Short conversations | O(k) | ⚡ Fast | ❌ No |
| **Window** | Long conversations | O(tokens) | ⚡ Fast | ✅ Yes |
| **Entity** | Multi-turn with recall | O(entities) | 🟡 Medium | ❌ No |
| **Summary** | Very long conversations | O(summaries) | 🔴 Slow | ✅ Yes |
| **Compressive** | Compression on demand | O(compressed) | 🔴 Slow | ✅ Yes |

---

## 🚀 Quick Start

### Simple In-Memory Storage

```go
import "github.com/rajveer43/goagentflow/memory/inmemory"

// Create memory
mem := inmemory.New()

// Add messages
mem.AddMessage(ctx, runtime.Message{Role: "user", Content: "Hi"})
mem.AddMessage(ctx, runtime.Message{Role: "assistant", Content: "Hello"})

// Retrieve messages
messages, _ := mem.GetMessages(ctx)
fmt.Println(messages)  // All messages
```

### Limited Window (Cost Control)

```go
import "github.com/rajveer43/goagentflow/memory/window"

// Keep only ~4096 tokens of history
mem := window.New(4096)

// Add messages - old ones are evicted when token budget exceeded
mem.AddMessage(ctx, runtime.Message{Role: "user", Content: "Long story..."})
mem.AddMessage(ctx, runtime.Message{Role: "user", Content: "Another message"})

// GetMessages returns only recent messages within token budget
messages, _ := mem.GetMessages(ctx)
```

### Entity Tracking (Remember Details)

```go
import "github.com/rajveer43/goagentflow/memory/entity"

base := inmemory.New()
mem := entity.New(base)

// Automatically extracts and remembers entities (people, places, concepts)
mem.AddMessage(ctx, runtime.Message{Role: "user", Content: "I'm Alice from NYC"})
mem.AddMessage(ctx, runtime.Message{Role: "user", Content: "What do you know about me?"})

// Entity memory injects extracted entities into context
messages, _ := mem.GetMessages(ctx)
// Includes: User: Alice, Location: NYC
```

### Auto-Summarization (Very Long Conversations)

```go
import "github.com/rajveer43/goagentflow/memory/summary"

base := inmemory.New()
llmSummarizer := anthropic.New(apiKey, "claude-opus-4-6")
mem := summary.New(base, llmSummarizer, 10)  // Summarize after 10 messages

// Add messages
for i := 0; i < 20; i++ {
    mem.AddMessage(ctx, runtime.Message{...})
}

// Old messages are summarized, recent messages kept verbatim
messages, _ := mem.GetMessages(ctx)
// Messages 1-10: Summarized
// Messages 11-20: Verbatim
```

---

## 🏗️ Composing Memory (Decorator Pattern)

Stack memory backends for complex scenarios:

### Pattern 1: Window + Entity Tracking

```go
base := inmemory.New()
withEntity := entity.New(base)           // Track entities
withWindow := window.New(withEntity, 4096)  // Respect token budget

agent := agent.New(llm, withWindow)
```

**Behavior:**
- InMemory stores all messages
- Entity layer extracts named entities
- Window layer evicts old messages to stay within 4096 tokens

### Pattern 2: Summary + Window

```go
base := inmemory.New()
withSummary := summary.New(base, llm, 10)    // Summarize after 10 messages
withWindow := window.New(withSummary, 8192)  // Keep within 8192 tokens

agent := agent.New(llm, withWindow)
```

**Behavior:**
- Stores all messages
- Every 10 messages, summarizes older ones
- Respects token budget

### Pattern 3: Full Stack

```go
base := inmemory.New()
withEntity := entity.New(base)
withWindow := window.New(withEntity, 4096)
withSummary := summary.New(withWindow, llm, 10)
withCompression := compressive.New(withSummary, llm, 0.9)  // Compress at 90% threshold

agent := agent.New(llm, withCompression)
```

**Behavior:**
- Tracks entities
- Respects token budget
- Auto-summarizes
- Compresses when 90% full

---

## 📊 Memory Type Comparison

### InMemory
Best for: Testing, short sessions

```go
mem := inmemory.New()
```

**Characteristics:**
- ✅ Simplest
- ✅ Fastest
- ❌ No size limit
- ❌ Grows unbounded

**When to use:**
- Unit tests
- Development
- Single-session apps
- Short conversations (<100 messages)

### Buffer

Best for: Fixed-size sliding window

```go
mem := buffer.New(10)  // Keep last 10 messages
```

**Characteristics:**
- ✅ Bounded memory (O(k))
- ✅ Fast
- ❌ Naive FIFO (doesn't understand token cost)

**When to use:**
- Short conversations with bounded memory
- Simple apps with known message count
- Example: Chat UI with "last 10 messages"

### Window

Best for: Token-aware conversations

```go
mem := window.New(4096)  // Keep ~4096 tokens
```

**Characteristics:**
- ✅ Respects token budget
- ✅ Fast
- ✅ Configurable threshold
- ⚠️ Approximate token counting

**When to use:**
- Long conversations
- Cost-sensitive applications
- Model context budget management

**Example: Token-aware chat**
```go
// Config: 4096 token budget for context
// Model: Claude with 200K context, but we only use 4K for history
mem := window.New(4096)

// Add 1000 token message
mem.AddMessage(ctx, Message{Content: longText})
// Add 2000 token message
mem.AddMessage(ctx, Message{Content: anotherLongText})
// Add 1500 token message
mem.AddMessage(ctx, Message{Content: thirdMessage})

// Now at 4500 tokens, window evicts oldest message (1000 tokens)
// Keeps newest 3 messages (2000 + 1500 = 3500 tokens)
mem.AddMessage(ctx, Message{Content: newMessage})
```

### Entity

Best for: Multi-turn conversations needing entity recall

```go
mem := entity.New(inmemory.New())  // Wrap another memory
```

**Characteristics:**
- ✅ Tracks entities (people, places, concepts)
- ✅ Injects entity summaries into context
- 🟡 Medium speed
- ⚠️ Entity extraction quality depends on LLM

**When to use:**
- Customer support (remember customer details)
- Long conversations with multiple entities
- Need to reference past mentions

**Example:**
```go
mem := entity.New(inmemory.New())

mem.AddMessage(ctx, Message{
    Role: "user",
    Content: "I'm John, I live in London, I work at Acme Corp",
})

mem.AddMessage(ctx, Message{
    Role: "assistant",
    Content: "Nice to meet you John!",
})

mem.AddMessage(ctx, Message{
    Role: "user",
    Content: "What do you know about me?",
})

// Memory injects extracted entities:
// - PERSON: John
// - LOCATION: London  
// - ORG: Acme Corp
// These are included when building LLM context
```

### Summary

Best for: Very long conversations

```go
mem := summary.New(inmemory.New(), llm, 10)  // Summarize every 10 messages
```

**Characteristics:**
- ✅ Handles very long conversations
- ✅ Semantic compression with LLM
- 🔴 Slow (calls LLM frequently)
- 💰 Expensive (LLM summarization costs)

**When to use:**
- Very long conversations (100+ messages)
- Need semantic compression
- Cost not critical

**Example:**
```go
mem := summary.New(inmemory.New(), llm, 10)

// Add messages 1-9: Stored verbatim
// Add message 10: Triggers summarization
// Messages 1-10 are summarized to "User discussed X, Y, Z"
// Messages 11+: Stored verbatim again

// GetMessages returns:
// - Summary of messages 1-10
// - Verbatim messages 11+
```

### Compressive

Best for: Compression on demand

```go
mem := compressive.New(inmemory.New(), llm, 0.9)  // Compress at 90% full
```

**Characteristics:**
- ✅ Compresses when threshold exceeded
- 🔴 LLM-based (slow, expensive)
- ⚠️ Semantic loss from compression

**When to use:**
- Need maximum control over compression
- Hybrid approach combining summarization and compression

---

## 🛠️ Adding Custom Memory

Want to add support for a new memory backend (Redis, Postgres, DynamoDB)? See [docs/EXTENDING.md](../docs/EXTENDING.md).

### Quick Template

```go
// memory/mybackend/memory.go
package mybackend

import (
    "context"
    "github.com/rajveer43/goagentflow/runtime"
)

type Memory struct {
    backend Backend  // Your storage backend
}

func New(backend Backend) *Memory {
    return &Memory{backend: backend}
}

func (m *Memory) AddMessage(ctx context.Context, msg runtime.Message) error {
    // Store message in backend
    return nil
}

func (m *Memory) GetMessages(ctx context.Context) ([]runtime.Message, error) {
    // Retrieve messages from backend
    return nil, nil
}

func (m *Memory) Set(ctx context.Context, key string, value any) error {
    // Store arbitrary key-value
    return nil
}

func (m *Memory) Get(ctx context.Context, key string) (any, error) {
    // Retrieve arbitrary key-value
    return nil, nil
}

var _ runtime.Memory = (*Memory)(nil)
```

### Backends You Could Implement

- **Redis** - Fast, distributed memory
- **Postgres** - Persistent, queryable
- **DynamoDB** - Serverless, scalable
- **MongoDB** - Document-based, flexible
- **SQLite** - File-based, embedded
- **Memcached** - Distributed caching

---

## 📚 Examples

### Example 1: Simple Chatbot

```go
import (
    "github.com/rajveer43/goagentflow/memory/inmemory"
    "github.com/rajveer43/goagentflow/chains/agent"
)

// Simple in-memory storage
memory := inmemory.New()
agent := agent.New(llm, memory)

// Turn 1
response1, _ := agent.Run(ctx, "My favorite color is blue")
// Turn 2
response2, _ := agent.Run(ctx, "What's my favorite color?")
// Agent remembers "blue"
```

### Example 2: Cost-Conscious App

```go
import "github.com/rajveer43/goagentflow/memory/window"

// Keep only 4096 tokens (roughly $0.004 in context)
memory := window.New(4096)
agent := agent.New(llm, memory)

// App can handle long conversations without exploding context size
for {
    input := getUserInput()
    response, _ := agent.Run(ctx, input)
    fmt.Println(response)
}
```

### Example 3: Customer Support Bot

```go
import (
    "github.com/rajveer43/goagentflow/memory/entity"
    "github.com/rajveer43/goagentflow/memory/inmemory"
)

// Remember customer details
memory := entity.New(inmemory.New())
agent := agent.New(llm, memory)

// User: "I'm John Doe, account 12345"
agent.Run(ctx, "I'm John Doe, account 12345")

// User: "What's my account number?"
// Agent remembers: account 12345 (extracted as entity)
agent.Run(ctx, "What's my account number?")
```

### Example 4: Efficient Long Conversation

```go
import (
    "github.com/rajveer43/goagentflow/memory/window"
    "github.com/rajveer43/goagentflow/memory/summary"
)

// Combine window (bounded) with summary (semantic compression)
base := inmemory.New()
withSummary := summary.New(base, llm, 20)     // Summarize every 20 messages
withWindow := window.New(withSummary, 8192)   // Keep within 8K tokens

agent := agent.New(llm, withWindow)

// Can have 100+ message conversations efficiently
// Old messages are summarized, recent ones kept verbatim
```

---

## 🔐 Security & Privacy

### Data Persistence

```go
// ✅ Good: Clear control
// InMemory - data lost on restart (good for testing)
mem := inmemory.New()

// ✅ Good: Explicit persistence
// Redis - persistent if configured
mem := custom.NewRedisMemory(redisClient)

// ⚠️ Be careful: Know where data goes
// May need encryption, backups, compliance handling
```

### API Keys in Memory

```go
// ❌ BAD: Don't store API keys in conversation memory
mem.AddMessage(ctx, Message{
    Content: "My API key is sk-...",  // DANGER!
})

// ✅ GOOD: Extract sensitive data before storing
sanitizedContent := removeSensitiveData(userMessage)
mem.AddMessage(ctx, Message{
    Content: sanitizedContent,
})
```

---

## 📈 Performance Tips

### Avoid Unbounded Growth

```go
// ❌ Bad: Will grow forever
mem := inmemory.New()

// ✅ Good: Bounded with token window
mem := window.New(4096)

// ✅ Good: Bounded with message count
mem := buffer.New(100)

// ✅ Good: Bounded with summarization
mem := summary.New(inmemory.New(), llm, 10)
```

### Optimize Token Counting

```go
// Window memory uses approximate token counting
// For exact counts, consider:

// Option 1: Use LLM with token counter interface
counter, ok := llm.(runtime.TokenCounter)
if ok {
    tokens := counter.CountTokens(text)
}

// Option 2: Use encoding library (tiktoken for Python-style)
// Note: Go implementation may vary
```

---

## 🐛 Troubleshooting

### Messages Not Persisting

```go
// Check: Are you using InMemory (lost on restart)?
mem := inmemory.New()  // ❌ Lost when program exits

// Use persistent backend instead
mem := custom.NewRedisMemory(client)  // ✅ Persists
```

### Memory Growing Too Large

```go
// Check: Are you using unbounded memory?
mem := inmemory.New()  // ❌ Grows forever

// Use window or buffer
mem := window.New(4096)  // ✅ Bounded to 4096 tokens
```

### Token Window Not Working

```go
// Window uses approximate token counting
// May not be exact, but is good enough for budgeting
// For production, consider:

// 1. Use Summary + Window together
mem := summary.New(
    window.New(inmemory.New(), 4096),
    llm,
    10,
)

// 2. Monitor actual token usage
memUsage, _ := mem.EstimateTokens()
if memUsage > 4096 * 0.9 {
    // Approaching limit
}
```

---

## 📖 More Info

- [Architecture Overview](../docs/ARCHITECTURE.md)
- [Extending Guide](../docs/EXTENDING.md)
- [Main README](../README.md)
- [Contributing Guide](../CONTRIBUTING.md)

---

Questions? [Start a discussion](https://github.com/rajveer43/goagentflow/discussions)!
