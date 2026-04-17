// Package goagentflow is an idiomatic, production-ready Go framework for building AI agents
// with LLM providers, retrieval-augmented generation (RAG), advanced memory management,
// and composable chains.
//
// # Philosophy
//
// goagentflow is inspired by OpenClaw's design principles:
//   - **Transparent**: All behavior is explicit and testable, no magic
//   - **Modular**: Components are composable and easily swappable
//   - **Hackable**: Easy to extend with custom LLMs, memory, and chains
//
// # Quick Start
//
// ## Simple LLM Completion
//
//	package main
//
//	import (
//		"context"
//		"fmt"
//		"github.com/rajveer43/goagentflow/provider/anthropic"
//	)
//
//	func main() {
//		ctx := context.Background()
//		llm := anthropic.New("your-api-key", "claude-opus-4-6")
//		response, _ := llm.Complete(ctx, "What is AI?")
//		fmt.Println(response)
//	}
//
// ## RAG Pipeline with Vector Search
//
//	vectorStore := memory.New(embedder)
//	retriever := retrieval.New(vectorStore, 3)
//	answer, _ := retrieval.RAG(ctx, llm, retriever, "What is machine learning?")
//
// ## Conversational Agent with Memory
//
//	mem := inmemory.New()
//	agent := agent.New(llm, mem)
//	response, _ := agent.Run(ctx, "Tell me about yourself")
//
// # Core Components
//
// The framework is built on interfaces, allowing you to swap implementations:
//
//   - [runtime.LLM] - Language models (Anthropic, OpenAI, Gemini, Ollama, etc.)
//   - [runtime.Memory] - Conversation state (InMemory, Buffer, Window, Entity, Summary, Compressive)
//   - [runtime.Chain] - Composable pipeline steps (QA, Summarization, SQL, Agent)
//   - [runtime.VectorStore] - Semantic search backends (Memory, Pinecone, Chroma)
//   - [runtime.Retriever] - Document retrieval for RAG
//   - [runtime.Embedder] - Text to vector conversion
//
// # Packages
//
//   - [runtime] - Core interfaces and types
//   - [provider] - LLM implementations (anthropic, openai, gemini, ollama, mistral, groq, cohere)
//   - [memory] - Memory backends (inmemory, buffer, window, entity, summary, compressive)
//   - [chains] - Pre-built chains (qa, summarization, sql, agent)
//   - [retrieval] - RAG pipelines and chains
//   - [embeddings] - Embedding providers (openai, cohere, huggingface)
//   - [vectorstore] - Vector store implementations (memory, pinecone, chroma)
//   - [loader] - Document loading and splitting (text, csv, pdf, url, html)
//   - [types] - Unified type definitions
//   - [observer] - Metrics, logging, and tracing
//
// # Examples
//
// Run examples to see all features in action:
//
//	go run examples/rag/main.go                    # RAG pipeline
//	go run examples/providers/main.go              # All LLM providers
//	go run examples/memory/main.go                 # Advanced memory types
//	go run examples/chains/main.go                 # Pre-built chains
//	go run examples/graph-workflow/main.go         # Complex agent graphs
//	go run examples/web-research-lite/main.go      # Web research agent
//	go run examples/planner-executor/main.go       # Planning and execution
//
// # Extending goagentflow
//
// The framework is designed for easy extension. See [docs/EXTENDING.md] for:
//   - How to implement a custom LLM provider
//   - How to implement custom memory
//   - How to create custom chains
//   - Design patterns used throughout
//
// # Documentation
//
//   - [docs/ARCHITECTURE.md] - High-level design overview
//   - [docs/EXTENDING.md] - How to extend with custom components
//   - [docs/API.md] - Full API reference
//   - [provider/README.md] - LLM provider guide
//   - [memory/README.md] - Memory backends guide
//   - [chains/README.md] - Chains guide
//   - [README.md] - Project overview and features
//
// # Getting Started
//
// Install goagentflow:
//
//	go get github.com/rajveer43/goagentflow
//
// Read the [README.md] for comprehensive documentation and examples.
// Visit [pkg.go.dev/github.com/rajveer43/goagentflow] for full API docs.
//
// # License
//
// MIT License - see LICENSE file for details
//
// # Contributing
//
// Contributions are welcome! See [CONTRIBUTING.md] for guidelines.
package goagentflow
