package loader

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// URLLoader fetches content from a web URL.
// Pattern: Strategy
type URLLoader struct {
	URL    string
	Client *http.Client
}

// NewURLLoader creates a new URL loader with the default HTTP client.
func NewURLLoader(url string) *URLLoader {
	return &URLLoader{
		URL:    url,
		Client: http.DefaultClient,
	}
}

// WithClient sets a custom HTTP client.
func (u *URLLoader) WithClient(client *http.Client) *URLLoader {
	u.Client = client
	return u
}

// Load fetches the URL and returns the response body as a Document.
func (u *URLLoader) Load(ctx context.Context) ([]Document, error) {
	// Create HTTP request with context for cancellation
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s: %w", u.URL, err)
	}

	// Set a reasonable user-agent
	req.Header.Set("User-Agent", "goagentflow/1.1.0")

	// Execute request
	resp, err := u.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", u.URL, err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch %s returned status %d", u.URL, resp.StatusCode)
	}

	// Read response body
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response from %s: %w", u.URL, err)
	}

	// Create document
	doc := Document{
		PageContent: string(content),
		Metadata: map[string]any{
			"source_url":   u.URL,
			"status_code":  resp.StatusCode,
			"content_type": resp.Header.Get("Content-Type"),
		},
	}

	return []Document{doc}, nil
}
