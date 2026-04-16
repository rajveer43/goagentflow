package ollama

import (
	"context"
	"fmt"

	"github.com/rajveer43/goagentflow/runtime"
)

// Client implements runtime.LLM using local Ollama server.
// Pattern: Strategy - interchangeable LLM provider
// Connects to a local Ollama instance via HTTP (default: http://localhost:11434)
// Supports: llama2, llama3.2, mistral, neural-chat, codellama, etc.
// Note: This is a stub implementation. For production use, implement HTTP calls to Ollama API.
type Client struct {
	endpoint string // Ollama server endpoint (e.g., "http://localhost:11434")
	model    string
}

// New creates a new Ollama LLM client.
// endpoint: Ollama server endpoint (e.g., "http://localhost:11434")
// model: model name (e.g., "llama2", "llama3.2", "mistral")
func New(endpoint, model string) *Client {
	if endpoint == "" {
		endpoint = "http://localhost:11434" // default local endpoint
	}
	return &Client{
		endpoint: endpoint,
		model:    model,
	}
}

// Complete sends a single prompt and returns the full response.
// Applies temperature from LLMConfig if provided.
// Pattern: Strategy interface implementation
// TODO: Implement actual Ollama HTTP call
func (c *Client) Complete(ctx context.Context, prompt string, opts ...runtime.LLMOption) (string, error) {
	// Apply options to config
	cfg := &runtime.LLMConfig{Temperature: 0.7} // default
	for _, opt := range opts {
		opt(cfg)
	}

	// TODO: Replace with actual Ollama API call
	// POST {c.endpoint}/api/generate
	// Body: {
	//   "model": c.model,
	//   "prompt": prompt,
	//   "temperature": cfg.Temperature,
	//   "stream": false
	// }

	// For now, return a placeholder response
	return fmt.Sprintf("Ollama (%s) response to: %s (temp: %.1f)", c.model, prompt, cfg.Temperature), nil
}

// Stream sends a prompt and returns token and error channels.
// Pattern: Strategy interface implementation
// TODO: Implement actual Ollama streaming API call
func (c *Client) Stream(ctx context.Context, prompt string, opts ...runtime.LLMOption) (<-chan string, <-chan error) {
	tokensCh := make(chan string, 10)   // buffered for backpressure
	errorsCh := make(chan error, 1)

	// Apply options
	cfg := &runtime.LLMConfig{Temperature: 0.7}
	for _, opt := range opts {
		opt(cfg)
	}

	// Run streaming in a goroutine
	go func() {
		defer close(tokensCh)
		defer close(errorsCh)

		// TODO: Replace with actual Ollama streaming API
		// POST {c.endpoint}/api/generate with "stream": true
		// Read streamed JSON objects with {"response": token} fields

		// For now, send a simple placeholder response
		tokens := []string{"Ollama", " response", " token", " by", " token"}
		for _, token := range tokens {
			select {
			case tokensCh <- token:
			case <-ctx.Done():
				errorsCh <- ctx.Err()
				return
			}
		}
	}()

	return tokensCh, errorsCh
}
