# Chains

Chains are composable pipeline steps that implement the `runtime.Chain` interface. Combine them to build complex workflows.

---

## 📋 Available Chains

| Chain | Input | Output | Use Case |
|-------|-------|--------|----------|
| **QA** | Question (string) | Answer + Sources | Document Q&A |
| **Summarization** | Documents | Summary | Summarize articles/docs |
| **SQL** | Natural language query | SQL query | Database queries |
| **Agent** | User message | Response | Multi-turn conversations |

---

## 🚀 Quick Start

### Simple Question Answering

```go
import "github.com/rajveer43/goagentflow/chains/qa"

// Create QA chain
retriever := retrieval.New(vectorStore, k=3)
qaChain := qa.New(retriever, llm, numResults=3)

// Ask question
answer, _ := qaChain.Run(ctx, "What is machine learning?")
```

### Document Summarization

```go
import "github.com/rajveer43/goagentflow/chains/summarization"

// Create summarizer
summarizer := summarization.New(llm, summarization.StuffStrategy)

// Summarize
summary, _ := summarizer.Run(ctx, documents)
```

### Natural Language to SQL

```go
import "github.com/rajveer43/goagentflow/chains/sql"

// Create SQL chain with schema
schema := "users (id, name, email), orders (id, user_id, total)"
sqlChain := sql.New(llm, schema)

// Generate SQL
sqlQuery, _ := sqlChain.Run(ctx, "Top 5 customers by spending")
```

### Multi-Turn Agent

```go
import "github.com/rajveer43/goagentflow/chains/agent"

// Create agent
agent := agent.New(llm, memory)
agent.RegisterTool(webSearchTool)
agent.RegisterTool(calculatorTool)

// Multi-turn conversation
response1, _ := agent.Run(ctx, "What's the capital of France?")
response2, _ := agent.Run(ctx, "What's its population?")
response3, _ := agent.Run(ctx, "Multiply by 2")
```

---

## 📖 Detailed Chain Documentation

### QA Chain

**Purpose:** Answer questions about documents using retrieval + LLM

```go
qaChain := qa.New(retriever, llm, numResults=3)
answer, _ := qaChain.Run(ctx, "What is Go?")
```

**What it does:**
1. Takes your question
2. Retrieves top 3 relevant documents
3. Builds context: "Question: {question}\nRelevant docs: {docs}"
4. Sends to LLM
5. Returns answer

**Configuration:**
```go
qaChain := qa.New(
    retriever,           // How to retrieve documents
    llm,                 // LLM for answering
    3,                   // Top K documents to retrieve
)
```

**Example: Document Q&A System**
```go
// Setup
documents := loader.LoadPDF("mybook.pdf")
vectorStore.Add(ctx, documents)
retriever := retrieval.New(vectorStore, 3)
qaChain := qa.New(retriever, llm, 3)

// Use
answer, _ := qaChain.Run(ctx, "What chapter covers authentication?")
```

### Summarization Chain

**Purpose:** Summarize documents (handles long documents)

```go
summarizer := summarization.New(llm, summarization.StuffStrategy)
summary, _ := summarizer.Run(ctx, documents)
```

**Two Strategies:**

#### 1. Stuff Strategy (Fast)
- Concatenate all documents
- Send to LLM once
- Best for: Short documents

```go
summarizer := summarization.New(llm, summarization.StuffStrategy)
summary, _ := summarizer.Run(ctx, documents)  // ~2 LLM calls
```

#### 2. MapReduce Strategy (For Large Documents)
- Summarize each chunk independently
- Combine summaries
- Summarize combined result
- Best for: Long documents that exceed context

```go
summarizer := summarization.New(llm, summarization.MapReduceStrategy)
summary, _ := summarizer.Run(ctx, documents)  // Many LLM calls
```

