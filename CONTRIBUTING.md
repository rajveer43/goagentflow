# Contributing to goagentflow

Thank you for your interest in contributing to goagentflow! We welcome all contributions—whether it's bug reports, feature requests, documentation improvements, or code.

---

## 🚀 Getting Started

### Prerequisites

- Go 1.18+
- git

### Setup Development Environment

```bash
# Clone the repository
git clone https://github.com/rajveer43/goagentflow.git
cd goagentflow

# Install dependencies
go mod download

# Verify setup
go test ./...
go build ./...
```

---

## 🐛 Reporting Bugs

Please open an issue on GitHub with:

1. **Clear title** - "LLM provider X fails with API error Y"
2. **Description** - What you expected vs. what happened
3. **Minimal reproducible example** - Code snippet showing the issue
4. **Environment** - Go version, OS, library version

Example:
```
Title: Anthropic provider fails with nil context

Description:
When calling llm.Complete(nil, "test"), it panics instead of returning an error.

Expected: Should return an error like "context is required"
Actual: runtime error: invalid memory address

Reproduction:
```go
ctx := context.Context(nil)  // Invalid context
llm := anthropic.New(apiKey, "claude-opus-4-6")
llm.Complete(ctx, "test")  // Panics
```

Environment:
- Go 1.25.0
- macOS
- goagentflow v1.6.0
```

---

## 💡 Feature Requests

Open an issue with:

1. **Clear title** - "Add feature X for use case Y"
2. **Motivation** - Why is this feature needed?
3. **Proposed solution** - How should it work?
4. **Alternatives considered** - Any other approaches?

Example:
```
Title: Add support for Azure OpenAI provider

Motivation:
Many organizations use Azure OpenAI instead of OpenAI directly. Currently goagentflow
only supports the public OpenAI API.

Proposed Solution:
Create provider/azure/provider.go implementing the runtime.LLM interface.

Alternatives:
1. Wrapper around openai provider (would require API key translation)
2. Community-maintained separate package (harder to maintain)
```

---

## 🛠️ Contributing Code

### 1. Fork and Create a Branch

```bash
git checkout -b feature/my-feature
# or for bug fixes:
git checkout -b fix/my-bug
```

### 2. Make Your Changes

