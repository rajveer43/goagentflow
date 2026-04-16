# goagentflow

An idiomatic, production-ready Go framework for building AI agents with LLM providers, retrieval-augmented generation (RAG), advanced memory management, and composable chains.

## Overview

`goagentflow` provides a complete toolkit for building sophisticated AI applications in Go, combining:
- **Multiple LLM Providers** (Anthropic, OpenAI, Gemini, Ollama, and more)
- **Vector Embeddings & RAG** (semantic search + LLM generation)
- **Advanced Memory Types** (buffer, summary, entity tracking, token-aware)
- **Pre-built Chains** (QA, summarization, SQL generation, agent orchestration)
- **Interface-First Architecture** (composable, extensible components)

## Philosophy

- **Idiomatic Go**: Structs, interfaces, and explicit control flow
- **No Magic**: All behavior is transparent and testable
- **Composable**: Components combine into pipelines and chains
- **Pluggable**: Swap implementations (memory backends, LLMs, vector stores)
- **Minimal Dependencies**: Core has no external dependencies beyond what's necessary

## Installation

```bash
go get github.com/rajveer43/goagentflow
```

## Quick Start

### Simple LLM Call

```go
package main

import (
    "context"
    "fmt"
    "github.com/rajveer43/goagentflow/provider/anthropic"
)

func main() {
    ctx := context.Background()
    llm := anthropic.New("your-api-key", "claude-opus-4-6")
    
    response, err := llm.Complete(ctx, "What is machine learning?")
    if err != nil {
        panic(err)
    }
    
    fmt.Println(response)
}
```

### RAG Pipeline

```go
// Create embedder and vector store
embedder := openai.New("your-api-key", "text-embedding-3-small")
vectorStore := memory.New(embedder)

// Create retriever
retriever := &retrieval.Retriever{
    VectorStore: vectorStore,
    K:           3,
}

// Create RAG chain
ragChain := retrieval.NewRAGChain(llm, retriever)
result, err := ragChain.Run(ctx, "Your question here")
```

### Conversational Agent with Memory

```go
// Create memory and agent
memory := inmemory.New()
agent := agent.New(llm, memory)

// First message
result1, _ := agent.Run(ctx, "I'm learning about AI")
fmt.Println(result1)

// Follow-up (context from memory)
result2, _ := agent.Run(ctx, "What are the key concepts?")
fmt.Println(result2)
```

## Features by Phase

### Phase 1: Embeddings & RAG (v1.3.0)

**Embeddings**
- OpenAI embeddings (text-embedding-3-small, text-embedding-3-large)
- Cohere embeddings API
- HuggingFace embeddings (local or hosted)

**Vector Stores**
- In-memory vector store with cosine similarity (pure Go, no external deps)
- Pinecone integration (stub ready for API)
- Chroma integration (stub ready for API)

**Document Processing**
- Character-based text splitting
- Recursive text splitting (respects boundaries)
- Configurable chunk size and overlap

**RAG Chains**
- `RetrieverChain`: Embed and search documents
- `RAGChain`: Full retrieval + LLM generation pipeline
- `ContextualRAGChain`: Filter by metadata before generation
- `StreamingRAGChain`: Stream token-by-token responses
- `RerankerChain`: Rerank results with MMR (maximal marginal relevance)

### Phase 2: LLM Providers (v1.4.0)

**7 LLM Providers** (25+ models total)
- **Anthropic**: Claude 3.5 Sonnet, Claude Opus 4.6, Claude 3 Haiku
- **OpenAI**: GPT-4, GPT-4 Turbo, GPT-3.5 Turbo
- **Google Gemini**: Gemini 2.0 Flash, Gemini 1.5 Pro
- **Ollama**: Run local LLMs (llama2, mistral, phi, etc.)
- **Mistral**: Mistral Large, Mistral Small
- **Groq**: Fast inference (llama-3.3-70b, mixtral-8x7b)
- **Cohere**: Command R+, Command R

**Model Registry**
- Centralized model catalog with metadata
- Cost tracking (input/output tokens)
- Capability filtering (streaming, vision, function calling)
- `GetModel()`, `ListByProvider()`, `ListCapable()` functions

**Optional Interfaces**
- `TokenCounter`: Count tokens before API calls
- `ModelInfoProvider`: Query model capabilities

### Phase 3: Advanced Memory Types (v1.5.0)

**6 Memory Backends** (all implement `runtime.Memory`)

1. **InMemory** (v1.0.0)
   - Simple map-based storage
   - Fast, thread-safe
   - Best for: Single-session apps, testing

2. **Buffer Memory**
   - FIFO sliding window (keep last N messages)
   - Automatic old message eviction
   - Best for: Short conversations with bounded memory

3. **Window Memory**
   - Token-aware sliding window
   - Approximate token counting
   - Best for: Long conversations with cost sensitivity

4. **Entity Memory**
   - Tracks named entities (people, places, concepts)
   - Extracts entities from messages
   - Injects entity summaries into context
   - Best for: Multi-turn conversations needing entity recall

5. **Summary Memory**
   - Auto-summarizes old messages using LLM
   - Wraps inner memory (decorator pattern)
   - Best for: Very long conversations with semantic compression

6. **Compressive Memory**
   - LLM-based message compression
   - Compresses when threshold exceeded
   - Best for: Conversations needing intelligent compression

**Memory Composition**
All memory types are composable:
```go
// Create a stack: inmemory -> entity tracking -> window management
base := inmemory.New()
withEntity := entity.New(base)
withWindow := window.New(withEntity, 4096) // 4096 token budget
```

