package embeddings

import (
	"context"
	"fmt"

	"github.com/rajveer43/goagentflow/runtime"
)

// OpenAIClient provides embeddings using OpenAI's embedding models.
// Pattern: Strategy - implements runtime.Embedder
// Supports: text-embedding-3-small, text-embedding-3-large, text-embedding-ada-002
type OpenAIClient struct {
	apiKey    string
	model     string
	dimension int
}

// New creates a new OpenAI embeddings client.
// apiKey: OpenAI API key
// model: embedding model name (e.g., "text-embedding-3-small")
func New(apiKey string, opts ...runtime.EmbedderOption) *OpenAIClient {
	cfg := &runtime.EmbedderConfig{
		Model: "text-embedding-3-small", // default
	}
	for _, opt := range opts {
		opt(cfg)
	}

	// Determine dimension based on model
	dimension := 1536 // default for 3-small and ada-002
	if cfg.Model == "text-embedding-3-large" {
		dimension = 3072
	}

	return &OpenAIClient{
		apiKey:    apiKey,
		model:     cfg.Model,
		dimension: dimension,
	}
}

// Embed converts a single text string into a vector embedding.
func (c *OpenAIClient) Embed(ctx context.Context, text string) ([]float32, error) {
	embeddings, err := c.EmbedBatch(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	return embeddings[0], nil
}

// EmbedBatch converts multiple texts into vector embeddings.
// Uses OpenAI's embedding API with batching support.
// TODO: Implement actual OpenAI API call using openai-go SDK
// For now, returns placeholder embeddings
func (c *OpenAIClient) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("empty texts slice")
	}

	// Placeholder implementation - TODO: use actual SDK
	// response, err := c.client.Embeddings.Create(ctx, &openai.EmbeddingCreateParams{
	//     Input: openai.F(texts),
	//     Model: openai.F(c.model),
	// })
	// if err != nil {
	//     return nil, err
	// }

	// For now, create dummy embeddings of correct dimension
	embeddings := make([][]float32, len(texts))
	for i := range texts {
		embeddings[i] = make([]float32, c.dimension)
		// Fill with placeholder values based on text length
		sum := float32(0)
		for j, ch := range texts[i] {
			embeddings[i][j%c.dimension] += float32(ch) / 256.0
			sum += float32(ch)
		}
		// Normalize
		magnitude := float32(0)
		for j := range embeddings[i] {
			magnitude += embeddings[i][j] * embeddings[i][j]
		}
		if magnitude > 0 {
			magnitude = float32(1.0) / float32(magnitude)
			for j := range embeddings[i] {
				embeddings[i][j] *= magnitude
			}
		}
	}
	return embeddings, nil
}

// Dimension returns the embedding dimension for the configured model.
func (c *OpenAIClient) Dimension() int {
	return c.dimension
}
