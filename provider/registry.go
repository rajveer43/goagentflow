package provider

import (
	"fmt"

	"github.com/rajveer43/goagentflow/runtime"
)

// Models is a global registry of known LLM models.
// Enables model discovery, cost estimation, and capability filtering.
var Models = map[string]runtime.ModelInfo{
	// OpenAI Models
	"gpt-4o": {
		Name:            "gpt-4o",
		Provider:        "openai",
		MaxTokens:       128000,
		ContextSize:     128000,
		CostPer1KInput:  0.005,
		CostPer1KOutput: 0.015,
		Capabilities:    []string{"vision", "function_calling", "json_mode", "streaming"},
		ReleaseDate:     "2024-11-20",
	},
	"gpt-4o-mini": {
		Name:            "gpt-4o-mini",
		Provider:        "openai",
		MaxTokens:       128000,
		ContextSize:     128000,
		CostPer1KInput:  0.00015,
		CostPer1KOutput: 0.0006,
		Capabilities:    []string{"vision", "function_calling", "json_mode", "streaming"},
		ReleaseDate:     "2024-07-18",
	},
	"gpt-4-turbo": {
		Name:            "gpt-4-turbo",
		Provider:        "openai",
		MaxTokens:       4096,
		ContextSize:     128000,
		CostPer1KInput:  0.01,
		CostPer1KOutput: 0.03,
		Capabilities:    []string{"vision", "function_calling", "json_mode", "streaming"},
		ReleaseDate:     "2023-11-06",
	},
	"gpt-3.5-turbo": {
		Name:            "gpt-3.5-turbo",
		Provider:        "openai",
		MaxTokens:       4096,
		ContextSize:     16384,
		CostPer1KInput:  0.0005,
		CostPer1KOutput: 0.0015,
		Capabilities:    []string{"function_calling", "streaming"},
		ReleaseDate:     "2023-03-15",
	},

	// Anthropic Models
	"claude-opus-4-6": {
		Name:            "claude-opus-4-6",
		Provider:        "anthropic",
		MaxTokens:       4096,
		ContextSize:     200000,
		CostPer1KInput:  0.015,
		CostPer1KOutput: 0.075,
		Capabilities:    []string{"vision", "function_calling", "streaming"},
		ReleaseDate:     "2025-02-27",
	},
	"claude-sonnet-4-6": {
		Name:            "claude-sonnet-4-6",
		Provider:        "anthropic",
		MaxTokens:       4096,
		ContextSize:     200000,
		CostPer1KInput:  0.003,
		CostPer1KOutput: 0.015,
		Capabilities:    []string{"vision", "function_calling", "streaming"},
		ReleaseDate:     "2025-02-27",
	},
	"claude-haiku-4-5-20251001": {
		Name:            "claude-haiku-4-5-20251001",
		Provider:        "anthropic",
		MaxTokens:       4096,
		ContextSize:     200000,
		CostPer1KInput:  0.00025,
		CostPer1KOutput: 0.00125,
		Capabilities:    []string{"vision", "function_calling", "streaming"},
		ReleaseDate:     "2025-10-01",
	},

	// Google Gemini Models
	"gemini-2.0-flash": {
		Name:            "gemini-2.0-flash",
		Provider:        "gemini",
		MaxTokens:       8000,
		ContextSize:     1000000,
		CostPer1KInput:  0.000075,
		CostPer1KOutput: 0.0003,
		Capabilities:    []string{"vision", "function_calling", "streaming", "json_mode"},
		ReleaseDate:     "2024-12-19",
	},
	"gemini-1.5-pro": {
		Name:            "gemini-1.5-pro",
		Provider:        "gemini",
		MaxTokens:       8000,
		ContextSize:     2000000,
		CostPer1KInput:  0.00125,
		CostPer1KOutput: 0.005,
		Capabilities:    []string{"vision", "function_calling", "streaming", "json_mode"},
		ReleaseDate:     "2024-05-14",
	},
	"gemini-1.5-flash": {
		Name:            "gemini-1.5-flash",
		Provider:        "gemini",
		MaxTokens:       8000,
		ContextSize:     1000000,
		CostPer1KInput:  0.000075,
		CostPer1KOutput: 0.0003,
		Capabilities:    []string{"vision", "function_calling", "streaming", "json_mode"},
		ReleaseDate:     "2024-09-24",
	},

	// Mistral Models
	"mistral-large-latest": {
		Name:            "mistral-large-latest",
		Provider:        "mistral",
		MaxTokens:       8192,
		ContextSize:     32000,
		CostPer1KInput:  0.0009,
		CostPer1KOutput: 0.0027,
		Capabilities:    []string{"function_calling", "streaming", "json_mode"},
		ReleaseDate:     "2024-09-17",
	},
	"mistral-medium-latest": {
		Name:            "mistral-medium-latest",
		Provider:        "mistral",
		MaxTokens:       8192,
		ContextSize:     32000,
		CostPer1KInput:  0.00027,
		CostPer1KOutput: 0.00081,
		Capabilities:    []string{"function_calling", "streaming"},
		ReleaseDate:     "2024-05-08",
	},
	"mistral-small-latest": {
		Name:            "mistral-small-latest",
		Provider:        "mistral",
		MaxTokens:       8192,
		ContextSize:     32000,
		CostPer1KInput:  0.000027,
		CostPer1KOutput: 0.000081,
		Capabilities:    []string{"function_calling", "streaming"},
		ReleaseDate:     "2024-11-22",
	},

	// Groq Models (very fast)
	"llama-3.3-70b-versatile": {
		Name:            "llama-3.3-70b-versatile",
		Provider:        "groq",
		MaxTokens:       8192,
		ContextSize:     8192,
		CostPer1KInput:  0.0,
		CostPer1KOutput: 0.0, // Free tier available
		Capabilities:    []string{"streaming"},
		ReleaseDate:     "2024-11-07",
	},
	"mixtral-8x7b-32768": {
		Name:            "mixtral-8x7b-32768",
		Provider:        "groq",
		MaxTokens:       32768,
		ContextSize:     32768,
		CostPer1KInput:  0.0,
		CostPer1KOutput: 0.0, // Free tier
		Capabilities:    []string{"streaming"},
		ReleaseDate:     "2024-01-09",
	},
	"gemma2-9b-it": {
		Name:            "gemma2-9b-it",
		Provider:        "groq",
		MaxTokens:       8192,
		ContextSize:     8192,
		CostPer1KInput:  0.0,
		CostPer1KOutput: 0.0,
		Capabilities:    []string{"streaming"},
		ReleaseDate:     "2024-06-27",
	},

	// Cohere Models
	"command-r-plus": {
		Name:            "command-r-plus",
		Provider:        "cohere",
		MaxTokens:       4096,
		ContextSize:     128000,
		CostPer1KInput:  0.003,
		CostPer1KOutput: 0.015,
		Capabilities:    []string{"function_calling", "streaming", "json_mode"},
		ReleaseDate:     "2024-03-22",
	},
	"command-r": {
		Name:            "command-r",
		Provider:        "cohere",
		MaxTokens:       4096,
		ContextSize:     128000,
		CostPer1KInput:  0.0005,
		CostPer1KOutput: 0.0015,
		Capabilities:    []string{"function_calling", "streaming", "json_mode"},
		ReleaseDate:     "2024-03-15",
	},

	// Ollama Local Models (no cost, requires local setup)
	"llama2": {
		Name:            "llama2",
		Provider:        "ollama",
		MaxTokens:       4096,
		ContextSize:     4096,
		CostPer1KInput:  0.0,
		CostPer1KOutput: 0.0,
		Capabilities:    []string{"streaming"},
		ReleaseDate:     "2023-07-18",
	},
	"llama3.2": {
		Name:            "llama3.2",
		Provider:        "ollama",
		MaxTokens:       4096,
		ContextSize:     8000,
		CostPer1KInput:  0.0,
		CostPer1KOutput: 0.0,
		Capabilities:    []string{"streaming"},
		ReleaseDate:     "2024-09-12",
	},
	"codellama": {
		Name:            "codellama",
		Provider:        "ollama",
		MaxTokens:       16000,
		ContextSize:     16000,
		CostPer1KInput:  0.0,
		CostPer1KOutput: 0.0,
		Capabilities:    []string{"streaming"},
		ReleaseDate:     "2023-08-24",
	},
	"mistral": {
		Name:            "mistral",
		Provider:        "ollama",
		MaxTokens:       8192,
		ContextSize:     8192,
		CostPer1KInput:  0.0,
		CostPer1KOutput: 0.0,
		Capabilities:    []string{"streaming"},
		ReleaseDate:     "2023-12-26",
	},
	"neural-chat": {
		Name:            "neural-chat",
		Provider:        "ollama",
		MaxTokens:       4096,
		ContextSize:     4096,
		CostPer1KInput:  0.0,
		CostPer1KOutput: 0.0,
		Capabilities:    []string{"streaming"},
		ReleaseDate:     "2023-06-15",
	},
}

