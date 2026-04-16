package loader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// HTMLLoader loads HTML from a local file and extracts text.
// Pattern: Strategy
type HTMLLoader struct {
	Path string
}

// NewHTMLLoader creates a new HTML file loader.
func NewHTMLLoader(path string) *HTMLLoader {
	return &HTMLLoader{Path: path}
}

// Load reads HTML file, strips tags, and returns clean text.
func (h *HTMLLoader) Load(ctx context.Context) ([]Document, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	file, err := os.Open(h.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open html file %s: %w", h.Path, err)
	}
	defer file.Close()

	text, title, err := extractHTMLText(file)
	if err != nil {
		return nil, err
	}

	doc := Document{
		PageContent: text,
		Metadata: map[string]any{
			"source":       h.Path,
			"title":        title,
			"content_type": "text/html",
		},
	}

	return []Document{doc}, nil
}

// HTMLURLLoader loads HTML from a remote URL and extracts text.
type HTMLURLLoader struct {
	URL    string
	Client *http.Client
}

// NewHTMLURLLoader creates a new HTML URL loader.
func NewHTMLURLLoader(url string) *HTMLURLLoader {
	return &HTMLURLLoader{
		URL:    url,
		Client: http.DefaultClient,
	}
}

// WithClient sets a custom HTTP client.
func (h *HTMLURLLoader) WithClient(client *http.Client) *HTMLURLLoader {
	h.Client = client
	return h
}

// Load fetches HTML from URL, strips tags, and returns clean text.
func (h *HTMLURLLoader) Load(ctx context.Context) ([]Document, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s: %w", h.URL, err)
	}
	req.Header.Set("User-Agent", "goagentflow/1.1.0")

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", h.URL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch %s returned status %d", h.URL, resp.StatusCode)
	}

	text, title, err := extractHTMLText(resp.Body)
	if err != nil {
		return nil, err
	}

	doc := Document{
		PageContent: text,
		Metadata: map[string]any{
			"source_url":   h.URL,
			"title":        title,
			"content_type": "text/html",
			"status_code":  resp.StatusCode,
		},
	}

	return []Document{doc}, nil
}

// extractHTMLText parses HTML and extracts clean text content + title.
// DSA: recursive tree traversal
func extractHTMLText(reader io.Reader) (string, string, error) {
	doc, err := html.Parse(reader)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse html: %w", err)
	}

	title := extractTitle(doc)
	text := extractText(doc)

	return strings.TrimSpace(text), title, nil
}

// extractTitle finds the <title> tag content.
func extractTitle(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "title" {
		if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
			return n.FirstChild.Data
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if title := extractTitle(c); title != "" {
			return title
		}
	}

	return ""
}

// extractText recursively extracts text content, skipping script/style tags.
// DSA: depth-first traversal
func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	// Skip script, style, and other non-content tags
	if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style" || n.Data == "noscript") {
		return ""
	}

	var text strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text.WriteString(extractText(c))
		// Add space between block elements
		if c.Type == html.ElementNode && isBlockElement(c.Data) {
			text.WriteString(" ")
		}
	}

	return text.String()
}

// isBlockElement checks if HTML tag is a block-level element.
func isBlockElement(tag string) bool {
	blockTags := map[string]bool{
		"p": true, "div": true, "h1": true, "h2": true, "h3": true, "h4": true,
		"h5": true, "h6": true, "blockquote": true, "pre": true, "ul": true,
		"ol": true, "li": true, "dl": true, "dt": true, "dd": true,
		"table": true, "tr": true, "td": true, "th": true, "section": true,
		"article": true, "nav": true, "aside": true, "main": true,
	}
	return blockTags[tag]
}
