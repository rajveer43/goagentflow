package summarization

import (
	"context"
	"fmt"
	"strings"

	"github.com/rajveer43/goagentflow/loader"
	"github.com/rajveer43/goagentflow/runtime"
)

// Strategy defines summarization approach.
type Strategy string

const (
	StuffStrategy      Strategy = "stuff"       // concat all + summarize once
	MapReduceStrategy  Strategy = "map_reduce"  // summarize each chunk, then combine
)

// Chain implements document summarization.
// Pattern: Composition - combines document processing + LLM into runtime.Chain
// Supports multiple strategies for handling large documents.
type Chain struct {
	llm           runtime.LLM
	strategy      Strategy
	chunkSize     int
	chunkOverlap  int
	promptTemplate string
}

// New creates a new summarization chain.
// llm: language model for summarization
// strategy: "stuff" (simple concat) or "map_reduce" (hierarchical)
func New(llm runtime.LLM, strategy Strategy) *Chain {
	if strategy != MapReduceStrategy {
		strategy = StuffStrategy // default
	}
	return &Chain{
		llm:            llm,
		strategy:       strategy,
		chunkSize:      1000,  // chars per chunk
		chunkOverlap:   100,   // overlap between chunks
		promptTemplate: defaultSummarizationPrompt,
	}
}

// SetPromptTemplate sets the summarization prompt template.
func (c *Chain) SetPromptTemplate(template string) {
	c.promptTemplate = template
}

// SetChunkSize configures chunk size for map_reduce strategy.
func (c *Chain) SetChunkSize(size, overlap int) {
	c.chunkSize = size
	c.chunkOverlap = overlap
}

// Run implements runtime.Chain interface.
// Input: string, []string, []loader.Document
// Output: string (summary)
func (c *Chain) Run(ctx context.Context, input any) (any, error) {
	texts := []string{}

	// Handle flexible input types
	switch v := input.(type) {
	case string:
		texts = []string{v}
	case []string:
		texts = v
	case []loader.Document:
		for _, doc := range v {
			texts = append(texts, doc.PageContent)
		}
	default:
		return nil, fmt.Errorf("expected string, []string, or []Document, got %T", input)
	}

	if len(texts) == 0 {
		return nil, fmt.Errorf("no documents to summarize")
	}

	// Summarize based on strategy
	switch c.strategy {
	case MapReduceStrategy:
		return c.mapReduceSummarize(ctx, texts)
	default:
		return c.stuffSummarize(ctx, texts)
	}
}

// stuffSummarize concatenates all text and summarizes once.
// Simple but may exceed context limits for large documents.
func (c *Chain) stuffSummarize(ctx context.Context, texts []string) (string, error) {
	combined := strings.Join(texts, "\n\n---\n\n")

	prompt := fmt.Sprintf(c.promptTemplate, combined)
	summary, err := c.llm.Complete(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("summarization failed: %w", err)
	}

	return strings.TrimSpace(summary), nil
}

// mapReduceSummarize chunks documents, summarizes each, then combines summaries.
// More robust for large documents but requires multiple LLM calls.
func (c *Chain) mapReduceSummarize(ctx context.Context, texts []string) (string, error) {
	// Chunk each document
	var chunks []string
	for _, text := range texts {
		chunked := c.chunkText(text)
		chunks = append(chunks, chunked...)
	}

	// Summarize each chunk ("map" phase)
	summaries := make([]string, 0, len(chunks))
	for _, chunk := range chunks {
		prompt := fmt.Sprintf(c.promptTemplate, chunk)
		summary, err := c.llm.Complete(ctx, prompt)
		if err != nil {
			return "", fmt.Errorf("chunk summarization failed: %w", err)
		}
		summaries = append(summaries, strings.TrimSpace(summary))
	}

	// Combine summaries and do final summarization ("reduce" phase)
	if len(summaries) == 1 {
		return summaries[0], nil
	}

	combined := strings.Join(summaries, "\n\n")
	finalPrompt := fmt.Sprintf(`Please create a final summary combining these summaries:

%s

Final Summary:`, combined)

	finalSummary, err := c.llm.Complete(ctx, finalPrompt)
	if err != nil {
		return "", fmt.Errorf("final summarization failed: %w", err)
	}

	return strings.TrimSpace(finalSummary), nil
}

// chunkText splits text into overlapping chunks.
func (c *Chain) chunkText(text string) []string {
	if len(text) <= c.chunkSize {
		return []string{text}
	}

	var chunks []string
	var current string

	for _, ch := range text {
		current += string(ch)

		if len(current) >= c.chunkSize {
			chunks = append(chunks, current)

			// Apply overlap
			if c.chunkOverlap > 0 && len(current) > c.chunkOverlap {
				current = current[len(current)-c.chunkOverlap:]
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

const defaultSummarizationPrompt = `Please summarize the following text concisely:

%s

Summary:`
