# LLM Providers

All LLM providers implement the `runtime.LLM` interface, allowing them to be used interchangeably.

---

## 📋 Built-in Providers

| Provider | Models | Status | Notes |
|----------|--------|--------|-------|
| **Anthropic** | Claude 3.5 Sonnet, Opus 4.6, Claude 3 Haiku | ✅ Production | Best quality, long context |
| **OpenAI** | GPT-4, GPT-4 Turbo, GPT-3.5 Turbo | ✅ Production | Most popular, great vision |
| **Google Gemini** | Gemini 2.0 Flash, Gemini 1.5 Pro | ✅ Production | Fast, multimodal |
| **Ollama** | llama2, mistral, phi, neural-chat, etc. | ✅ Production | Local inference, free |
| **Mistral** | Mistral Large, Mistral Small | ✅ Production | Good performance, efficient |
| **Groq** | llama-3.3-70b, mixtral-8x7b | ✅ Production | Fastest token generation |
| **Cohere** | Command R+, Command R | ✅ Production | Good for production use |

---

## 🚀 Quick Start

### Using Different Providers

```go
import (
    "context"
    "github.com/rajveer43/goagentflow/provider/anthropic"
    "github.com/rajveer43/goagentflow/provider/openai"
    "github.com/rajveer43/goagentflow/provider/ollama"
)

ctx := context.Background()

// Anthropic Claude (best quality, high cost)
llm := anthropic.New("sk-ant-...", "claude-opus-4-6")
response, _ := llm.Complete(ctx, "What is AI?")

// OpenAI GPT-4 (popular, moderate cost)
llm := openai.New("sk-...", "gpt-4")
response, _ := llm.Complete(ctx, "What is AI?")

// Ollama (local, free)
llm := ollama.New("http://localhost:11434", "mistral")
response, _ := llm.Complete(ctx, "What is AI?")
```

### Using Model Registry

```go
import "github.com/rajveer43/goagentflow/provider"

// Get specific model
model := provider.GetModel("claude-opus-4-6")
fmt.Printf("Cost: $%v per input token\n", model.CostInput)

// List all Anthropic models
models := provider.ListByProvider("anthropic")

// List only streaming models
models := provider.ListCapable("streaming")

// List models with vision support
models := provider.ListCapable("vision")
```

---

## 🔧 Provider-Specific Setup

### Anthropic

```bash
export ANTHROPIC_API_KEY="sk-ant-..."
```

```go
llm := anthropic.New(os.Getenv("ANTHROPIC_API_KEY"), "claude-opus-4-6")
```

**Features:**
- ✅ Streaming
- ✅ Vision (Claude 3.5)
- ✅ Function calling
- Context: 200K tokens (Claude 3.5, Opus 4.6)

### OpenAI

```bash
export OPENAI_API_KEY="sk-..."
```

```go
llm := openai.New(os.Getenv("OPENAI_API_KEY"), "gpt-4")
```

**Features:**
- ✅ Streaming
- ✅ Vision (GPT-4V)
- ✅ Function calling
- Context: 128K tokens (GPT-4 Turbo)

### Google Gemini

```bash
export GOOGLE_API_KEY="..."
```

```go
llm := gemini.New(os.Getenv("GOOGLE_API_KEY"), "gemini-2.0-flash")
```

**Features:**
- ✅ Streaming
- ✅ Vision
- ✅ Function calling
- Context: 1M tokens (Gemini 2.0 Flash)

### Ollama (Local)

```bash
# Start Ollama server
ollama serve

# In another terminal, pull a model
ollama pull mistral
```

```go
llm := ollama.New("http://localhost:11434", "mistral")
```

**Features:**
- ✅ Local inference
- ✅ Free (after download)
- ✅ No API key needed
- Models: llama2, mistral, phi, neural-chat, etc.

### Mistral

```bash
export MISTRAL_API_KEY="..."
```