**Configuration:**
```go
summarizer := summarization.New(llm, summarization.StuffStrategy)
summarizer.SetChunkSize(500, 50)  // 500 char chunks, 50 char overlap
summarizer.SetMaxConcurrency(5)   // Parallel requests
```

**Example: Multi-Document Summarization**
```go
// Load multiple documents
docs := []runtime.Document{
    {Content: "Article 1..."},
    {Content: "Article 2..."},
    {Content: "Article 3..."},
}

summarizer := summarization.New(llm, summarization.MapReduceStrategy)
summary, _ := summarizer.Run(ctx, docs)
fmt.Println(summary)  // Unified summary of all articles
```

### SQL Chain

**Purpose:** Convert natural language to SQL queries

```go
schema := "users (id, name, email), orders (id, user_id, total)"
sqlChain := sql.New(llm, schema)
sqlQuery, _ := sqlChain.Run(ctx, "Top 5 customers by spending")
```

**What it does:**
1. Takes natural language query
2. Includes database schema
3. Sends to LLM: "Convert this to SQL: {query}\n\nSchema: {schema}"
4. Returns SQL query

**Configuration:**
```go
sqlChain := sql.New(llm, schema)
sqlChain.SetDialect("PostgreSQL")   // SQL flavor
sqlChain.SetIncludeSchema(true)     // Include full schema
```

**Example: Interactive Database Query**
```go
// Database schema
schema := `
users (id INTEGER, name VARCHAR, email VARCHAR)
orders (id INTEGER, user_id INTEGER, total DECIMAL, created_at DATE)
`

sqlChain := sql.New(llm, schema)
sqlChain.SetDialect("PostgreSQL")

// User queries
queries := []string{
    "Find all users named John",
    "Top 10 users by total spending",
    "Orders from last month",
    "Users with more than 5 orders",
}

for _, query := range queries {
    sql, _ := sqlChain.Run(ctx, query)
    fmt.Println(sql)
    
    // Execute and show results...
}
```

### Agent Chain

**Purpose:** Multi-turn conversations with tools and memory

```go
agent := agent.New(llm, memory)
agent.RegisterTool(webSearch)
agent.RegisterTool(calculator)

response, _ := agent.Run(ctx, "What's the GDP of France times 2?")
```

**What it does:**
1. Takes user message
2. Adds conversation history (from memory)
3. Lists available tools
4. Sends to LLM
5. If LLM requests tool, executes it
6. Repeats until LLM gives final answer
7. Stores message + response in memory

**Configuration:**
```go
agent := agent.New(llm, memory)
agent.RegisterTool(&Tool{
    Name: "calculator",
    Description: "Perform math",
    Run: func(ctx context.Context, input string) (string, error) {
        // Your implementation
        return result, nil
    },
})
agent.SetMaxSteps(5)      // Max iterations
agent.SetSystemPrompt("") // Custom system message
```

**Example: Research Agent**
```go
// Tools
webSearch := &Tool{
    Name: "web_search",
    Description: "Search the web for information",
    Run: func(ctx context.Context, query string) (string, error) {
        // Use real search API
        return results, nil
    },
}

summarizer := &Tool{
    Name: "summarize",
    Description: "Summarize long text",
    Run: func(ctx context.Context, text string) (string, error) {
        // Summarize using chain
        return summary, nil
    },
}

// Agent
memory := inmemory.New()
agent := agent.New(llm, memory)
agent.RegisterTool(webSearch)
agent.RegisterTool(summarizer)
agent.SetMaxSteps(10)

// Multi-turn conversation
response1, _ := agent.Run(ctx, "Tell me about Go programming language")
response2, _ := agent.Run(ctx, "What are its main benefits?")
response3, _ := agent.Run(ctx, "Compare with Rust")
```

---

## 🔗 Chain Composition (Pipeline)

Compose chains into workflows using `runtime.NewChainPipeline`:

