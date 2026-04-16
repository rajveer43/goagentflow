package runtime

import (
	"context"

	"github.com/rajveer43/goagentflow/types"
)

// SearchResult represents a document returned from vector search with similarity score.
type SearchResult struct {
	Document types.Document
	Score    float32 // cosine similarity, typically 0.0 to 1.0
	ID       string  // unique identifier in the vector store
}

// VectorStore manages storing and retrieving documents by semantic similarity.
// Pattern: Repository - abstract away vector storage implementation
type VectorStore interface {
	// Add indexes documents with their embeddings.
	// Returns the IDs assigned to each document, in the same order as input.
	Add(ctx context.Context, docs []types.Document, embeddings [][]float32) ([]string, error)

	// Search finds the top-k documents most similar to the query embedding.
	// Returns results sorted by similarity (highest first).
	Search(ctx context.Context, embedding []float32, k int) ([]SearchResult, error)

	// Delete removes documents by their IDs.
	Delete(ctx context.Context, ids []string) error

	// Clear removes all documents from the store.
	Clear(ctx context.Context) error

	// Size returns the number of documents currently stored.
	Size(ctx context.Context) (int, error)
}
