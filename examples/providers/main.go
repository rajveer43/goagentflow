package main

import (
	"context"
	"fmt"
	"log"

	"github.com/rajveer43/goagentflow/provider"
	"github.com/rajveer43/goagentflow/provider/anthropic"
	"github.com/rajveer43/goagentflow/provider/cohere"
	"github.com/rajveer43/goagentflow/provider/gemini"
	"github.com/rajveer43/goagentflow/provider/groq"
	"github.com/rajveer43/goagentflow/provider/mistral"
	"github.com/rajveer43/goagentflow/provider/ollama"
	"github.com/rajveer43/goagentflow/provider/openai"
)

// Example demonstrates available LLM providers and the model registry.
// This example shows:
// 1. Creating clients for each provider
// 2. Using the model registry to discover available models
// 3. Querying models by provider, capability, etc.
// 4. Generating responses from different providers
func main() {
	ctx := context.Background()

	// Example 1: Create clients for each provider
	fmt.Println("=== LLM Provider Clients ===\n")

	openaiClient := openai.New("your-openai-api-key", "gpt-4o")
	fmt.Printf("OpenAI client created for model: %s\n", "gpt-4o")

	anthropicClient := anthropic.New("your-anthropic-api-key", "claude-opus-4-6")
	fmt.Printf("Anthropic client created for model: %s\n", "claude-opus-4-6")

	geminiClient := gemini.New("your-google-api-key", "gemini-2.0-flash")
	fmt.Printf("Gemini client created for model: %s\n", "gemini-2.0-flash")

	ollamaClient := ollama.New("http://localhost:11434", "llama3.2")
	fmt.Printf("Ollama client created for model: %s\n", "llama3.2")

	mistralClient := mistral.New("your-mistral-api-key", "mistral-large-latest")
	fmt.Printf("Mistral client created for model: %s\n", "mistral-large-latest")

	groqClient := groq.New("your-groq-api-key", "llama-3.3-70b-versatile")
	fmt.Printf("Groq client created for model: %s\n", "llama-3.3-70b-versatile")

	cohereClient := cohere.New("your-cohere-api-key", "command-r-plus")
	fmt.Printf("Cohere client created for model: %s\n\n", "command-r-plus")

	// Example 2: Use model registry to discover models
	fmt.Println("=== Model Registry Discovery ===\n")

	// List all providers
	fmt.Println("Available Providers:")
	providers := provider.Providers()
	for _, p := range providers {
		fmt.Printf("  - %s\n", p)
	}
	fmt.Println()

	// List models by provider
	fmt.Println("OpenAI Models:")
	openaiModels := provider.ListByProvider("openai")
	for _, model := range openaiModels {
		fmt.Printf("  - %s (context: %d tokens, cost: $%.4f/$%.4f per 1K)\n",
			model.Name, model.ContextSize, model.CostPer1KInput, model.CostPer1KOutput)
	}
	fmt.Println()

	// List models with specific capability
	fmt.Println("Vision-Capable Models:")
	visionModels := provider.ListCapable("vision")
	for _, model := range visionModels {
		fmt.Printf("  - %s (%s)\n", model.Name, model.Provider)
	}
	fmt.Println()

	// Example 3: Get detailed info about a model
	fmt.Println("=== Model Details ===\n")
	if modelInfo, ok := provider.GetModel("gpt-4o"); ok {
		fmt.Printf("Model: %s\n", modelInfo.Name)
		fmt.Printf("Provider: %s\n", modelInfo.Provider)
		fmt.Printf("Context Size: %d tokens\n", modelInfo.ContextSize)
		fmt.Printf("Max Output: %d tokens\n", modelInfo.MaxTokens)
		fmt.Printf("Input Cost: $%.6f per 1K tokens\n", modelInfo.CostPer1KInput)
		fmt.Printf("Output Cost: $%.6f per 1K tokens\n", modelInfo.CostPer1KOutput)
		fmt.Printf("Capabilities: %v\n", modelInfo.Capabilities)
		fmt.Printf("Released: %s\n\n", modelInfo.ReleaseDate)
	}

	// Example 4: Generate completions from different providers
	fmt.Println("=== Generating Responses ===\n")

	prompt := "What is machine learning?"

	// OpenAI
	fmt.Println("OpenAI (gpt-4o):")
	response, err := openaiClient.Complete(ctx, prompt)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("  %s\n\n", response)
	}

	// Anthropic
	fmt.Println("Anthropic (claude-opus-4-6):")
	response, err = anthropicClient.Complete(ctx, prompt)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("  %s\n\n", response)
	}

	// Gemini
	fmt.Println("Google Gemini (gemini-2.0-flash):")
	response, err = geminiClient.Complete(ctx, prompt)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("  %s\n\n", response)
	}

	// Ollama (local)
	fmt.Println("Ollama (llama3.2) - local:")
	response, err = ollamaClient.Complete(ctx, prompt)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("  %s\n\n", response)
	}

	// Mistral
	fmt.Println("Mistral (mistral-large-latest):")
	response, err = mistralClient.Complete(ctx, prompt)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("  %s\n\n", response)
	}

	// Groq
	fmt.Println("Groq (llama-3.3-70b-versatile):")
	response, err = groqClient.Complete(ctx, prompt)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("  %s\n\n", response)
	}

	// Cohere
	fmt.Println("Cohere (command-r-plus):")
	response, err = cohereClient.Complete(ctx, prompt)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("  %s\n\n", response)
	}

	// Example 5: Streaming responses
	fmt.Println("=== Streaming Responses ===\n")
	fmt.Println("OpenAI Streaming (gpt-4o):")
	tokens, errors := openaiClient.Stream(ctx, prompt)
	for {
		select {
		case token, ok := <-tokens:
			if !ok {
				fmt.Println()
				tokens = nil
				continue
			}
			fmt.Print(token)
		case err, ok := <-errors:
			if ok {
				fmt.Printf("Error: %v\n", err)
			}
			if tokens == nil {
				return
			}
		}
		if tokens == nil {
			break
		}
	}
	fmt.Println()

	// Example 6: Cost comparison
	fmt.Println("\n=== Cost Comparison ===\n")
	fmt.Println("Estimated cost for 1000 input tokens + 1000 output tokens:")
	models := []string{"gpt-4o", "claude-opus-4-6", "gemini-2.0-flash", "mistral-large-latest", "command-r-plus"}
	for _, modelName := range models {
		if info, ok := provider.GetModel(modelName); ok {
			inputCost := (float64(1000) / 1000) * info.CostPer1KInput
			outputCost := (float64(1000) / 1000) * info.CostPer1KOutput
			totalCost := inputCost + outputCost
			fmt.Printf("  %s: $%.6f\n", info.Name, totalCost)
		}
	}
	fmt.Println()

	// Example 7: List all capabilities
	fmt.Println("=== Available Capabilities ===")
	capabilities := provider.AvailableCapabilities()
	for _, cap := range capabilities {
		models := provider.ListCapable(cap)
		fmt.Printf("  %s: %d models\n", cap, len(models))
	}
}
