# Architecture Overview

goagentflow is designed around **transparent, composable interfaces** inspired by OpenClaw's philosophy. This document explains the high-level architecture and design decisions.

---

## 🏗️ Core Design Principles

### 1. **Transparency** (No Magic)
- All behavior is explicit
- No hidden side effects
- Easy to test and reason about
- Error handling is explicit, not silent

### 2. **Modularity** (Composable)
- Everything is an interface
- Implementations are swappable
- No hard dependencies between components
- Decorator pattern for enhancement

### 3. **Hackability** (Extensible)
- Easy to implement custom components
- Clear templates and examples
- Minimal boilerplate
- Rich examples for all patterns

### 4. **Idiomatic Go**
- Structs, interfaces, explicit control flow
- No reflection-based magic
- Standard library where possible
- Minimal external dependencies in core

---

## 🎯 Core Interfaces

All major concepts are interfaces, allowing multiple implementations:

### LLM Interface

```go
type LLM interface {
    Complete(ctx context.Context, prompt string) (string, error)
}
```

**Purpose:** Abstract language model interaction
**Implementations:** 7 providers (Anthropic, OpenAI, Gemini, Ollama, Mistral, Groq, Cohere)
**Why interface?** Easy to swap providers, test with mocks

### Memory Interface

```go
type Memory interface {
    AddMessage(ctx context.Context, msg Message) error
    GetMessages(ctx context.Context) ([]Message, error)
    Set(ctx context.Context, key string, value any) error
    Get(ctx context.Context, key string) (any, error)
}
```

**Purpose:** Manage conversation state
**Implementations:** 6 backends (InMemory, Buffer, Window, Entity, Summary, Compressive)
**Why interface?** Easy to compose (decorator pattern), swap backends

### Chain Interface

```go
type Chain interface {
    Run(ctx context.Context, input any) (any, error)
}
```

**Purpose:** Composable pipeline steps
**Implementations:** QA, Summarization, SQL, Agent, plus custom
**Why interface?** Chains compose into complex workflows via ChainPipeline

### VectorStore Interface

```go
type VectorStore interface {
    Add(ctx context.Context, docs []Document) error
    Search(ctx context.Context, query []float32, k int) ([]SearchResult, error)
}
```

**Purpose:** Semantic search
**Implementations:** Memory (pure Go), Pinecone, Chroma
**Why interface?** Swap between local and cloud backends

### Retriever Interface

```go
type Retriever interface {
    Retrieve(ctx context.Context, query string, k int) ([]Document, error)
}
```

**Purpose:** Document retrieval for RAG
**Implementations:** Simple retriever, contextual retriever, streaming retriever
**Why interface?** Different retrieval strategies

### Embedder Interface

```go
type Embedder interface {
    Embed(ctx context.Context, text string) ([]float32, error)
}
```

**Purpose:** Text to vector conversion
**Implementations:** OpenAI, Cohere, HuggingFace
**Why interface?** Easy to swap embedding models

---

## 📦 Package Organization

### `runtime/`
Core interfaces and types. **Zero external dependencies.**

```
runtime/
├── llm.go           # LLM interface
├── memory.go        # Memory interface
├── chain.go         # Chain interface
├── vectorstore.go   # VectorStore interface
├── retriever.go     # Retriever interface
├── embedder.go      # Embedder interface
├── errors.go        # Common error types
├── tool.go          # Tool definition
├── tool_registry.go # Tool registry
└── agent.go         # Agent orchestration
```

**Design:** Public interfaces, concrete types in implementation packages

### `provider/`
LLM provider implementations.

```
provider/
├── anthropic/       # Anthropic Claude
├── openai/          # OpenAI GPT
├── gemini/          # Google Gemini
├── ollama/          # Ollama (local)
├── mistral/         # Mistral
├── groq/            # Groq
├── cohere/          # Cohere
├── registry.go      # Model registry
└── cache/           # Optional caching layer
```

**Design:** Each provider is independent package
**Pattern:** `New(apiKey string, modelName string) LLM`

### `memory/`
Memory backend implementations.

```
memory/
├── inmemory/        # Simple map-based
├── buffer/          # FIFO sliding window
├── window/          # Token-aware window
├── entity/          # Entity tracking
├── summary/         # Auto-summarization
├── compressive/     # LLM compression
└── README.md        # How to add custom
```

**Design:** Each backend implements `Memory` interface
**Pattern:** Decorator pattern for composition

### `chains/`
Pre-built chain implementations.

```
chains/
├── qa/              # Question answering
├── summarization/   # Document summarization
├── sql/             # Natural language to SQL
├── agent/           # Multi-turn agent
└── README.md        # How to add custom
```

