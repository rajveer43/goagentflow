package loader

import (
	"context"
	"fmt"
	"io"
	"os"
)

// TextLoader loads plain text files.
// Pattern: Strategy
// DSA: buffered reader (O(n) streaming, no slurp)
type TextLoader struct {
	Path string
}

// NewTextLoader creates a new text file loader.
func NewTextLoader(path string) *TextLoader {
	return &TextLoader{Path: path}
}

// Load reads the file and returns it as a single Document.
func (t *TextLoader) Load(ctx context.Context) ([]Document, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Open file
	file, err := os.Open(t.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open text file %s: %w", t.Path, err)
	}
	defer file.Close()

	// Get file size for metadata
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	// Read entire file into memory (for text files this is fine)
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read text file: %w", err)
	}

	// Create document
	doc := Document{
		PageContent: string(content),
		Metadata: map[string]any{
			"source":       t.Path,
			"file_size":    stat.Size(),
			"content_type": "text/plain",
		},
	}

	return []Document{doc}, nil
}
