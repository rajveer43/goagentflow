package loader

import (
	"context"

	"github.com/rajveer43/goagentflow/types"
)

// Document is an alias to types.Document for backward compatibility.
// This allows existing code to use loader.Document while avoiding circular imports.
type Document = types.Document

// Loader is the single interface all loaders implement.
// Pattern: Strategy — swap sources freely
type Loader interface {
	Load(ctx context.Context) ([]Document, error)
}

// LoaderFunc is a convenience adapter for single-function loaders.
// Pattern: mirrors runtime.ChainFunc
type LoaderFunc func(ctx context.Context) ([]Document, error)

// Load implements the Loader interface.
func (f LoaderFunc) Load(ctx context.Context) ([]Document, error) {
	return f(ctx)
}

// MultiLoader loads from multiple sources and concatenates results.
// Pattern: Composite — combine multiple loaders
type MultiLoader struct {
	loaders []Loader
}

// NewMultiLoader creates a loader that runs multiple loaders in sequence.
func NewMultiLoader(loaders ...Loader) *MultiLoader {
	return &MultiLoader{loaders: loaders}
}

// Load runs each loader in sequence and concatenates results.
func (m *MultiLoader) Load(ctx context.Context) ([]Document, error) {
	var docs []Document
	for _, loader := range m.loaders {
		loaded, err := loader.Load(ctx)
		if err != nil {
			return nil, err
		}
		docs = append(docs, loaded...)
	}
	return docs, nil
}