```go
import "github.com/rajveer43/goagentflow/runtime"

// Create individual chains
splitter := loader.NewCharacterSplitter(500, 50)
summarizer := summarization.New(llm, summarization.StuffStrategy)
qaChain := qa.New(retriever, llm, 3)

// Compose into pipeline
pipeline := runtime.NewChainPipeline(
    splitter,      // Step 1: Split documents
    summarizer,    // Step 2: Summarize
    qaChain,       // Step 3: Answer questions
)

// Execute
result, _ := pipeline.Run(ctx, documents)
```

**Example: Document Processing Pipeline**

```go
// Pipeline: Load → Split → Summarize → Store → QA

// Step 1: Load documents
loader := loader.New()
documents, _ := loader.LoadPDF("report.pdf")

// Step 2: Split into chunks
splitter := loader.NewRecursiveSplitter(1000, 200)
chunks, _ := splitter.Split(ctx, documents)

// Step 3: Embed and store
embedder := openai.New(apiKey, "text-embedding-3-small")
vectorStore := memory.New(embedder)
vectorStore.Add(ctx, chunks)

// Step 4: Create chains
summarizer := summarization.New(llm, summarization.StuffStrategy)
qaChain := qa.New(retrieval.New(vectorStore, 3), llm, 3)

// Step 5: Compose pipeline
pipeline := runtime.NewChainPipeline(summarizer, qaChain)

// Step 6: Execute
summary, _ := pipeline.Run(ctx, chunks)
```

---

## 🛠️ Creating Custom Chains

Want to add custom chain logic? See [docs/EXTENDING.md](../docs/EXTENDING.md).

### Quick Template

```go
// chains/custom/custom.go
package custom

import (
    "context"
    "github.com/rajveer43/goagentflow/runtime"
)

type MyChain struct {
    llm runtime.LLM
}

func New(llm runtime.LLM) *MyChain {
    return &MyChain{llm: llm}
}

func (c *MyChain) Run(ctx context.Context, input any) (any, error) {
    // 1. Validate and type-assert input
    docs, ok := input.([]runtime.Document)
    if !ok {
        return nil, fmt.Errorf("expected []Document")
    }
    
    // 2. Process
    var results []string
    for _, doc := range docs {
        result, _ := c.process(ctx, doc)
        results = append(results, result)
    }
    
    // 3. Return
    return results, nil
}

func (c *MyChain) process(ctx context.Context, doc runtime.Document) (string, error) {
    prompt := fmt.Sprintf("Analyze: %s", doc.Content)
    return c.llm.Complete(ctx, prompt)
}

var _ runtime.Chain = (*MyChain)(nil)
```

### Custom Chain Ideas

- **Translation Chain** - Translate documents to different languages
- **Extraction Chain** - Extract specific information from documents
- **Classification Chain** - Classify documents by category
- **Validation Chain** - Validate outputs of other chains
- **Formatting Chain** - Format outputs into specific formats (JSON, CSV, etc.)

---

## 📚 Examples

### Example 1: Customer Support Chatbot

```go
// Multi-turn with memory and web search
memory := inmemory.New()
agent := agent.New(llm, memory)
agent.RegisterTool(webSearchTool)

// Conversation
response1, _ := agent.Run(ctx, "I have a billing question")
response2, _ := agent.Run(ctx, "How do I pay my invoice?")
response3, _ := agent.Run(ctx, "What payment methods do you accept?")
// Agent remembers context, searches web for answers
```

### Example 2: Research Paper Analyzer

```go
// Load → Split → Summarize → Extract Key Points → Q&A

documents, _ := loader.LoadPDF("paper.pdf")

// Summarize
summarizer := summarization.New(llm, summarization.MapReduceStrategy)
summary, _ := summarizer.Run(ctx, documents)
fmt.Println("Summary:", summary)

// Extract key points (custom chain)
extractor := custom.NewExtractionChain(llm)
keyPoints, _ := extractor.Run(ctx, documents)
fmt.Println("Key Points:", keyPoints)

// Answer questions
vectorStore.Add(ctx, documents)
retriever := retrieval.New(vectorStore, 3)
qaChain := qa.New(retriever, llm, 3)

answer, _ := qaChain.Run(ctx, "What methodology did they use?")
fmt.Println("Answer:", answer)
```

