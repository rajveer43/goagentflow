package runtime

import "context"

// Embedder converts text into vector embeddings.
// Pattern: Strategy - interchangeable embedding provider
type Embedder interface {
	// Embed converts a single text string into a vector embedding.
	Embed(ctx context.Context, text string) ([]float32, error)
	// EmbedBatch converts multiple texts into vector embeddings.
	// Returns a slice of embeddings in the same order as the input texts.
	EmbedBatch(ctx context.Context, texts []string) ([][]float32, error)
	// Dimension returns the size of embeddings produced by this embedder.
	Dimension() int
}

// EmbedderOption is a functional option for embedder configuration.
type EmbedderOption func(*EmbedderConfig)

// EmbedderConfig holds embedder-specific configuration.
type EmbedderConfig struct {
	Model string
	// additional provider-specific options can be added here
}

// WithModel sets the embedding model name.
func WithModel(model string) EmbedderOption {
	return func(cfg *EmbedderConfig) {
		cfg.Model = model
	}
}
