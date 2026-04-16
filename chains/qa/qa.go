package qa

import (
	"context"
	"fmt"
	"strings"

	"github.com/rajveer43/goagentflow/loader"
	"github.com/rajveer43/goagentflow/runtime"
)

// Chain implements a question-answering chain using retrieval-augmented generation.
// Pattern: Composition - combines Retriever + LLM into a runtime.Chain
// Input: question string; Output: answer string with optional source documents
type Chain struct {
	retriever       runtime.Retriever
	llm             runtime.LLM
	k               int    // number of documents to retrieve
	promptTemplate  string // template for combining context and query
	contextSeparator string
}

// Input represents the input to a QA chain.
type Input struct {
	Question string
}

// Output represents the output from a QA chain.
type Output struct {
	Answer  string
	Sources []loader.Document
}

// New creates a new QA chain.
// retriever: retrieves relevant documents for the question
// llm: generates answers based on retrieved context
// k: number of top documents to retrieve
func New(retriever runtime.Retriever, llm runtime.LLM, k int) *Chain {
	if k <= 0 {
		k = 3 // default
	}
	return &Chain{
		retriever:        retriever,
		llm:              llm,
		k:                k,
		promptTemplate:   defaultPrompt,
		contextSeparator: "\n\n---\n\n",
	}
}

// SetPromptTemplate sets a custom prompt template.
// Template should include {context} and {question} placeholders.
func (c *Chain) SetPromptTemplate(template string) {
	c.promptTemplate = template
}

// Run implements runtime.Chain interface.
// Input: question string or Input struct
// Output: answer string or Output struct (with sources)
func (c *Chain) Run(ctx context.Context, input any) (any, error) {
	// Handle flexible input types
	question := ""
	switch v := input.(type) {
	case string:
		question = v
	case Input:
		question = v.Question
	default:
		return nil, fmt.Errorf("expected string or Input, got %T", input)
	}

	if question == "" {
		return nil, fmt.Errorf("question cannot be empty")
	}

	// Retrieve relevant documents
	docs, err := c.retriever.Retrieve(ctx, question, c.k)
	if err != nil {
		return nil, fmt.Errorf("retrieval failed: %w", err)
	}

	// Format documents as context
	context := c.formatDocuments(docs)

	// Build final prompt
	prompt := strings.NewReplacer(
		"{context}", context,
		"{question}", question,
	).Replace(c.promptTemplate)

	// Get LLM response
	answer, err := c.llm.Complete(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("completion failed: %w", err)
	}

	return Output{
		Answer:  answer,
		Sources: docs,
	}, nil
}

// formatDocuments formats retrieved documents into a context string.
func (c *Chain) formatDocuments(docs []loader.Document) string {
	if len(docs) == 0 {
		return ""
	}

	var parts []string
	for i, doc := range docs {
		part := fmt.Sprintf("Document %d:\n%s", i+1, doc.PageContent)
		parts = append(parts, part)
	}

	return strings.Join(parts, c.contextSeparator)
}

const defaultPrompt = `You are a helpful assistant. Answer the question based on the provided documents.

Context:
{context}

Question: {question}

Answer:`
