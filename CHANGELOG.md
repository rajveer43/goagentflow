# Changelog

All notable changes to goagentflow will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Planned
- GoQL query language for agent workflows
- Distributed agent execution (multi-agent coordination)
- Persistent streaming with checkpoint/resume
- Visual agent graph editor
- Integration with more vector databases (Weaviate, Milvus)

---

## [1.6.0] - 2026-04-17

### Added

- **Pre-built Chains (Phase 4)**
  - `QA Chain` - Document Q&A with retriever
  - `Summarization Chain` - Multi-document summarization with Stuff and MapReduce strategies
  - `SQL Chain` - Natural language to SQL with dialect support
  - `Agent Chain` - Multi-turn agent with tool registration and memory
  - `ChainPipeline` - Compose chains into workflows

- **Examples**
  - `examples/chains/main.go` - Demonstrates all pre-built chains
  - `examples/graph-workflow/main.go` - Complex agent workflow with graphs
  - `examples/web-research-lite/main.go` - RAG-based web research agent
  - `examples/planner-executor/main.go` - Planning and execution pattern

- **Documentation**
  - `docs/ARCHITECTURE.md` - High-level design overview
  - `docs/EXTENDING.md` - Extension guide with templates
  - Package-level documentation in `doc.go`

### Changed

- Memory interface now supports `Set()` and `Get()` for arbitrary state
- Tool registry accepts variadic tool definitions
- Chain interface signature remains stable but implementations expanded

### Fixed

- Entity extraction handles nested entities correctly
- Summary memory respects compression thresholds
- Streaming chains properly handle EOF

---

## [1.5.0] - 2026-04-10

### Added

- **Advanced Memory Types (Phase 3)**
  - `Entity Memory` - Extracts and tracks named entities (people, places, concepts)
  - `Summary Memory` - Auto-summarizes old messages with LLM
  - `Compressive Memory` - LLM-based message compression with thresholds
  - `Buffer Memory` - FIFO sliding window (keep last N messages)
  - `Window Memory` - Token-aware sliding window (respects token budgets)

- **Memory Composition**
  - Decorator pattern for stacking memory backends
  - Example: `inmemory → entity tracking → window management`
  - Composable memory chains for complex scenarios

- **LLM-based Memory Operations**
  - `LLMSummarizer` - Uses LLM to summarize conversations
  - Entity extraction pipeline
  - Message compression strategies

### Changed

- Memory interface extended with `Clear()` and `Size()` methods
- Message type now includes metadata (e.g., entities, tokens)

### Fixed

- Memory backends now properly handle concurrent access
- Token counting in Window memory is more accurate
- Entity extraction works with multi-sentence inputs

---

## [1.4.0] - 2026-04-01

### Added

- **Multiple LLM Providers (Phase 2)**
  - Anthropic Claude (3.5 Sonnet, Opus 4.6, 3 Haiku)
  - OpenAI GPT (GPT-4, GPT-4 Turbo, GPT-3.5 Turbo)
  - Google Gemini (Gemini 2.0 Flash, Gemini 1.5 Pro)
  - Ollama (local LLMs: llama2, mistral, phi, neural-chat, etc.)
  - Mistral AI (Mistral Large, Mistral Small)
  - Groq (fast inference: llama-3.3-70b, mixtral-8x7b)
  - Cohere (Command R+, Command R)

- **Model Registry**
  - Centralized model catalog with 25+ models
  - Model metadata (cost, capabilities, context window)
  - Query by provider, capability, or exact model
  - Functions: `GetModel()`, `ListByProvider()`, `ListCapable()`

- **Provider Features**
  - Token counting (for cost estimation)
  - Streaming support (where available)
  - Vision capabilities (GPT-4V, Claude 3.5, Gemini Vision)
  - Function calling support
  - Rate limiting and retry logic

- **Examples**
  - `examples/providers/main.go` - All 7 providers in action

### Changed

- LLM interface remains stable, new providers added
- Provider registry accessible via `provider.GetModel()`, `provider.ListByProvider()`

### Fixed

- Provider timeout handling improved
- API error messages are more descriptive

---

## [1.3.0] - 2026-03-25

### Added

- **Vector Embeddings & RAG (Phase 1)**
  - Embedding providers:
    - OpenAI (text-embedding-3-small, text-embedding-3-large)
    - Cohere
    - HuggingFace (local and hosted)

  - Vector stores:
    - In-memory vector store with cosine similarity (pure Go, no external deps)
    - Pinecone integration (stub, ready for API)
    - Chroma integration (stub, ready for API)

  - Document processing:
    - CharacterSplitter (fixed-size chunks)
    - RecursiveSplitter (respects boundaries)
    - Configurable chunk size and overlap

  - RAG chains:
    - RetrieverChain - Embed and search documents
    - RAGChain - Full retrieval + LLM generation
    - ContextualRAGChain - Filter by metadata
    - StreamingRAGChain - Token-by-token responses
    - RerankerChain - Maximal Marginal Relevance (MMR) reranking

  - Document loaders:
    - Text loader
    - CSV loader
    - PDF loader
    - URL loader (fetch HTML)
    - HTML loader

- **Comprehensive Examples**
  - `examples/rag/main.go` - Full RAG pipeline

### Changed

- Initial stable release of core interfaces

### Fixed

- Cosine similarity computation handles edge cases (zero vectors)
- PDF extraction properly handles unicode characters

---

## [1.0.0] - 2026-03-15

### Added

- Initial release of goagentflow
- Core interfaces: LLM, Memory, Chain, VectorStore, Retriever, Embedder
- Basic in-memory implementations
- Project structure and documentation

---

## Versioning

- **Major (X.0.0)** - Breaking API changes
- **Minor (1.Y.0)** - New features without breaking changes
- **Patch (1.0.Z)** - Bug fixes only

---

## Migration Guides

### v1.4.0 → v1.5.0
No breaking changes. New memory types are opt-in.

### v1.3.0 → v1.4.0
No breaking changes. Multiple providers now available via `provider.GetModel()`.

### v1.0.0 → v1.3.0
No breaking changes. New RAG features are additive.

---

## Future Roadmap

### Phase 5: Multi-Agent Orchestration (v1.7.0)
- Agent graph execution
- Inter-agent communication
- Consensus mechanisms

### Phase 6: Structured Output & Validation (v1.8.0)
- JSON schema generation
- Structured output parsing
- Pydantic-like validation

### Phase 7: Observability & Monitoring (v1.9.0)
- Enhanced tracing with OpenTelemetry
- Cost tracking per operation
- Performance metrics

### Phase 8: Production Hardening (v2.0.0)
- Connection pooling
- Circuit breakers
- Comprehensive error recovery

---

For questions or suggestions, please [open an issue](https://github.com/rajveer43/goagentflow/issues).
