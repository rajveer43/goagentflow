package agent

import (
	"context"
	"fmt"
	"strings"

	"github.com/rajveer43/goagentflow/runtime"
)

// Chain is an enhanced agent orchestrator combining LLM, tools, memory, and retrieval.
// Pattern: Composition - orchestrates multiple components into a runtime.Chain
// Input: user message string; Output: response string
type Chain struct {
	llm       runtime.LLM
	memory    runtime.Memory
	retriever runtime.Retriever
	tools     map[string]runtime.Tool
	k         int // number of docs to retrieve for context
	maxSteps  int
}

// New creates a new agent orchestrator chain.
// llm: language model for reasoning
// memory: conversation memory backend
// retriever: optional retriever for context (can be nil)
func New(llm runtime.LLM, memory runtime.Memory) *Chain {
	return &Chain{
		llm:      llm,
		memory:   memory,
		tools:    make(map[string]runtime.Tool),
		k:        3,
		maxSteps: 10,
	}
}

// SetRetriever sets the retriever for context retrieval.
func (c *Chain) SetRetriever(retriever runtime.Retriever) {
	c.retriever = retriever
}

// RegisterTool registers a tool that the agent can use.
func (c *Chain) RegisterTool(tool runtime.Tool) error {
	if tool.Name() == "" {
		return fmt.Errorf("tool name cannot be empty")
	}
	c.tools[tool.Name()] = tool
	return nil
}

// SetMaxSteps sets the maximum number of reasoning steps.
func (c *Chain) SetMaxSteps(steps int) {
	if steps > 0 {
		c.maxSteps = steps
	}
}

// Run implements runtime.Chain interface.
// Input: user message string
// Output: agent response string
func (c *Chain) Run(ctx context.Context, input any) (any, error) {
	message, ok := input.(string)
	if !ok {
		return nil, fmt.Errorf("expected string message, got %T", input)
	}

	if message == "" {
		return nil, fmt.Errorf("message cannot be empty")
	}

	// Store user message in memory
	if c.memory != nil {
		if err := c.memory.AddMessage(ctx, runtime.Message{
			Role:    "user",
			Content: message,
		}); err != nil {
			return nil, fmt.Errorf("failed to store message: %w", err)
		}
	}

	// Build context from memory and retriever
	context := c.buildContext(ctx, message)

	// Build prompt for agent reasoning
	prompt := c.buildPrompt(context, message)

	// Get LLM response (in production, this would be agentic loop with tool calling)
	response, err := c.llm.Complete(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("agent reasoning failed: %w", err)
	}

	response = strings.TrimSpace(response)

	// Store response in memory
	if c.memory != nil {
		if err := c.memory.AddMessage(ctx, runtime.Message{
			Role:    "assistant",
			Content: response,
		}); err != nil {
			return nil, fmt.Errorf("failed to store response: %w", err)
		}
	}

	return response, nil
}

// buildContext constructs the context for the agent from memory and retrieval.
func (c *Chain) buildContext(ctx context.Context, message string) string {
	var sb strings.Builder

	// Add conversation history if memory available
	if c.memory != nil {
		if messages, err := c.memory.GetMessages(ctx); err == nil && len(messages) > 0 {
			sb.WriteString("Conversation History:\n")
			// Show last few messages for context
			start := len(messages) - 5
			if start < 0 {
				start = 0
			}
			for _, msg := range messages[start:] {
				sb.WriteString(fmt.Sprintf("%s: %s\n", msg.Role, msg.Content))
			}
			sb.WriteString("\n")
		}
	}

	// Add retrieved context if retriever available
	if c.retriever != nil {
		if docs, err := c.retriever.Retrieve(ctx, message, c.k); err == nil && len(docs) > 0 {
			sb.WriteString("Retrieved Context:\n")
			for i, doc := range docs {
				sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, doc.PageContent[:min(len(doc.PageContent), 200)]))
			}
			sb.WriteString("\n")
		}
	}

	// Add available tools if any
	if len(c.tools) > 0 {
		sb.WriteString("Available Tools:\n")
		for name, tool := range c.tools {
			sb.WriteString(fmt.Sprintf("- %s: %s\n", name, tool.Description()))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// buildPrompt constructs the prompt for agent reasoning.
func (c *Chain) buildPrompt(context, message string) string {
	var sb strings.Builder

	sb.WriteString(`You are a helpful assistant with access to information and tools.
Use the provided context and tools to help the user.

`)

	if context != "" {
		sb.WriteString(context)
	}

	sb.WriteString(fmt.Sprintf("User Message: %s\n\nResponse:", message))

	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
