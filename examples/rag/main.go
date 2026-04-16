package main

import (
	"context"
	"fmt"
	"log"

	"github.com/rajveer43/goagentflow/embeddings"
	"github.com/rajveer43/goagentflow/loader"
	"github.com/rajveer43/goagentflow/provider/anthropic"
	"github.com/rajveer43/goagentflow/retrieval"
	"github.com/rajveer43/goagentflow/types"
	"github.com/rajveer43/goagentflow/vectorstore/memory"
)

// Example demonstrates a complete RAG (Retrieval-Augmented Generation) workflow.
// This example:
// 1. Loads some sample documents
// 2. Chunks them into smaller pieces
// 3. Embeds the chunks using OpenAI embeddings
// 4. Stores them in an in-memory vector store
// 5. Retrieves relevant documents for a query
// 6. Uses an LLM to generate an answer based on the retrieved context
func main() {
	ctx := context.Background()

	// Step 1: Create some sample documents
	docs := []types.Document{
		{
			PageContent: "Go is a compiled, statically typed programming language. It was created by Google and first released in 2009.",
			Metadata: map[string]any{
				"source": "wikipedia",
				"topic":  "golang",
			},
		},
		{
			PageContent: "Go's concurrency model is based on goroutines and channels. Goroutines are lightweight threads managed by the Go runtime.",
			Metadata: map[string]any{
				"source": "golang.org",
				"topic":  "concurrency",
			},
		},
		{
			PageContent: "The Go standard library is comprehensive and includes packages for networking, encryption, JSON processing, and more.",
			Metadata: map[string]any{
				"source": "golang.org",
				"topic":  "stdlib",
			},
		},
	}

	// Step 2: Split documents into chunks (optional, for larger documents)
	splitter := loader.NewCharacterSplitter(500, 50)
	chunkedDocs, err := splitter.SplitDocuments(ctx, docs)
	if err != nil {
		log.Fatalf("Failed to split documents: %v", err)
	}
	fmt.Printf("Split %d documents into %d chunks\n\n", len(docs), len(chunkedDocs))

	// Step 3: Create embedder
	// Note: In a real application, provide a valid OpenAI API key
	embeddingClient := embeddings.New("your-openai-api-key")
	fmt.Printf("Using embedder with dimension: %d\n\n", embeddingClient.Dimension())

	// Step 4: Embed the documents
	embeddings := make([][]float32, len(chunkedDocs))
	for i, doc := range chunkedDocs {
		emb, err := embeddingClient.Embed(ctx, doc.PageContent)
		if err != nil {
			log.Fatalf("Failed to embed document %d: %v", i, err)
		}
		embeddings[i] = emb
	}
	fmt.Printf("Embedded %d documents\n\n", len(embeddings))

	// Step 5: Create and populate vector store
	vectorStore := memory.New()
	ids, err := vectorStore.Add(ctx, chunkedDocs, embeddings)
	if err != nil {
		log.Fatalf("Failed to add documents to vector store: %v", err)
	}
	fmt.Printf("Added %d documents to vector store with IDs: %v\n\n", len(ids), ids[:min(len(ids), 3)])

	// Step 6: Create retriever
	retriever := retrieval.New(vectorStore, embeddingClient)

	// Step 7: Retrieve documents for a query
	query := "What is Go's concurrency model?"
	retrievedDocs, err := retriever.Retrieve(ctx, query, 2)
	if err != nil {
		log.Fatalf("Failed to retrieve documents: %v", err)
	}
	fmt.Printf("Query: %s\n", query)
	fmt.Printf("Retrieved %d documents:\n", len(retrievedDocs))
	for i, doc := range retrievedDocs {
		fmt.Printf("  %d. %s\n", i+1, doc.PageContent[:min(len(doc.PageContent), 80)]+"...")
	}
	fmt.Println()

	// Step 8: Create LLM client (using Anthropic as example)
	// Note: Provide valid API key for real usage
	llmClient := anthropic.New("your-anthropic-api-key", "claude-3-haiku")

	// Step 9: Create RAG chain
	ragChain := retrieval.NewRAGChain(retriever, llmClient, 2)

	// Step 10: Generate answer using RAG
	answer, err := ragChain.Run(ctx, query)
	if err != nil {
		log.Fatalf("Failed to generate answer: %v", err)
	}
	fmt.Printf("Generated answer:\n%v\n\n", answer)

	// Step 11: Try streaming response (optional)
	fmt.Println("Streaming example (would show token-by-token output with real LLM):")
	streamingChain := retrieval.NewStreamingRAGChain(retriever, llmClient, 2)
	tokens, errors := streamingChain.Stream(ctx, query)

	for {
		select {
		case token, ok := <-tokens:
			if !ok {
				fmt.Println()
				break
			}
			fmt.Print(token)
		case err, ok := <-errors:
			if ok {
				fmt.Printf("Error: %v\n", err)
			}
			return
		}
	}

	// Step 12: Vector store operations
	size, _ := vectorStore.Size(ctx)
	fmt.Printf("\nVector store now contains %d documents\n", size)

	// Clean up
	vectorStore.Clear(ctx)
	fmt.Println("Vector store cleared")
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