**Design:** Each chain implements `Chain` interface
**Pattern:** Strategy pattern for different approaches (e.g., MapReduce in summarization)

### `retrieval/`
RAG pipeline implementations.

```
retrieval/
├── chain.go         # Simple RAG chain
├── rag.go           # Full RAG pipeline
├── reranker.go      # Result reranking
└── streaming.go     # Streaming RAG
```

**Design:** Composition of Embedder + VectorStore + LLM

### `embeddings/`
Embedding provider implementations.

```
embeddings/
├── openai/          # OpenAI embeddings
├── cohere/          # Cohere embeddings
└── huggingface/     # HuggingFace embeddings
```

### `vectorstore/`
Vector store implementations.

```
vectorstore/
├── memory/          # Pure Go, in-memory
├── pinecone/        # Pinecone cloud
├── chroma/          # Chroma server
```

### `loader/`
Document loading and processing. **Zero external dependencies in core.**

```
loader/
├── loader.go        # Base loader interface
├── text.go          # Text file loader
├── csv.go           # CSV loader
├── pdf.go           # PDF loader
├── html.go          # HTML loader
├── url.go           # URL fetcher
└── splitter.go      # Text splitting strategies
```

**Design:** Chainable loaders, text splitting strategies

### `types/`
Unified type definitions used across packages.

```go
type Document struct {
    Content  string
    Metadata map[string]any
}

type Message struct {
    Role    string
    Content string
}

type SearchResult struct {
    Document Document
    Score    float32
}
```

### `observer/`
Observability: metrics, logging, tracing.

```
observer/
├── logging/         # Structured logging
├── metrics/         # Metrics collection
└── tracing/         # OpenTelemetry tracing
```

### `internal/`
Internal helpers (not part of public API).

```
internal/
├── backoff/         # Exponential backoff
├── stream/          # Streaming utilities
├── idempotency/     # Idempotency keys
└── logger/          # Logger wrapper
```

---

## 🔄 Common Patterns

### 1. **Provider Pattern** (Strategy)

Different implementations of same interface:

```go
// All implement runtime.LLM
llm1 := anthropic.New(key1, "claude-opus-4-6")
llm2 := openai.New(key2, "gpt-4")
llm3 := ollama.New("localhost:11434", "mistral")

// Use identically
response, _ := llm1.Complete(ctx, prompt)
response, _ := llm2.Complete(ctx, prompt)
response, _ := llm3.Complete(ctx, prompt)
```

### 2. **Decorator Pattern** (Composition)

Wrap one implementation with another:

```go
// Chain memory implementations
base := inmemory.New()
withEntity := entity.New(base)           // Track entities
withWindow := window.New(withEntity, 4096)  // Token window
withSummary := summary.New(withWindow, llm) // Auto-summarize

// Use composed memory
agent := agent.New(llm, withSummary)
```

### 3. **Strategy Pattern** (Behavior Variation)

Multiple strategies for same task:

```go
// Summarization supports different strategies
summarizer := summarychain.New(llm, summarychain.StuffStrategy)      // Fast
summarizer := summarychain.New(llm, summarychain.MapReduceStrategy)  // Handles large docs
```

### 4. **Pipeline Pattern** (Composition)

Compose chains into workflows:

```go
splitter := loader.NewCharacterSplitter(500, 50)
summarizer := summarychain.New(llm, summarychain.StuffStrategy)
qaChain := qachain.New(retriever, llm, 3)

// Pipeline: Load → Split → Summarize → Answer
pipeline := runtime.NewChainPipeline(
    splitter,
    summarizer,
    qaChain,
)

result, _ := pipeline.Run(ctx, documents)
```

### 5. **Adapter Pattern** (Bridge APIs)

Adapt different APIs to unified interface:

```go
// Different APIs, unified interface
anthropicLLM := anthropic.New(apiKey, modelName)   // Anthropic API
openaiLLM := openai.New(apiKey, modelName)         // OpenAI API
geminLLM := gemini.New(apiKey, modelName)          // Google API

// All implement runtime.LLM, so all work the same way
```

---

## 📊 Data Flow Example: RAG Pipeline

```
Document Input
    ↓
[Loader] - Load & split documents
    ↓
[Embedder] - Convert text to vectors
    ↓
[VectorStore] - Store embeddings
    ↓
User Query
    ↓
[Embedder] - Convert query to vector
    ↓
[VectorStore.Search] - Find similar documents
    ↓
[LLM] - Generate answer with context
    ↓
Response
```

**Component independence:** Each component is replaceable

---

## 🧠 Memory Architecture

Memory is composed via decorator pattern:

