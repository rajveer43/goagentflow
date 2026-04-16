package main

import (
	"context"
	"fmt"

	"github.com/rajveer43/goagentflow/chains/agent"
	qachain "github.com/rajveer43/goagentflow/chains/qa"
	sqlchain "github.com/rajveer43/goagentflow/chains/sql"
	summarychain "github.com/rajveer43/goagentflow/chains/summarization"
	"github.com/rajveer43/goagentflow/memory/inmemory"
	"github.com/rajveer43/goagentflow/provider/anthropic"
	"github.com/rajveer43/goagentflow/runtime"
	"github.com/rajveer43/goagentflow/types"
)

// Example demonstrates all pre-built chains in goagentflow.
// Shows:
// 1. QA Chain - Question answering with retrieval
// 2. Summarization Chain - Document and text summarization
// 3. SQL Chain - Natural language to SQL generation
// 4. Agent Chain - Enhanced agent orchestrator
func main() {
	ctx := context.Background()

	fmt.Println("=== Pre-Built Chains Demo ===\n")

	// Setup: Create an LLM provider for all chains
	llm := anthropic.New("your-api-key", "claude-opus-4-6")

	// Example 1: QA Chain (requires a retriever)
	fmt.Println("1. QA Chain (Question Answering with Retrieval)")
	fmt.Println("---")
	demonstrateQAChain(ctx, llm)
	fmt.Println()

	// Example 2: Summarization Chain
	fmt.Println("2. Summarization Chain (Document & Text Summarization)")
	fmt.Println("---")
	demonstrateSummarizationChain(ctx, llm)
	fmt.Println()

	// Example 3: SQL Chain
	fmt.Println("3. SQL Chain (Natural Language to SQL)")
	fmt.Println("---")
	demonstrateSQLChain(ctx, llm)
	fmt.Println()

	// Example 4: Agent Chain
	fmt.Println("4. Agent Chain (Orchestrator)")
	fmt.Println("---")
	demonstrateAgentChain(ctx, llm)
	fmt.Println()

	// Example 5: Chaining chains together
	fmt.Println("5. Composing Chains (Chaining chains together)")
	fmt.Println("---")
	demonstrateChainComposition(ctx, llm)
	fmt.Println()

	fmt.Println("=== Pre-Built Chains Summary ===")
	fmt.Println("✓ QA Chain - Retrieve context + answer questions")
	fmt.Println("✓ Summarization Chain - Stuff or map-reduce summarization strategies")
	fmt.Println("✓ SQL Chain - Generate SQL from natural language")
	fmt.Println("✓ Agent Chain - Orchestrate LLM + tools + memory + retrieval")
	fmt.Println("\nAll chains implement runtime.Chain and can be composed with ChainPipeline!")
}

// demonstrateQAChain shows how to use the QA chain.
func demonstrateQAChain(ctx context.Context, llm runtime.LLM) {
	// For this demo, we'll create a mock retriever
	mockRetriever := &mockRetriever{}

	// Create QA chain
	qaChain := qachain.New(mockRetriever, llm, 2)

	// Custom prompt (optional)
	qaChain.SetPromptTemplate(`Answer based on context:
{context}

Q: {question}
A:`)

	// Run the chain
	input := "What is machine learning?"
	result, err := qaChain.Run(ctx, input)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if output, ok := result.(qachain.Output); ok {
		fmt.Printf("Question: %s\n", input)
		fmt.Printf("Answer: %s\n", output.Answer)
		fmt.Printf("Sources: %d documents\n", len(output.Sources))
	} else {
		fmt.Printf("Response: %v\n", result)
	}
}

