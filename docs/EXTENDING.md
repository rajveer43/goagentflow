# Extending goagentflow

goagentflow is designed for easy extension. This guide shows how to implement custom components.

---

## 📋 Table of Contents

1. [Custom LLM Provider](#custom-llm-provider)
2. [Custom Memory Backend](#custom-memory-backend)
3. [Custom Chain](#custom-chain)
4. [Custom VectorStore](#custom-vectorstore)
5. [Testing Custom Components](#testing-custom-components)
6. [Publishing Your Extension](#publishing-your-extension)

---

## Custom LLM Provider

### Step 1: Create Package Structure

```bash
mkdir -p provider/myservice
touch provider/myservice/provider.go
touch provider/myservice/provider_test.go
touch provider/myservice/errors.go
```

### Step 2: Implement runtime.LLM Interface

**File: `provider/myservice/provider.go`**

```go
package myservice

import (
    "context"
    "fmt"
    "github.com/rajveer43/goagentflow/runtime"
)

// Provider wraps the MyService API client.
type Provider struct {
    apiKey   string
    endpoint string
    client   *MyServiceClient  // Your HTTP client
    timeout  time.Duration
}

// New creates a new MyService provider.
//
// apiKey: Your MyService API key
// endpoint: API endpoint (e.g., "https://api.myservice.com")
func New(apiKey, endpoint string) *Provider {
    return &Provider{
        apiKey:   apiKey,
        endpoint: endpoint,
        client:   NewMyServiceClient(apiKey, endpoint),
        timeout:  30 * time.Second,
    }
}

// SetTimeout sets the API timeout (default 30s).
func (p *Provider) SetTimeout(d time.Duration) {
    p.timeout = d
}

// Complete sends a prompt to MyService and returns the response.
func (p *Provider) Complete(ctx context.Context, prompt string) (string, error) {
    // 1. Validate inputs
    if prompt == "" {
        return "", fmt.Errorf("prompt cannot be empty")
    }
    
    // 2. Create context with timeout
    ctx, cancel := context.WithTimeout(ctx, p.timeout)
    defer cancel()
    
    // 3. Build request
    req := &MyServiceRequest{
        Model:  "myservice-1.0",
        Prompt: prompt,
    }
    
    // 4. Call API
    resp, err := p.client.Complete(ctx, req)
    if err != nil {
        return "", fmt.Errorf("api call failed: %w", err)
    }
    
    // 5. Validate response
    if resp.Text == "" {
        return "", fmt.Errorf("empty response from API")
    }
    
    // 6. Return response
    return resp.Text, nil
}

// Verify that Provider implements runtime.LLM
var _ runtime.LLM = (*Provider)(nil)
```

### Step 3: Add Error Types

**File: `provider/myservice/errors.go`**

```go
package myservice

import "errors"

var (
    ErrInvalidAPIKey     = errors.New("invalid API key")
    ErrServiceUnavailable = errors.New("MyService is unavailable")
    ErrRateLimited       = errors.New("rate limited by MyService")
    ErrInvalidPrompt     = errors.New("invalid prompt")
)
```

### Step 4: Write Tests

**File: `provider/myservice/provider_test.go`**

```go
package myservice_test

import (
    "context"
    "testing"
    "time"
    "github.com/rajveer43/goagentflow/provider/myservice"
)

func TestProviderComplete(t *testing.T) {
    provider := myservice.New("test-key", "http://localhost:8000")
    
    ctx := context.Background()
    response, err := provider.Complete(ctx, "What is AI?")
    
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if response == "" {
        t.Error("expected non-empty response")
    }
}

func TestProviderEmptyPrompt(t *testing.T) {
    provider := myservice.New("test-key", "http://localhost:8000")
    
    ctx := context.Background()
    _, err := provider.Complete(ctx, "")
    
    if err == nil {
        t.Error("expected error for empty prompt")
    }
}

func TestProviderTimeout(t *testing.T) {
    provider := myservice.New("test-key", "http://localhost:8000")
    provider.SetTimeout(1 * time.Millisecond)
    
    ctx := context.Background()
    _, err := provider.Complete(ctx, "test")
    
    // Should timeout or fail
    if err == nil {
        t.Error("expected timeout error")
    }
}

func BenchmarkComplete(b *testing.B) {
    provider := myservice.New("test-key", "http://localhost:8000")
    ctx := context.Background()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        provider.Complete(ctx, "test")
    }
}
```

### Step 5: Add to Model Registry (Optional)

If you want your models to appear in `provider.ListByProvider()`:

**File: `provider/registry.go`** (update)

```go
func init() {
    // Register MyService models
    models["myservice-1.0"] = &ModelInfo{
        Provider:        "myservice",
        ModelName:       "myservice-1.0",
        ContextWindow:   4096,
        MaxTokens:       2048,
        CostInput:       0.0001,
        CostOutput:      0.0002,
        SupportsStreaming: false,
        SupportsVision: false,
        SupportsTools:  false,
    }
}
```

### Step 6: Create Example

**File: `examples/myservice/main.go`**

```go
package main

import (
    "context"
    "fmt"
    "os"
    
    "github.com/rajveer43/goagentflow/provider/myservice"
)

func main() {
    apiKey := os.Getenv("MYSERVICE_API_KEY")
    if apiKey == "" {
        fmt.Println("Please set MYSERVICE_API_KEY environment variable")
        os.Exit(1)
    }
    
    ctx := context.Background()
    
    // Create provider
    llm := myservice.New(apiKey, "https://api.myservice.com")
    
    // Simple completion
    response, err := llm.Complete(ctx, "What is machine learning?")
    if err != nil {
        panic(err)
    }
    
    fmt.Println("Response:", response)
}
```

### Step 7: Update Documentation

Add to `provider/README.md`:

```markdown
## MyService Provider

Fast, cheap local inference via MyService.

### Setup

```bash
export MYSERVICE_API_KEY="your-key"
go run examples/myservice/main.go
```

### Features

- ✅ Streaming support
- ❌ Vision
- ❌ Function calling

### Pricing

$0.0001 / input token, $0.0002 / output token
```

---

## Custom Memory Backend

### Step 1: Create Package

```bash
mkdir -p memory/custom
touch memory/custom/memory.go
touch memory/custom/memory_test.go
```

### Step 2: Implement runtime.Memory

**File: `memory/custom/memory.go`**

```go
package custom

import (
    "context"
    "fmt"
    "sync"
    
    "github.com/rajveer43/goagentflow/runtime"
)

// Memory stores messages in a custom backend (e.g., Redis, Postgres, DynamoDB).
type Memory struct {
    backend Backend  // Interface to your backend
    mu      sync.RWMutex
}

// Backend is the interface for custom storage backends.
type Backend interface {
    Save(ctx context.Context, key string, value []byte) error
    Load(ctx context.Context, key string) ([]byte, error)
    Delete(ctx context.Context, key string) error
    List(ctx context.Context, prefix string) ([]string, error)
}

// New creates a new custom memory with the given backend.
func New(backend Backend) *Memory {
    return &Memory{backend: backend}
}

// AddMessage adds a message to the store.
func (m *Memory) AddMessage(ctx context.Context, msg runtime.Message) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    // Serialize message
    data, err := msg.MarshalJSON()
    if err != nil {
        return fmt.Errorf("marshal error: %w", err)
    }
    
    // Store with timestamped key
    key := fmt.Sprintf("message:%d", time.Now().UnixNano())
    return m.backend.Save(ctx, key, data)
}

// GetMessages retrieves all messages from the store.
func (m *Memory) GetMessages(ctx context.Context) ([]runtime.Message, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    // List all message keys
    keys, err := m.backend.List(ctx, "message:")
    if err != nil {
        return nil, fmt.Errorf("list error: %w", err)
    }
    
    var messages []runtime.Message
    for _, key := range keys {
        data, err := m.backend.Load(ctx, key)
        if err != nil {
            return nil, fmt.Errorf("load error: %w", err)
        }
        
        var msg runtime.Message
        if err := msg.UnmarshalJSON(data); err != nil {
            return nil, fmt.Errorf("unmarshal error: %w", err)
        }
        
        messages = append(messages, msg)
    }
    
    return messages, nil
}

// Set stores arbitrary key-value data.
func (m *Memory) Set(ctx context.Context, key string, value any) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    
    return m.backend.Save(ctx, fmt.Sprintf("kv:%s", key), data)
}

// Get retrieves arbitrary key-value data.
func (m *Memory) Get(ctx context.Context, key string) (any, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    data, err := m.backend.Load(ctx, fmt.Sprintf("kv:%s", key))
    if err != nil {
        return nil, err
    }
    
    var value any
    if err := json.Unmarshal(data, &value); err != nil {
        return nil, err
    }
    
    return value, nil
}

// Verify implements runtime.Memory
var _ runtime.Memory = (*Memory)(nil)
```

### Step 3: Write Tests

```go
func TestMemoryAddMessage(t *testing.T) {
    backend := NewMockBackend()
    mem := custom.New(backend)
    
    ctx := context.Background()
    msg := runtime.Message{Role: "user", Content: "hello"}
    
    err := mem.AddMessage(ctx, msg)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}

func TestMemoryGetMessages(t *testing.T) {
    backend := NewMockBackend()
    mem := custom.New(backend)
    
    ctx := context.Background()
    
    // Add messages
    mem.AddMessage(ctx, runtime.Message{Role: "user", Content: "hello"})
    mem.AddMessage(ctx, runtime.Message{Role: "assistant", Content: "hi"})
    
    // Retrieve
    messages, err := mem.GetMessages(ctx)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if len(messages) != 2 {
        t.Errorf("expected 2 messages, got %d", len(messages))
    }
}
```

---

## Custom Chain

### Step 1: Create Package

```bash
mkdir -p chains/custom
touch chains/custom/chain.go
touch chains/custom/chain_test.go
```

### Step 2: Implement runtime.Chain

**File: `chains/custom/chain.go`**

```go
package custom

import (
    "context"
    "fmt"
    
    "github.com/rajveer43/goagentflow/runtime"
)

// MyChain processes documents through a custom pipeline.
type MyChain struct {
    llm  runtime.LLM
    name string
}

// New creates a new custom chain.
func New(llm runtime.LLM) *MyChain {
    return &MyChain{
        llm:  llm,
        name: "custom-chain",
    }
}

// Run executes the chain on input.
func (c *MyChain) Run(ctx context.Context, input any) (any, error) {
    // 1. Type assert input
    docs, ok := input.([]runtime.Document)
    if !ok {
        return nil, fmt.Errorf("expected []Document, got %T", input)
    }
    
    // 2. Process documents
    var results []string
    for _, doc := range docs {
        // Your custom logic here
        result, err := c.processDocument(ctx, doc)
        if err != nil {
            return nil, err
        }
        results = append(results, result)
    }
    
    // 3. Return results
    return results, nil
}

// processDocument implements custom logic for a single document.
func (c *MyChain) processDocument(ctx context.Context, doc runtime.Document) (string, error) {
    prompt := fmt.Sprintf("Analyze this document:\n\n%s", doc.Content)
    return c.llm.Complete(ctx, prompt)
}

// Verify implements runtime.Chain
var _ runtime.Chain = (*MyChain)(nil)
```

### Step 3: Write Tests

```go
func TestMyChainRun(t *testing.T) {
    llm := &MockLLM{Response: "Analysis result"}
    chain := custom.New(llm)
    
    ctx := context.Background()
    docs := []runtime.Document{
        {Content: "test document"},
    }
    
    result, err := chain.Run(ctx, docs)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if result == nil {
        t.Error("expected non-nil result")
    }
}
```

---

## Custom VectorStore

Implement `runtime.VectorStore`:

```go
type MyVectorStore struct {
    // Your storage backend
}

func (vs *MyVectorStore) Add(ctx context.Context, docs []runtime.Document) error {
    // Store documents with embeddings
    return nil
}

func (vs *MyVectorStore) Search(ctx context.Context, query []float32, k int) ([]runtime.SearchResult, error) {
    // Find k nearest neighbors to query vector
    return nil, nil
}

var _ runtime.VectorStore = (*MyVectorStore)(nil)
```

---

## Testing Custom Components

### Mock Implementations

Create mock implementations for testing:

```go
// MockLLM implements runtime.LLM for testing
type MockLLM struct {
    Response string
    Err      error
    Calls    int
}

func (m *MockLLM) Complete(ctx context.Context, prompt string) (string, error) {
    m.Calls++
    return m.Response, m.Err
}

// Use in tests
func TestWithMockLLM(t *testing.T) {
    llm := &MockLLM{Response: "test response"}
    chain := custom.New(llm)
    
    result, _ := chain.Run(context.Background(), input)
    
    if llm.Calls != 1 {
        t.Errorf("expected 1 call, got %d", llm.Calls)
    }
}
```

### Integration Tests

Test with real components:

```go
func TestIntegration(t *testing.T) {
    // Use real LLM provider
    llm := anthropic.New(os.Getenv("ANTHROPIC_API_KEY"), "claude-opus-4-6")
    
    // Use real memory
    mem := inmemory.New()
    
    // Test integration
    agent := agent.New(llm, mem)
    response, _ := agent.Run(context.Background(), "test")
    
    if response == "" {
        t.Error("expected response")
    }
}
```

---

## Publishing Your Extension

### Option 1: Separate Repository

Create a standalone package users can import:

```bash
mkdir myservice-provider
cd myservice-provider
git init
```

**go.mod:**
```go
module github.com/yourname/myservice-provider

require github.com/rajveer43/goagentflow v1.6.0
```

**Usage:**
```bash
go get github.com/yourname/myservice-provider
```

```go
import "github.com/yourname/myservice-provider"

llm := myservice.New(apiKey)
```

### Option 2: Contribute to goagentflow

Submit a PR to the main repository!

---

## Best Practices

1. **Always implement interfaces completely** - Don't skip methods
2. **Add comprehensive tests** - Unit + integration
3. **Document with examples** - Especially for custom components
4. **Error handling** - Explicit, checked errors
5. **Performance** - Consider benchmarks
6. **Type safety** - Use strong types, avoid `any` where possible
7. **Thread safety** - Use `sync.RWMutex` where needed
8. **Context handling** - Always respect `context.Context`

---

## Examples

See actual implementations:
- LLM: `provider/anthropic/`, `provider/openai/`
- Memory: `memory/entity/`, `memory/window/`
- Chain: `chains/qa/`, `chains/summarization/`
- VectorStore: `vectorstore/memory/`

---

## Questions?

- 💬 Start a [discussion](https://github.com/rajveer43/goagentflow/discussions)
- 🐛 [Open an issue](https://github.com/rajveer43/goagentflow/issues)
- 📖 Check [ARCHITECTURE.md](ARCHITECTURE.md) for design patterns
