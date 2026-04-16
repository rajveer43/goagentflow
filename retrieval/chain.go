package retrieval

import (
	"context"
	"fmt"

	"github.com/rajveer43/goagentflow/runtime"
	"github.com/rajveer43/goagentflow/types"
)

// RetrieverChain implements runtime.Retriever using an Embedder and VectorStore.
// Pattern: Adapter - bridges Embedder + VectorStore -> Retriever interface
type RetrieverChain struct {
	vectorStore runtime.VectorStore
	embedder    runtime.Embedder
}

// New creates a new retriever chain.
// vectorStore: where to search for documents
// embedder: converts query text to embedding
func New(vectorStore runtime.VectorStore, embedder runtime.Embedder) *RetrieverChain {
	return &RetrieverChain{
		vectorStore: vectorStore,
		embedder:    embedder,
	}
}

// Retrieve returns the top-k documents most relevant to the query.
// Implements runtime.Retriever.
func (rc *RetrieverChain) Retrieve(ctx context.Context, query string, k int) ([]types.Document, error) {
	if query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	if k <= 0 {
		return nil, fmt.Errorf("k must be positive, got %d", k)
	}

	// Embed the query
	embedding, err := rc.embedder.Embed(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	// Search the vector store
	results, err := rc.vectorStore.Search(ctx, embedding, k)
	if err != nil {
		return nil, fmt.Errorf("failed to search vector store: %w", err)
	}

	// Extract documents
	docs := make([]types.Document, len(results))
	for i, result := range results {
		docs[i] = result.Document
	}

	return docs, nil
}

// ChainFunc implements runtime.Retriever as a function.
// Pattern: mirrors runtime.ChainFunc for functional style
type ChainFunc func(ctx context.Context, query string, k int) ([]types.Document, error)

// Retrieve implements runtime.Retriever.
func (f ChainFunc) Retrieve(ctx context.Context, query string, k int) ([]types.Document, error) {
	return f(ctx, query, k)
}

// ChainStep implements runtime.Chain to use a retriever in a pipeline.
// Takes input query string, returns []loader.Document.
type ChainStep struct {
	retriever runtime.Retriever
	k         int // number of documents to retrieve
}

// NewChainStep creates a new chain step for retrieval.
// Can be used in a ChainPipeline for RAG workflows.
func NewChainStep(retriever runtime.Retriever, k int) *ChainStep {
	return &ChainStep{
		retriever: retriever,
		k:         k,
	}
}

// Run implements runtime.Chain.
// Input: query string; Output: []loader.Document
func (cs *ChainStep) Run(ctx context.Context, input any) (any, error) {
	query, ok := input.(string)
	if !ok {
		return nil, fmt.Errorf("expected string query, got %T", input)
	}

	return cs.retriever.Retrieve(ctx, query, cs.k)
}