// GetModel returns model information by name.
// Returns (ModelInfo, false) if model not found.
func GetModel(name string) (runtime.ModelInfo, bool) {
	model, ok := Models[name]
	return model, ok
}

// ListByProvider returns all models from a specific provider.
func ListByProvider(provider string) []runtime.ModelInfo {
	var results []runtime.ModelInfo
	for _, model := range Models {
		if model.Provider == provider {
			results = append(results, model)
		}
	}
	return results
}

// ListCapable returns all models with a specific capability.
func ListCapable(capability string) []runtime.ModelInfo {
	var results []runtime.ModelInfo
	for _, model := range Models {
		for _, cap := range model.Capabilities {
			if cap == capability {
				results = append(results, model)
				break
			}
		}
	}
	return results
}

// ListAll returns all registered models.
func ListAll() []runtime.ModelInfo {
	var results []runtime.ModelInfo
	for _, model := range Models {
		results = append(results, model)
	}
	return results
}

// Providers returns a list of unique provider names.
func Providers() []string {
	providerSet := make(map[string]bool)
	for _, model := range Models {
		providerSet[model.Provider] = true
	}
	var providers []string
	for provider := range providerSet {
		providers = append(providers, provider)
	}
	return providers
}

// AvailableCapabilities returns a list of unique capabilities across all models.
func AvailableCapabilities() []string {
	capSet := make(map[string]bool)
	for _, model := range Models {
		for _, cap := range model.Capabilities {
			capSet[cap] = true
		}
	}
	var caps []string
	for cap := range capSet {
		caps = append(caps, cap)
	}
	return caps
}

// RegisterModel adds or updates a model in the registry.
func RegisterModel(info runtime.ModelInfo) error {
	if info.Name == "" {
		return fmt.Errorf("model name cannot be empty")
	}
	if info.Provider == "" {
		return fmt.Errorf("model provider cannot be empty")
	}
	Models[info.Name] = info
	return nil
}
