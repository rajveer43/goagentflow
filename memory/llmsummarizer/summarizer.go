package llmsummarizer

import (
	"context"
	"fmt"
	"strings"

	"github.com/rajveer43/goagentflow/runtime"
)

// Summarizer wraps an LLM to create a Summarizer interface.
// Pattern: Adapter - adapts runtime.LLM to runtime.Summarizer interface
type Summarizer struct {
	llm            runtime.LLM
	systemPrompt   string
}

// New creates a new LLM-based summarizer.
// llm: the language model to use for summarization
func New(llm runtime.LLM) *Summarizer {
	return &Summarizer{
		llm: llm,
		systemPrompt: `You are a concise conversation summarizer.
Summarize the following conversation messages into a brief, factual summary that captures the key points, decisions, and important information.
Keep the summary under 200 words.
Focus on what was discussed, not how many messages there were.`,
	}
}

// NewWithPrompt creates a summarizer with a custom system prompt.
func NewWithPrompt(llm runtime.LLM, systemPrompt string) *Summarizer {
	if systemPrompt == "" {
		systemPrompt = `You are a conversation summarizer. Summarize the key points of this conversation.`
	}
	return &Summarizer{
		llm:          llm,
		systemPrompt: systemPrompt,
	}
}

// Summarize generates a summary of the given messages using the LLM.
func (s *Summarizer) Summarize(ctx context.Context, messages []runtime.Message) (string, error) {
	if len(messages) == 0 {
		return "", nil
	}

	// Build prompt from messages
	prompt := s.buildSummarizationPrompt(messages)

	// Call LLM
	summary, err := s.llm.Complete(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("summarization failed: %w", err)
	}

	return strings.TrimSpace(summary), nil
}

// buildSummarizationPrompt constructs the prompt for summarization.
func (s *Summarizer) buildSummarizationPrompt(messages []runtime.Message) string {
	var sb strings.Builder

	sb.WriteString(s.systemPrompt)
	sb.WriteString("\n\n")
	sb.WriteString("Conversation to summarize:\n")
	sb.WriteString("---\n")

	for _, msg := range messages {
		sb.WriteString(fmt.Sprintf("%s: %s\n", msg.Role, msg.Content))
	}

	sb.WriteString("---\n\n")
	sb.WriteString("Summary:")

	return sb.String()
}

// SetSystemPrompt updates the system prompt used for summarization.
func (s *Summarizer) SetSystemPrompt(prompt string) {
	if prompt != "" {
		s.systemPrompt = prompt
	}
}