// demonstrateSummarizationChain shows how to use the summarization chain.
func demonstrateSummarizationChain(ctx context.Context, llm runtime.LLM) {
	// Create summarization chain with "stuff" strategy (simple concat)
	stuffChain := summarychain.New(llm, summarychain.StuffStrategy)

	// Create with "map_reduce" strategy (hierarchical)
	mapReduceChain := summarychain.New(llm, summarychain.MapReduceStrategy)
	mapReduceChain.SetChunkSize(500, 50) // 500 char chunks with 50 char overlap

	// Example 1: Simple text summarization
	text := `Artificial Intelligence (AI) is transforming industries worldwide.
Machine learning enables systems to learn from data without explicit programming.
Deep learning uses neural networks to process complex patterns.
Natural language processing allows computers to understand human language.`

	result, err := stuffChain.Run(ctx, text)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Input text: %d chars\n", len(text))
	fmt.Printf("Summary (%v strategy): %v\n", summarychain.StuffStrategy, result)

	// Example 2: Summarize documents
	docs := []types.Document{
		{PageContent: text, Metadata: map[string]any{"source": "article1"}},
	}

	result2, err := mapReduceChain.Run(ctx, docs)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("\nMap-reduce summary: %v\n", result2)
}

// demonstrateSQLChain shows how to use the SQL generation chain.
func demonstrateSQLChain(ctx context.Context, llm runtime.LLM) {
	// Define a database schema
	schema := `
Users table: id (int), name (string), email (string), created_at (date)
Posts table: id (int), user_id (int), title (string), content (text), created_at (date)
Comments table: id (int), post_id (int), user_id (int), content (text), created_at (date)
`

	// Create SQL chain
	sqlChain := sqlchain.New(llm, schema)
	sqlChain.SetDialect("PostgreSQL")

	// Generate SQL
	question := "Find all posts created in the last 30 days"
	result, err := sqlChain.Run(ctx, question)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if output, ok := result.(sqlchain.Output); ok {
		fmt.Printf("Question: %s\n", question)
		fmt.Printf("Generated SQL:\n%s\n", output.SQL)
		if output.Explanation != "" {
			fmt.Printf("Explanation: %s\n", output.Explanation)
		}
	} else {
		fmt.Printf("SQL: %v\n", result)
	}
}

// demonstrateAgentChain shows how to use the agent orchestrator.
func demonstrateAgentChain(ctx context.Context, llm runtime.LLM) {
	// Create memory for conversation
	memory := inmemory.New()

	// Create agent chain
	agentChain := agent.New(llm, memory)
	agentChain.SetMaxSteps(5)

	// Run agent
	message := "I need help understanding cloud computing"
	result, err := agentChain.Run(ctx, message)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("User: %s\n", message)
	fmt.Printf("Agent: %v\n", result)

	// Continue conversation (memory retains context)
	followUp := "What are the main benefits?"
	result2, err := agentChain.Run(ctx, followUp)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("\nUser: %s\n", followUp)
	fmt.Printf("Agent: %v\n", result2)
}

// demonstrateChainComposition shows how to compose multiple chains.
func demonstrateChainComposition(ctx context.Context, llm runtime.LLM) {
	// Create a pipeline that summarizes then answers questions
	summarizer := summarychain.New(llm, summarychain.StuffStrategy)

	// In a real scenario, you'd chain these using runtime.ChainPipeline:
	// pipeline := runtime.NewChainPipeline(summarizer, qaChain)
	// result := pipeline.Run(ctx, input)

	// For demo, show the concept
	text := `The Internet of Things (IoT) connects billions of devices.
Smart homes use IoT for automation and efficiency.
Industrial IoT improves manufacturing processes.
Healthcare IoT enables remote patient monitoring.`

	summary, err := summarizer.Run(ctx, text)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Original text: %d chars\n", len(text))
	fmt.Printf("After summarization: %v\n", summary)
	fmt.Println("\nYou can now chain this summary as input to other chains!")
}

// mockRetriever implements runtime.Retriever for demo purposes.
type mockRetriever struct{}

func (m *mockRetriever) Retrieve(ctx context.Context, query string, k int) ([]types.Document, error) {
	// Return mock documents
	return []types.Document{
		{
			PageContent: "Machine learning is a subset of AI that enables learning from data.",
			Metadata:    map[string]any{"source": "doc1"},
		},
		{
			PageContent: "Deep learning uses neural networks for complex pattern recognition.",
			Metadata:    map[string]any{"source": "doc2"},
		},
	}, nil
}