**Code Guidelines:**

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Format code: `gofmt -w .`
- Lint: `golangci-lint run ./...` (if available)
- All public functions/types must have doc comments
- Implement interfaces completely (don't partial implementations)

**Example: Adding a Custom LLM Provider**

See `docs/EXTENDING.md` and `provider/README.md` for detailed guide.

```bash
# Create new package
mkdir -p provider/myprovider
touch provider/myprovider/provider.go
```

```go
package myprovider

import (
    "context"
    "github.com/rajveer43/goagentflow/runtime"
)

// Provider implements runtime.LLM for MyService.
type Provider struct {
    apiKey   string
    endpoint string
}

// New creates a new MyService provider.
func New(apiKey, endpoint string) *Provider {
    return &Provider{apiKey: apiKey, endpoint: endpoint}
}

// Complete sends a prompt to MyService.
func (p *Provider) Complete(ctx context.Context, prompt string) (string, error) {
    // 1. Validate inputs
    if prompt == "" {
        return "", runtime.ErrEmptyPrompt
    }
    
    // 2. Make request
    // 3. Handle errors
    // 4. Return response
    return "", nil  // TODO: implement
}

// Verify implements runtime.LLM
var _ runtime.LLM = (*Provider)(nil)
```

**Example: Adding Custom Memory**

See `memory/README.md` for detailed guide.

```bash
mkdir -p memory/mymemory
touch memory/mymemory/memory.go
```

```go
package mymemory

import (
    "context"
    "github.com/rajveer43/goagentflow/runtime"
)

// Memory stores messages in custom backend.
type Memory struct {
    backend string
}

// New creates new custom memory.
func New(backend string) *Memory {
    return &Memory{backend: backend}
}

// AddMessage adds a message.
func (m *Memory) AddMessage(ctx context.Context, msg runtime.Message) error {
    // TODO: implement
    return nil
}

// GetMessages retrieves all messages.
func (m *Memory) GetMessages(ctx context.Context) ([]runtime.Message, error) {
    // TODO: implement
    return nil, nil
}

// Verify implements runtime.Memory
var _ runtime.Memory = (*Memory)(nil)
```

### 3. Write Tests

All new code must include tests:

```bash
# For new LLM provider
touch provider/myprovider/provider_test.go
```

```go
package myprovider_test

import (
    "context"
    "testing"
    "github.com/rajveer43/goagentflow/provider/myprovider"
)

func TestProviderComplete(t *testing.T) {
    provider := myprovider.New("test-key", "http://localhost:8000")
    
    ctx := context.Background()
    response, err := provider.Complete(ctx, "What is AI?")
    
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if response == "" {
        t.Error("expected non-empty response")
    }
}

func TestProviderEmptyPrompt(t *testing.T) {
    provider := myprovider.New("test-key", "http://localhost:8000")
    
    ctx := context.Background()
    _, err := provider.Complete(ctx, "")
    
    if err == nil {
        t.Error("expected error for empty prompt")
    }
}
```

Run tests:

```bash
go test ./...
go test -cover ./...
go test -race ./...  # Check for race conditions
```

### 4. Create an Example (for new public APIs)

If you add new functionality, include an example:

```bash
touch examples/myprovider/main.go
```

```go
package main

import (
    "context"
    "fmt"
    "github.com/rajveer43/goagentflow/provider/myprovider"
)

func main() {
    ctx := context.Background()
    
    // Create provider
    llm := myprovider.New("your-api-key", "http://localhost:8000")
    
    // Use it
    response, err := llm.Complete(ctx, "What is AI?")
    if err != nil {
        panic(err)
    }
    
    fmt.Println(response)
}
```

### 5. Update Documentation

Update relevant documentation files:

- **New provider?** Update `provider/README.md`
- **New memory type?** Update `memory/README.md`
- **New chain?** Update `chains/README.md`
- **New feature?** Add to relevant `docs/*.md`

### 6. Commit Changes

Use conventional commit format:

```bash
# Features
git commit -m "feat: add MyService LLM provider"

# Bug fixes
git commit -m "fix: handle nil context in Anthropic provider"

# Documentation
git commit -m "docs: update LLM provider guide"

# Tests
git commit -m "test: add integration tests for MyService"

# Performance
git commit -m "perf: optimize vector similarity computation"

# Refactoring
git commit -m "refactor: simplify memory interface implementation"
```

### 7. Push and Open Pull Request

```bash
git push origin feature/my-feature
```

Then [open a pull request](https://github.com/rajveer43/goagentflow/pulls) with:

1. **Clear title** - "Add MyService LLM provider"
2. **Description** - What does this PR do?
3. **Related issues** - Links to issues (e.g., "Fixes #123")
4. **Testing** - How to test the changes
5. **Screenshots** (if UI changes)

Example PR description:
```markdown
## Description
Adds support for MyService as an LLM provider, enabling fast local inference.

## Motivation
MyService is widely used for local LLMs. This PR allows goagentflow users to use it.

## Changes
- Added `provider/myprovider/` with full implementation
- Added comprehensive tests in `provider/myprovider/*_test.go`
- Added example in `examples/myprovider/main.go`
- Updated `provider/README.md` with new provider

## Testing
```bash
go test ./provider/myprovider/...
go run examples/myprovider/main.go
```

## Related Issues
Fixes #123
```

---

## 📋 PR Review Process

Once you open a PR:

1. **Automated checks** run (tests, build, linting)
2. **Code review** - We'll review your changes
3. **Feedback** - We may request changes
4. **Merge** - Once approved, your PR is merged!

**Expectations:**
- Tests must pass
- Code must build cleanly
- No external dependencies in core packages
- Examples must work
- Documentation updated

---

## 🎯 What We're Looking For

We especially welcome contributions in these areas:

1. **New LLM Providers**
   - See `provider/README.md` for list and guide

2. **New Memory Types**
   - See `memory/README.md` for list and guide

3. **New Chains**
   - See `chains/README.md` for list and guide

4. **Bug Fixes**
   - Any issue marked "bug" is fair game

5. **Documentation**
   - Improving examples
   - Clarifying docs
   - Adding tutorials

6. **Tests**
   - Increasing coverage
   - Adding edge cases
   - Performance tests

7. **Performance**
   - Optimizing hot paths
   - Reducing allocations
   - Better algorithms

---

## 💬 Getting Help

- **Questions?** [Open a discussion](https://github.com/rajveer43/goagentflow/discussions)
- **Stuck?** Comment on the issue or PR
- **Design feedback?** Start a discussion before implementing

---

## 📝 Style Guide

### Go Style

Follow [Effective Go](https://golang.org/doc/effective_go):

```go
// ✅ Good: Clear, concise, idiomatic Go
func (p *Provider) Complete(ctx context.Context, prompt string) (string, error) {
    if prompt == "" {
        return "", errors.New("prompt cannot be empty")
    }
    
    // Do work
    result, err := p.call(ctx, prompt)
    if err != nil {
        return "", fmt.Errorf("api call failed: %w", err)
    }
    
    return result, nil
}

// ❌ Bad: Verbose, confusing error handling
func (p *Provider) Complete(ctx context.Context, prompt string) (string, error) {
    if len(prompt) == 0 {
        err := errors.New("prompt is empty!")
        return "", err
    }
    
    result, err := p.call(ctx, prompt)
    if err != nil {
        fmt.Println("ERROR: " + err.Error())
        return "", err
    }
    
    return result, nil
}
```

### Comments

```go
// ✅ Good: Clear, explains why
// Use buffer instead of unbuffered channel to avoid goroutine leak
ch := make(chan bool, 1)

// ❌ Bad: Obvious or unclear
// This is a channel
ch := make(chan bool)
```

### Naming

```go
// ✅ Good: Clear names
type Provider struct {
    apiKey string
    timeout time.Duration
}

// ❌ Bad: Abbreviations, unclear
type Prov struct {
    key string
    t time.Duration
}
```

---

## 🤝 Code of Conduct

- Be respectful and inclusive
- Assume good intent
- Accept constructive criticism
- Focus on the code, not the person

---

## 📄 License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

## Questions?

- 💬 Start a [discussion](https://github.com/rajveer43/goagentflow/discussions)
- 🐛 [Open an issue](https://github.com/rajveer43/goagentflow/issues)
- 📧 Reach out to the maintainer

---

Thank you for making goagentflow better! 🙏