```go
llm := mistral.New(os.Getenv("MISTRAL_API_KEY"), "mistral-large")
```

**Features:**
- ✅ Streaming
- ✅ Function calling
- Context: 32K tokens

### Groq (Fast Inference)

```bash
export GROQ_API_KEY="..."
```

```go
llm := groq.New(os.Getenv("GROQ_API_KEY"), "llama-3.3-70b-versatile")
```

**Features:**
- ✅ Fastest token generation (>1000 tokens/sec)
- ✅ Streaming
- Context: 8K tokens

### Cohere

```bash
export COHERE_API_KEY="..."
```

```go
llm := cohere.New(os.Getenv("COHERE_API_KEY"), "command-r-plus")
```

**Features:**
- ✅ Streaming
- ✅ Function calling
- Context: 128K tokens

---

## 📊 Comparison Matrix

### Performance

| Provider | Latency | Tokens/sec | Best For |
|----------|---------|-----------|----------|
| Groq | <100ms | >1000 | Fast responses |
| Ollama | <50ms | ~500 | Local inference |
| OpenAI | 300-500ms | ~100 | Production |
| Gemini | 400-600ms | ~120 | Fast API |
| Anthropic | 500-700ms | ~80 | Quality |
| Mistral | 400-500ms | ~100 | Balance |
| Cohere | 300-400ms | ~100 | Production |

### Cost per 1M tokens

| Provider | Input | Output | Total | Notes |
|----------|-------|--------|-------|-------|
| Ollama | Free | Free | Free | Local only |
| Groq | Free | Free | Free | API calls free |
| Mistral Small | $0.14 | $0.42 | $0.56 | Most affordable |
| Cohere | $0.50 | $1.50 | $2.00 | Good value |
| OpenAI (GPT-3.5) | $0.50 | $1.50 | $2.00 | Popular |
| Gemini | $3.50 | $10.50 | $14.00 | Competitive |
| Anthropic (Claude 3.5) | $3.00 | $15.00 | $18.00 | Best quality |

### Features

| Feature | Anthropic | OpenAI | Gemini | Ollama | Mistral | Groq | Cohere |
|---------|-----------|--------|--------|--------|---------|------|--------|
| Streaming | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Vision | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ |
| Function Calls | ✅ | ✅ | ✅ | ❌ | ✅ | ❌ | ✅ |
| Context | 200K | 128K | 1M | Varies | 32K | 8K | 128K |

---

## 🛠️ Adding a Custom LLM Provider

Want to add support for your own LLM service? See [docs/EXTENDING.md](../docs/EXTENDING.md) for full guide.

### Quick Template

```go
// provider/myservice/provider.go
package myservice

import (
    "context"
    "github.com/rajveer43/goagentflow/runtime"
)

type Provider struct {
    apiKey string
}

func New(apiKey string) *Provider {
    return &Provider{apiKey: apiKey}
}

func (p *Provider) Complete(ctx context.Context, prompt string) (string, error) {
    // Your implementation
    return "", nil
}

var _ runtime.LLM = (*Provider)(nil)
```

### Integration Checklist

- [ ] Implement `runtime.LLM` interface
- [ ] Add comprehensive tests
- [ ] Create example in `examples/myservice/`
- [ ] Update `provider/README.md`
- [ ] Add to model registry (optional)
- [ ] Submit PR or publish separate package

---

## 🎯 Choosing a Provider

### Best for Quality
Use **Anthropic Claude**:
- Highest quality responses
- Excellent reasoning
- Long context (200K tokens)
- Best for complex tasks

### Best for Cost
Use **Mistral Small** or **Ollama**:
- Very affordable
- Good quality
- Fast enough for most use cases

### Best for Speed
Use **Groq**:
- Fastest token generation (>1000 tokens/sec)
- Great for real-time applications
- Free API calls

### Best for Local Inference
Use **Ollama**:
- Runs on your machine
- No API costs
- Private (data stays local)
- Requires GPU/CPU power

