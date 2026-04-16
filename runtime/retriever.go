package runtime

import (
	"context"

	"github.com/rajveer43/goagentflow/types"
)

// Retriever searches for documents relevant to a query.
// Pattern: Strategy - swap between different retrieval implementations
type Retriever interface {
	// Retrieve returns the top-k documents most relevant to the query.
	Retrieve(ctx context.Context, query string, k int) ([]types.Document, error)
}
