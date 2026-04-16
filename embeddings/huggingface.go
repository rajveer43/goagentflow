package embeddings

import (
	"context"
	"fmt"

	"github.com/rajveer43/goagentflow/runtime"
)

// HuggingFaceClient provides embeddings using HuggingFace Inference API or local models.
// Pattern: Strategy - implements runtime.Embedder
// Can point to HuggingFace hosted API or a local Hugging Face server.
type HuggingFaceClient struct {
	endpoint  string // e.g., "http://localhost:8080" for local, "https://api-inference.huggingface.co/models/..." for hosted
	apiKey    string // required for hosted API
	model     string
	dimension int
}

// NewHuggingFace creates a new HuggingFace embeddings client.
// endpoint: API endpoint URL
// apiKey: HuggingFace API token (can be empty for local endpoints)
// model: model name or identifier
func NewHuggingFace(endpoint, apiKey string, opts ...runtime.EmbedderOption) *HuggingFaceClient {
	cfg := &runtime.EmbedderConfig{
		Model: "sentence-transformers/all-MiniLM-L6-v2", // default
	}
	for _, opt := range opts {
		opt(cfg)
	}

	// Determine dimension based on model
	dimension := 384 // default for many sentence-transformers models
	switch cfg.Model {
	case "sentence-transformers/all-mpnet-base-v2":
		dimension = 768
	case "sentence-transformers/paraphrase-multilingual-mpnet-base-v2":
		dimension = 768
	case "sentence-transformers/all-MiniLM-L6-v2":
		dimension = 384
	}

	return &HuggingFaceClient{
		endpoint:  endpoint,
		apiKey:    apiKey,
		model:     cfg.Model,
		dimension: dimension,
	}
}

// Embed converts a single text string into a vector embedding.
func (c *HuggingFaceClient) Embed(ctx context.Context, text string) ([]float32, error) {
	embeddings, err := c.EmbedBatch(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	return embeddings[0], nil
}

// EmbedBatch converts multiple texts into vector embeddings.
// TODO: Implement actual HuggingFace API call using HTTP client
// For local servers: POST {endpoint}/predict with request: {"inputs": texts}
// For hosted API: POST {endpoint} with Authorization header and request: {"inputs": texts}
func (c *HuggingFaceClient) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("empty texts slice")
	}

	// Placeholder implementation
	// TODO: Implement HTTP call to HuggingFace endpoint
	// POST {c.endpoint} with JSON body containing texts

	embeddings := make([][]float32, len(texts))
	for i := range texts {
		embeddings[i] = make([]float32, c.dimension)
		for j := range embeddings[i] {
			embeddings[i][j] = 0.0
		}
	}
	return embeddings, nil
}

// Dimension returns the embedding dimension for the configured model.
func (c *HuggingFaceClient) Dimension() int {
	return c.dimension
}
