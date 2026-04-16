package types

// Document represents a loaded chunk of content from any source.
// Shared across loader, runtime, and vector store packages to avoid circular imports.
// Pattern: mirrors LangChain's document structure
type Document struct {
	PageContent string         // the actual text content
	Metadata    map[string]any // source, page number, content-type, etc.
}
