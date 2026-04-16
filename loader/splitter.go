package loader

import (
	"context"
	"fmt"
	"strings"

	"github.com/rajveer43/goagentflow/types"
)

// TextSplitter chunks documents into smaller pieces before embedding.
// Supports both character-based and token-based chunking.
type TextSplitter interface {
	Split(ctx context.Context, text string) ([]string, error)
	SplitDocuments(ctx context.Context, docs []types.Document) ([]types.Document, error)
}

// CharacterSplitter chunks text by character count with optional overlap.
// Pattern: Strategy - different chunking strategies can be swapped
type CharacterSplitter struct {
	ChunkSize    int // size of each chunk in characters
	ChunkOverlap int // overlap between chunks in characters
	Separator    string // separator to split on (e.g., "\n\n")
}

// NewCharacterSplitter creates a new character-based text splitter.
// chunkSize: target chunk size in characters
// chunkOverlap: number of overlapping characters between chunks (must be < chunkSize)
func NewCharacterSplitter(chunkSize, chunkOverlap int) *CharacterSplitter {
	if chunkOverlap >= chunkSize {
		chunkOverlap = chunkSize - 1
	}
	return &CharacterSplitter{
		ChunkSize:    chunkSize,
		ChunkOverlap: chunkOverlap,
		Separator:    "\n\n",
	}
}

// Split chunks text into pieces.
func (cs *CharacterSplitter) Split(ctx context.Context, text string) ([]string, error) {
	if cs.ChunkSize <= 0 {
		return nil, fmt.Errorf("chunk size must be positive, got %d", cs.ChunkSize)
	}

	// Try to split by separator first
	parts := strings.Split(text, cs.Separator)
	if len(parts) == 1 {
		// No separator found, fall back to character split
		return cs.splitByCharacter(text), nil
	}

	// Merge parts to create chunks of appropriate size
	var chunks []string
	var current string

	for _, part := range parts {
		// If adding this part would exceed chunk size, save current and start new
		if len(current) > 0 && len(current)+len(cs.Separator)+len(part) > cs.ChunkSize {
			chunks = append(chunks, current)
			// Add overlap for next chunk
			if cs.ChunkOverlap > 0 && len(current) > cs.ChunkOverlap {
				current = current[len(current)-cs.ChunkOverlap:]
			} else {
				current = ""
			}
		}

		// Add separator if current is not empty
		if len(current) > 0 {
			current += cs.Separator
		}
		current += part
	}

	// Add remaining text
	if len(current) > 0 {
		chunks = append(chunks, current)
	}

	return chunks, nil
}

// SplitDocuments splits documents into smaller documents with overlapping chunks.
// Preserves metadata from original documents.
func (cs *CharacterSplitter) SplitDocuments(ctx context.Context, docs []types.Document) ([]types.Document, error) {
	var result []types.Document

	for _, doc := range docs {
		chunks, err := cs.Split(ctx, doc.PageContent)
		if err != nil {
			return nil, err
		}

		for _, chunk := range chunks {
			newDoc := types.Document{
				PageContent: chunk,
				Metadata:    doc.Metadata, // copy metadata
			}
			result = append(result, newDoc)
		}
	}

	return result, nil
}

// splitByCharacter splits text by character count when no natural separator exists.
func (cs *CharacterSplitter) splitByCharacter(text string) []string {
	var chunks []string
	var current string

	for _, ch := range text {
		current += string(ch)

		if len(current) >= cs.ChunkSize {
			chunks = append(chunks, current)
			// Apply overlap
			if cs.ChunkOverlap > 0 && len(current) > cs.ChunkOverlap {
				current = current[len(current)-cs.ChunkOverlap:]
			} else {
				current = ""
			}
		}
	}

	if len(current) > 0 {
		chunks = append(chunks, current)
	}

	return chunks
}

// RecursiveCharacterSplitter chunks text recursively by trying multiple separators.
// More sophisticated than CharacterSplitter.
type RecursiveCharacterSplitter struct {
	ChunkSize    int
	ChunkOverlap int
	Separators   []string // ordered list of separators to try (e.g., ["\n\n", "\n", " ", ""])
}

// NewRecursiveCharacterSplitter creates a new recursive character splitter.
func NewRecursiveCharacterSplitter(chunkSize, chunkOverlap int) *RecursiveCharacterSplitter {
	if chunkOverlap >= chunkSize {
		chunkOverlap = chunkSize - 1
	}
	return &RecursiveCharacterSplitter{
		ChunkSize:    chunkSize,
		ChunkOverlap: chunkOverlap,
		Separators:   []string{"\n\n", "\n", " ", ""}, // default separators
	}
}

// Split chunks text recursively.
func (rs *RecursiveCharacterSplitter) Split(ctx context.Context, text string) ([]string, error) {
	return rs.splitRecursive(text, rs.Separators), nil
}

// SplitDocuments splits documents into smaller documents.
func (rs *RecursiveCharacterSplitter) SplitDocuments(ctx context.Context, docs []types.Document) ([]types.Document, error) {
	var result []types.Document

	for _, doc := range docs {
		chunks, err := rs.Split(ctx, doc.PageContent)
		if err != nil {
			return nil, err
		}

		for _, chunk := range chunks {
			newDoc := types.Document{
				PageContent: chunk,
				Metadata:    doc.Metadata,
			}
			result = append(result, newDoc)
		}
	}

	return result, nil
}

// splitRecursive recursively splits text by trying separators in order.
func (rs *RecursiveCharacterSplitter) splitRecursive(text string, separators []string) []string {
	separator := separators[len(separators)-1] // empty string fallback

	for _, s := range separators {
		if s == "" {
			break
		}
		if !strings.Contains(text, s) {
			continue
		}

		separator = s
		break
	}

	// Split by the best separator
	var splits []string
	if separator != "" {
		splits = strings.Split(text, separator)
	} else {
		splits = []string{text}
	}

	// Merge splits to create appropriately-sized chunks
	return rs.mergeSplits(splits, separator)
}

// mergeSplits merges splits recursively to create appropriately-sized chunks.
func (rs *RecursiveCharacterSplitter) mergeSplits(splits []string, separator string) []string {
	var merged []string
	separator = strings.TrimSpace(separator)

	var goodSplits []string
	for _, split := range splits {
		split = strings.TrimSpace(split)
		if len(split) == 0 {
			continue
		}
		goodSplits = append(goodSplits, split)
	}

	// Group splits into chunks of appropriate size
	var current string
	for _, split := range goodSplits {
		candidate := current
		if len(candidate) > 0 {
			candidate += separator + split
		} else {
			candidate = split
		}

		if len(candidate) <= rs.ChunkSize {
			current = candidate
		} else {
			// Current chunk would exceed size, save it and start new
			if len(current) > 0 {
				merged = append(merged, current)
			}
			// Apply overlap if needed
			if rs.ChunkOverlap > 0 && len(current) > rs.ChunkOverlap {
				current = current[len(current)-rs.ChunkOverlap:] + separator + split
			} else {
				current = split
			}
		}
	}

	// Add remaining text
	if len(current) > 0 {
		merged = append(merged, current)
	}

	return merged
}
