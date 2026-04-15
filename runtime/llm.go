package goagentflow

import "context"

type LLM interface {
	Complete(ctx context.Context, prompt string, opts ...LLMOption) (string, error)
	Stream(ctx context.Context, prompt string, opts ...LLMOption) (<-chan string, <-chan error)
}

type LLMOption func(*LLMConfig)

type LLMConfig struct {
	Temperature float64
}
