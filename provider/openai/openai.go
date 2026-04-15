package openai

import (
	"context"
	"errors"
)

type Client struct {
	BaseURL string
	APIKey  string
}

func New(baseURL, apiKey string) *Client {
	return &Client{BaseURL: baseURL, APIKey: apiKey}
}

func (c *Client) Complete(ctx context.Context, prompt string, _ ...interface{}) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}
	if c == nil {
		return "", errors.New("nil client")
	}
	return prompt, nil
}