### Example 3: Data Pipeline

```go
// Load CSV → Summarize rows → Generate SQL → Execute → Format

// Load data
data := loader.LoadCSV("sales.csv")

// Generate SQL
schema := "sales (id, product, amount, date)"
sqlChain := sql.New(llm, schema)
query, _ := sqlChain.Run(ctx, "Total revenue by product")
fmt.Println("Generated SQL:", query)

// Execute query and format results
results := executeSQL(query)
formatter := custom.NewFormatterChain(llm)
formatted, _ := formatter.Run(ctx, results)
fmt.Println("Formatted:", formatted)
```

---

## 🔐 Best Practices

### Error Handling

```go
// ✅ Good: Handle errors
result, err := chain.Run(ctx, input)
if err != nil {
    log.Errorf("chain failed: %v", err)
    return err
}

// ❌ Bad: Ignore errors
result, _ := chain.Run(ctx, input)
```

### Context Management

```go
// ✅ Good: Use context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
result, _ := chain.Run(ctx, input)

// ❌ Bad: No timeout
result, _ := chain.Run(context.Background(), input)
```

### Type Safety

```go
// ✅ Good: Assert types
if docs, ok := input.([]runtime.Document); ok {
    // Process docs
}

// ❌ Bad: Assume type
docs := input.([]runtime.Document)  // Panics if wrong type
```

---

## 📈 Performance Tips

### Batch Processing

```go
// Process multiple inputs efficiently
inputs := []string{"query1", "query2", "query3"}
results := make([]string, len(inputs))

for i, input := range inputs {
    result, _ := chain.Run(ctx, input)
    results[i] = result
}
```

### Parallel Execution

```go
// Use goroutines for parallel chains
wg := sync.WaitGroup{}
results := make([]string, len(inputs))

for i, input := range inputs {
    wg.Add(1)
    go func(idx int, inp string) {
        defer wg.Done()
        result, _ := chain.Run(ctx, inp)
        results[idx] = result
    }(i, input)
}

wg.Wait()
```

### Caching Results

```go
// Cache chain outputs
cache := make(map[string]string)
key := fmt.Sprintf("chain:%s", hashString(input))

if cached, ok := cache[key]; ok {
    return cached  // Use cached result
}

result, _ := chain.Run(ctx, input)
cache[key] = result
return result
```

---

## 🐛 Troubleshooting

### Chain Returns Unexpected Output

```go
// Check: Are you handling the output type correctly?
result, _ := chain.Run(ctx, input)

// Verify type
switch v := result.(type) {
case string:
    fmt.Println("String:", v)
case []string:
    fmt.Println("Strings:", v)
default:
    fmt.Printf("Unknown type: %T\n", v)
}
```

### Chain Too Slow

```go
// Check: Are you using the right strategy?

// Fast (Stuff)
summarizer := summarization.New(llm, summarization.StuffStrategy)

// Slow (MapReduce)
summarizer := summarization.New(llm, summarization.MapReduceStrategy)

// Use Stuff for small documents, MapReduce for large
```

### Agent Not Using Tools

```go
// Check: Tools are registered?
agent.RegisterTool(myTool)

// Check: Tool description is clear?
tool := &Tool{
    Name: "my_tool",
    Description: "What this tool does",  // Must be descriptive
    Run: func(...) {...},
}
```

---

## 📖 More Info

- [Architecture Overview](../docs/ARCHITECTURE.md)
- [Extending Guide](../docs/EXTENDING.md)
- [Main README](../README.md)
- [Contributing Guide](../CONTRIBUTING.md)

---

Questions? [Start a discussion](https://github.com/rajveer43/goagentflow/discussions)!
