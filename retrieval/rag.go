package retrieval

import (
	"context"
	"fmt"
	"strings"

	"github.com/rajveer43/goagentflow/runtime"
	"github.com/rajveer43/goagentflow/types"
)

// RAGChain combines retrieval and LLM prompting into a single chain.
// Pattern: Adapter - orchestrates Retriever + LLM into a single runtime.Chain
// Input: query string; Output: LLM response string with context
type RAGChain struct {
	retriever       runtime.Retriever
	llm             runtime.LLM
	k               int    // number of documents to retrieve
	promptTemplate  string // template for combining context and query
	contextSeparator string
}

// NewRAGChain creates a new RAG chain.
// retriever: retrieves relevant documents
// llm: generates answers based on context
// k: number of documents to retrieve
func NewRAGChain(retriever runtime.Retriever, llm runtime.LLM, k int) *RAGChain {
	return &RAGChain{
		retriever:        retriever,
		llm:              llm,
		k:                k,
		promptTemplate:   defaultPromptTemplate,
		contextSeparator: "\n\n---\n\n",
	}
}

// SetPromptTemplate sets a custom prompt template.
// Template should include {context} and {query} placeholders.
func (rc *RAGChain) SetPromptTemplate(template string) {
	rc.promptTemplate = template
}

// Run implements runtime.Chain.
// Input: query string; Output: LLM response string
func (rc *RAGChain) Run(ctx context.Context, input any) (any, error) {
	query, ok := input.(string)
	if !ok {
		return nil, fmt.Errorf("expected string query, got %T", input)
	}

	// Retrieve relevant documents
	docs, err := rc.retriever.Retrieve(ctx, query, rc.k)
	if err != nil {
		return nil, fmt.Errorf("retrieval failed: %w", err)
	}

	// Format documents as context
	context := rc.formatDocuments(docs)

	// Build final prompt
	prompt := strings.NewReplacer(
		"{context}", context,
		"{query}", query,
	).Replace(rc.promptTemplate)

	// Get LLM response
	response, err := rc.llm.Complete(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("llm completion failed: %w", err)
	}

	return response, nil
}

// formatDocuments formats retrieved documents into a context string.
func (rc *RAGChain) formatDocuments(docs []types.Document) string {
	if len(docs) == 0 {
		return ""
	}

	var parts []string
	for i, doc := range docs {
		part := fmt.Sprintf("Document %d:\n%s", i+1, doc.PageContent)
		parts = append(parts, part)
	}

	return strings.Join(parts, rc.contextSeparator)
}

const defaultPromptTemplate = `You are a helpful assistant. Use the following documents to answer the question.

Context:
{context}

Question: {query}

Answer:`

// ContextualRAGChain extends RAGChain with metadata awareness.
// Can filter documents by metadata before using them.
type ContextualRAGChain struct {
	*RAGChain
	metadataFilter func(doc types.Document) bool
}

// NewContextualRAGChain creates a RAG chain with metadata filtering.
func NewContextualRAGChain(retriever runtime.Retriever, llm runtime.LLM, k int) *ContextualRAGChain {
	return &ContextualRAGChain{
		RAGChain: NewRAGChain(retriever, llm, k),
		metadataFilter: func(doc types.Document) bool {
			return true // default: include all
		},
	}
}

// SetMetadataFilter sets a function to filter documents by metadata.
func (crc *ContextualRAGChain) SetMetadataFilter(filter func(doc types.Document) bool) {
	crc.metadataFilter = filter
}

// Run implements runtime.Chain with metadata filtering.
func (crc *ContextualRAGChain) Run(ctx context.Context, input any) (any, error) {
	query, ok := input.(string)
	if !ok {
		return nil, fmt.Errorf("expected string query, got %T", input)
	}

	// Retrieve documents
	docs, err := crc.retriever.Retrieve(ctx, query, crc.k)
	if err != nil {
		return nil, fmt.Errorf("retrieval failed: %w", err)
	}

	// Filter by metadata
	var filtered []types.Document
	for _, doc := range docs {
		if crc.metadataFilter(doc) {
			filtered = append(filtered, doc)
		}
	}

	if len(filtered) == 0 {
		return nil, fmt.Errorf("no documents passed metadata filter")
	}

	// Format and prompt
	context := crc.formatDocuments(filtered)
	prompt := strings.NewReplacer(
		"{context}", context,
		"{query}", query,
	).Replace(crc.promptTemplate)

	response, err := crc.llm.Complete(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("llm completion failed: %w", err)
	}

	return response, nil
}

// StreamingRAGChain streams LLM responses instead of returning full response.
type StreamingRAGChain struct {
	*RAGChain
}

// NewStreamingRAGChain creates a RAG chain with streaming support.
func NewStreamingRAGChain(retriever runtime.Retriever, llm runtime.LLM, k int) *StreamingRAGChain {
	return &StreamingRAGChain{
		RAGChain: NewRAGChain(retriever, llm, k),
	}
}

// Stream retrieves documents and streams LLM response.
// Returns two channels: tokens channel and errors channel.
func (src *StreamingRAGChain) Stream(ctx context.Context, query string) (<-chan string, <-chan error) {
	tokensCh := make(chan string, 10)
	errorsCh := make(chan error, 1)

	go func() {
		defer close(tokensCh)
		defer close(errorsCh)

		// Retrieve documents
		docs, err := src.retriever.Retrieve(ctx, query, src.k)
		if err != nil {
			errorsCh <- fmt.Errorf("retrieval failed: %w", err)
			return
		}

		// Format context and build prompt
		context := src.formatDocuments(docs)
		prompt := strings.NewReplacer(
			"{context}", context,
			"{query}", query,
		).Replace(src.promptTemplate)

		// Stream from LLM
		tokens, errors := src.llm.Stream(ctx, prompt)

		for {
			select {
			case token, ok := <-tokens:
				if !ok {
					tokens = nil
					continue
				}
				select {
				case tokensCh <- token:
				case <-ctx.Done():
					errorsCh <- ctx.Err()
					return
				}
			case err, ok := <-errors:
				if !ok {
					return
				}
				errorsCh <- err
				return
			case <-ctx.Done():
				errorsCh <- ctx.Err()
				return
			}

			// Exit when both channels are closed
			if tokens == nil {
				return
			}
		}
	}()

	return tokensCh, errorsCh
}
