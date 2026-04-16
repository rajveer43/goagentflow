package runtime

import "context"

// LLM is the core interface for language model providers.
// Pattern: Strategy - interchangeable LLM implementation
type LLM interface {
	Complete(ctx context.Context, prompt string, opts ...LLMOption) (string, error)
	Stream(ctx context.Context, prompt string, opts ...LLMOption) (<-chan string, <-chan error)
}

// TokenCounter is an optional interface for LLMs that support token counting.
// Implementations can use model-specific tokenization logic.
type TokenCounter interface {
	CountTokens(ctx context.Context, text string) (int, error)
}

// ModelInfoProvider is an optional interface for LLMs that provide model information.
// Useful for discovering model capabilities, costs, context size, etc.
type ModelInfoProvider interface {
	GetModelInfo() ModelInfo
}

// ModelInfo contains metadata about an LLM model.
type ModelInfo struct {
	Name             string    // e.g., "gpt-4o", "claude-opus-4-6"
	Provider         string    // e.g., "openai", "anthropic", "gemini"
	MaxTokens        int       // maximum tokens the model can process
	ContextSize      int       // context window size in tokens
	CostPer1KInput   float64   // cost per 1K input tokens (USD)
	CostPer1KOutput  float64   // cost per 1K output tokens (USD)
	Capabilities     []string  // "vision", "function_calling", "streaming", "json_mode", etc.
	ReleaseDate      string    // release date (YYYY-MM-DD)
}

type LLMOption func(*LLMConfig)

type LLMConfig struct {
	Temperature float64
}