### Best for Vision
Use **OpenAI (GPT-4V)** or **Anthropic (Claude 3.5)**:
- Excellent image understanding
- High accuracy

### Best for Production
Use **OpenAI** or **Anthropic**:
- Proven reliability
- SLA guarantees
- Great documentation
- Large user base

---

## 📚 Examples

```go
// Simple chat
llm := anthropic.New(apiKey, "claude-opus-4-6")
response, _ := llm.Complete(ctx, "What is machine learning?")

// Streaming responses
llm := openai.New(apiKey, "gpt-4")
stream, _ := llm.Stream(ctx, "Tell me a story...")
for chunk := range stream {
    fmt.Print(chunk)
}

// With memory for context
memory := inmemory.New()
agent := agent.New(llm, memory)
response1, _ := agent.Run(ctx, "My name is Alice")
response2, _ := agent.Run(ctx, "What's my name?")  // Remembers "Alice"

// In RAG pipeline
retriever := retrieval.New(vectorStore, k=3)
ragChain := retrieval.RAG(llm, retriever)
answer, _ := ragChain.Run(ctx, "What is Go?")
```

---

## 🔐 Security Best Practices

### API Keys
```go
// ✅ Good: Use environment variables
apiKey := os.Getenv("ANTHROPIC_API_KEY")
if apiKey == "" {
    panic("ANTHROPIC_API_KEY not set")
}
llm := anthropic.New(apiKey, "claude-opus-4-6")

// ❌ Bad: Hardcoded keys
llm := anthropic.New("sk-ant-abc123", "claude-opus-4-6")
```

### Rate Limiting
```go
// Built-in retry logic with exponential backoff
// Handles rate limits automatically
response, err := llm.Complete(ctx, prompt)  // Retries on rate limit
```

### Timeout Handling
```go
// Always use context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

response, err := llm.Complete(ctx, prompt)
if err == context.DeadlineExceeded {
    fmt.Println("LLM request timed out")
}
```

---

## 📈 Performance Tips

### Token Counting (Cost Estimation)
```go
// Use provider.TokenCounter if available
if counter, ok := llm.(runtime.TokenCounter); ok {
    tokens := counter.CountTokens("Your text here")
    estimatedCost := float64(tokens) * 0.001  // Example: $0.001 per token
    fmt.Printf("Estimated cost: $%.4f\n", estimatedCost)
}
```

### Caching Results
```go
// Avoid re-calling LLM with same prompt
cache := make(map[string]string)
key := "prompt:" + hashString(prompt)

if cached, ok := cache[key]; ok {
    return cached  // Use cached response
}

response, _ := llm.Complete(ctx, prompt)
cache[key] = response
return response
```

### Streaming for Long Responses
```go
// Use streaming for better UX on long responses
stream, _ := llm.Stream(ctx, prompt)
for chunk := range stream {
    fmt.Print(chunk)  // Show tokens as they arrive
}
```

---

## 🐛 Troubleshooting

### "Invalid API Key" Error
```go
// Check API key is set
apiKey := os.Getenv("ANTHROPIC_API_KEY")
if apiKey == "" {
    panic("API key not set!")
}

// Check key format (should start with provider prefix)
// Anthropic: sk-ant-...
// OpenAI: sk-...
// Gemini: (API key without prefix)
```

### Rate Limiting
```go
// Automatic retry with backoff, but you can also:
// - Use Groq for free API calls
// - Use Ollama for local inference
// - Batch requests efficiently
```

### Timeout Issues
```go
// Increase timeout for slow networks
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
defer cancel()
```

---

## 📖 More Info

- [Architecture Overview](../docs/ARCHITECTURE.md)
- [Extending Guide](../docs/EXTENDING.md)
- [Main README](../README.md)
- [Contributing Guide](../CONTRIBUTING.md)

---

Questions? [Start a discussion](https://github.com/rajveer43/goagentflow/discussions)!