### Phase 4: Pre-Built Chains (v1.6.0)

All chains implement `runtime.Chain` and can be composed with pipelines.

1. **QA Chain** (`chains/qa`)
   - Input: question string
   - Output: answer + sources
   - Uses: Retriever + LLM
   - Best for: Document Q&A systems

2. **Summarization Chain** (`chains/summarization`)
   - Input: text or documents
   - Output: summary string
   - Strategies: 
     - `StuffStrategy`: Concatenate + summarize once (fast)
     - `MapReduceStrategy`: Summarize each chunk, then combine (handles large docs)
   - Best for: Document/article summarization

3. **SQL Chain** (`chains/sql`)
   - Input: natural language question
   - Output: SQL query
   - Uses: LLM + schema context
   - Configurable SQL dialect (PostgreSQL, MySQL, SQLite)
   - Best for: Natural language database queries

4. **Agent Chain** (`chains/agent`)
   - Input: user message string
   - Output: agent response string
   - Combines: LLM + tools + memory + retriever
   - Features: Conversation history, tool registration, context building
   - Best for: Multi-turn agent interactions

**Chain Composition**
Compose chains into pipelines:
```go
summarizer := summarychain.New(llm, summarychain.StuffStrategy)
qaChain := qachain.New(retriever, llm, 3)
// Pipeline summarizes input, then answers questions on summary
pipeline := runtime.NewChainPipeline(summarizer, qaChain)
```

## Architecture

### Core Interfaces (`runtime/`)

```go
// LLM: Any language model
type LLM interface {
    Complete(ctx context.Context, prompt string) (string, error)
}

// Embedder: Text to vector
type Embedder interface {
    Embed(ctx context.Context, text string) ([]float32, error)
}

// VectorStore: Semantic search
type VectorStore interface {
    Add(ctx context.Context, docs []Document) error
    Search(ctx context.Context, query []float32, k int) ([]SearchResult, error)
}

// Memory: Conversation state
type Memory interface {
    AddMessage(ctx context.Context, msg Message) error
    GetMessages(ctx context.Context) ([]Message, error)
    Set(ctx context.Context, key string, value any) error
    Get(ctx context.Context, key string) (any, error)
}

// Chain: Composable pipeline step
type Chain interface {
    Run(ctx context.Context, input any) (any, error)
}

// Retriever: Document retrieval
type Retriever interface {
    Retrieve(ctx context.Context, query string, k int) ([]Document, error)
}
```

### Directory Structure

```
├── runtime/           # Core interfaces (LLM, Memory, Chain, VectorStore, etc.)
├── provider/          # LLM implementations (anthropic, openai, gemini, ollama, etc.)
├── embeddings/        # Embedding providers (openai, cohere, huggingface)
├── vectorstore/       # Vector stores (memory, pinecone, chroma)
├── memory/            # Memory backends (inmemory, buffer, window, entity, etc.)
├── chains/            # Pre-built chains (qa, summarization, sql, agent)
├── retrieval/         # RAG pipelines and chains
├── loader/            # Document loading and splitting
├── types/             # Unified type definitions
├── examples/          # Runnable examples
└── internal/          # Internal helpers
```

## Examples

Run the examples to see everything in action:

```bash
# RAG pipeline with embeddings and vector store
go run examples/rag/main.go

# All LLM providers and model registry
go run examples/providers/main.go

# Advanced memory types
go run examples/memory/main.go

# Pre-built chains (QA, summarization, SQL, agent)
go run examples/chains/main.go
```

## Key Design Patterns

- **Interface-First**: All major concepts are interfaces, enabling composition and testing
- **Decorator Pattern**: Memory types can wrap other memory backends (entity → window → inmemory)
- **Strategy Pattern**: Summarization chains support multiple strategies (stuff vs. map-reduce)
- **Composition**: Chains combine via `NewChainPipeline()` into complex workflows
- **Adapter Pattern**: LLM providers adapt different APIs to unified `runtime.LLM` interface

## Configuration

### LLM Provider Selection

```go
// Use model registry to find the right model
models := provider.ListCapable("streaming")  // Only streaming models
bestModel := provider.GetModel("claude-opus-4-6")

// Or instantiate directly
llm := anthropic.New(apiKey, "claude-opus-4-6")
```

### Memory Configuration

```go
// Buffer: keep last 5 messages
bufMem := buffer.New(5)

// Window: keep ~4096 tokens
winMem := window.New(4096)

// Entity: track entities with base memory
entityMem := entity.New(inmemory.New())

// Summary: auto-summarize after 3 messages
summaryMem := summary.New(inmemory.New(), llmSummarizer, 3)
```

### Chain Configuration

```go
// Summarization with custom chunk size
summarizer := summarychain.New(llm, summarychain.MapReduceStrategy)
summarizer.SetChunkSize(500, 50)  // 500 char chunks, 50 char overlap

// SQL with specific dialect
sqlChain := sqlchain.New(llm, schema)
sqlChain.SetDialect("PostgreSQL")

// Agent with tool registration
agent := agent.New(llm, memory)
agent.RegisterTool(myTool)
agent.SetMaxSteps(5)
```

## Requirements

- Go 1.18+
- No external dependencies in core (`runtime/`, `memory/`, `loader/`, `retrieval/`)
- Provider-specific packages have dependencies on their respective SDKs

## Development

Build and test:
```bash
go build ./...
go test ./...
```

## License

MIT

## Contributing

Contributions welcome! Please ensure:
- All tests pass (`go test ./...`)
- Code builds cleanly (`go build ./...`)
- New components implement appropriate interfaces
- Examples are provided for new features