```
User Message → Agent → Memory → LLM Context
                         ↑
                    [Summary Memory]
                         ↑
                    [Entity Memory]
                         ↑
                    [Window Memory]
                         ↑
                    [InMemory Storage]
```

Each layer adds behavior:
- **InMemory**: Stores messages
- **Window**: Respects token budget
- **Entity**: Extracts named entities
- **Summary**: Summarizes old messages

---

## 🔌 Extensibility Points

### Add Custom LLM Provider

Implement `runtime.LLM`:

```go
type MyProvider struct { /* ... */ }
func (p *MyProvider) Complete(ctx context.Context, prompt string) (string, error) { /* ... */ }
```

See `provider/README.md` for template.

### Add Custom Memory

Implement `runtime.Memory`:

```go
type MyMemory struct { /* ... */ }
func (m *MyMemory) AddMessage(ctx context.Context, msg Message) error { /* ... */ }
func (m *MyMemory) GetMessages(ctx context.Context) ([]Message, error) { /* ... */ }
func (m *MyMemory) Set(ctx context.Context, key string, value any) error { /* ... */ }
func (m *MyMemory) Get(ctx context.Context, key string) (any, error) { /* ... */ }
```

See `memory/README.md` for template.

### Add Custom Chain

Implement `runtime.Chain`:

```go
type MyChain struct { /* ... */ }
func (c *MyChain) Run(ctx context.Context, input any) (any, error) { /* ... */ }
```

See `chains/README.md` for template.

---

## 🎯 Error Handling

Explicit, checked error handling (no panic/silent failures):

```go
// ✅ Good: Errors are explicit and tested
response, err := llm.Complete(ctx, prompt)
if err != nil {
    return fmt.Errorf("llm failed: %w", err)
}

// ❌ Bad: Error ignored
response, _ := llm.Complete(ctx, prompt)

// ❌ Bad: Error causes panic
response := must(llm.Complete(ctx, prompt))
```

**Error types in `runtime/errors.go`:**
- `ErrEmptyPrompt`
- `ErrContextCanceled`
- `ErrAPIError`
- `ErrInvalidModel`
- etc.

---

## 🏥 Testing Strategy

### Unit Tests
Test individual components in isolation:

```go
func TestProviderComplete(t *testing.T) {
    provider := anthropic.New(apiKey, "claude-opus-4-6")
    response, err := provider.Complete(ctx, "test")
    // Assert
}
```

### Integration Tests
Test component interactions:

```go
func TestRAGPipeline(t *testing.T) {
    // Load documents
    // Embed them
    // Search
    // Generate answer
    // Assert end-to-end result
}
```

### Mock Implementations
All interfaces have mock implementations for testing:

```go
type MockLLM struct {
    Response string
    Err      error
}

func (m *MockLLM) Complete(ctx context.Context, prompt string) (string, error) {
    return m.Response, m.Err
}
```

---

## 🚀 Performance Considerations

### Memory
- **InMemory**: O(n) space, fast
- **Buffer**: O(k) space (bounded window)
- **Window**: O(tokens) space
- **Entity**: O(entities) space
- **Summary**: O(summaries) space

### Embedding
- **Local (HuggingFace)**: Slow, free, private
- **API (OpenAI, Cohere)**: Fast, costs, network
- **Cache results** to avoid re-embedding

### LLM
- **Token counting**: Estimate cost before calling
- **Caching**: Reuse responses for same prompt
- **Batching**: Multiple requests in single API call (if provider supports)
- **Streaming**: Use for long responses to improve perceived latency

### Vector Search
- **Memory VectorStore**: Fast for small datasets (<100k docs)
- **Pinecone/Chroma**: Required for large datasets (millions of docs)

---

## 🔐 Security

### API Keys
- Never hardcode API keys
- Use environment variables or secure config
- Example: `os.Getenv("ANTHROPIC_API_KEY")`

### Input Validation
- All user inputs should be validated
- Sanitize prompts before sending to LLM
- Validate API responses before using

### Output Safety
- LLMs can generate harmful content
- Consider output filtering in production
- Monitor for prompt injection attacks

---

## 📈 Future Architecture

### Phase 5: Multi-Agent Orchestration
- Agent graph execution
- Inter-agent communication patterns
- Consensus mechanisms

### Phase 6: Structured Output
- JSON schema generation
- Output validation
- Pydantic-like constraints

### Phase 7: Observability
- Enhanced tracing
- Cost tracking
- Performance metrics

### Phase 8: Production Hardening
- Connection pooling
- Circuit breakers
- Comprehensive error recovery

---

## References

- [Design Patterns Used](#-common-patterns)
- [Extending Guide](EXTENDING.md)
- [Contributing Guide](../CONTRIBUTING.md)
- [API Reference](API.md)
