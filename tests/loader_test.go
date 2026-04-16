package tests

import (
	"context"
	"os"
	"testing"

	"github.com/rajveer43/goagentflow/loader"
	"github.com/rajveer43/goagentflow/memory/inmemory"
	"github.com/rajveer43/goagentflow/runtime"
)

// TestTextLoader tests loading plain text files
func TestTextLoader(t *testing.T) {
	// Create temporary test file
	tmpfile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := "Hello, world!\nThis is a test document."
	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Load document
	ldr := loader.NewTextLoader(tmpfile.Name())
	docs, err := ldr.Load(context.Background())
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(docs) != 1 {
		t.Fatalf("Expected 1 document, got %d", len(docs))
	}

	if docs[0].PageContent != content {
		t.Fatalf("Content mismatch: got %q", docs[0].PageContent)
	}

	if docs[0].Metadata["source"] != tmpfile.Name() {
		t.Fatalf("Metadata mismatch: source not set")
	}
}

// TestCSVLoader tests loading CSV files
func TestCSVLoader(t *testing.T) {
	// Create temporary CSV file
	tmpfile, err := os.CreateTemp("", "test*.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := "name,age\nAlice,30\nBob,25\n"
	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Load CSV
	ldr := loader.NewCSVLoader(tmpfile.Name())
	docs, err := ldr.Load(context.Background())
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(docs) != 2 {
		t.Fatalf("Expected 2 documents, got %d", len(docs))
	}

	if !contains(docs[0].PageContent, "Alice") {
		t.Fatalf("Content mismatch: got %q", docs[0].PageContent)
	}
}

// TestHTMLLoader tests loading and parsing HTML files
func TestHTMLLoader(t *testing.T) {
	// Create temporary HTML file
	tmpfile, err := os.CreateTemp("", "test*.html")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	html := `<html>
	<head><title>Test Page</title></head>
	<body><p>Hello from HTML</p></body>
	</html>`
	if _, err := tmpfile.WriteString(html); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Load HTML
	ldr := loader.NewHTMLLoader(tmpfile.Name())
	docs, err := ldr.Load(context.Background())
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(docs) != 1 {
		t.Fatalf("Expected 1 document, got %d", len(docs))
	}

	if !contains(docs[0].PageContent, "Hello from HTML") {
		t.Fatalf("Content mismatch: got %q", docs[0].PageContent)
	}

	if docs[0].Metadata["title"] != "Test Page" {
		t.Fatalf("Title mismatch: got %q", docs[0].Metadata["title"])
	}
}

// TestLoaderChain tests that loaders work as runtime.Chain
func TestLoaderChain(t *testing.T) {
	// Create temporary test file
	tmpfile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString("test content"); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Run as chain
	ldr := loader.NewTextLoader(tmpfile.Name())
	chain := loader.NewLoaderChain(ldr)

	result, err := chain.Run(context.Background(), nil)
	if err != nil {
		t.Fatalf("Chain run failed: %v", err)
	}

	docs, ok := result.([]loader.Document)
	if !ok {
		t.Fatalf("Expected []Document, got %T", result)
	}

	if len(docs) != 1 {
		t.Fatalf("Expected 1 document, got %d", len(docs))
	}
}

// TestInjectIntoMemory tests injecting documents into memory
func TestInjectIntoMemory(t *testing.T) {
	mem := inmemory.New()

	docs := []loader.Document{
		{PageContent: "Doc 1", Metadata: map[string]any{"id": 1}},
		{PageContent: "Doc 2", Metadata: map[string]any{"id": 2}},
	}

	// Inject documents
	err := loader.InjectIntoMemory(context.Background(), mem, "docs", docs, true)
	if err != nil {
		t.Fatalf("InjectIntoMemory failed: %v", err)
	}

	// Verify documents are stored
	stored, err := mem.Get(context.Background(), "docs")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	storedDocs, ok := stored.([]loader.Document)
	if !ok {
		t.Fatalf("Expected []Document, got %T", stored)
	}

	if len(storedDocs) != 2 {
		t.Fatalf("Expected 2 documents, got %d", len(storedDocs))
	}

	// Verify messages were added
	messages, err := mem.GetMessages(context.Background())
	if err != nil {
		t.Fatalf("GetMessages failed: %v", err)
	}

	if len(messages) < 2 {
		t.Fatalf("Expected at least 2 messages, got %d", len(messages))
	}
}

// TestInjectIntoState tests injecting documents into state
func TestInjectIntoState(t *testing.T) {
	state := runtime.NewState("input")

	docs := []loader.Document{
		{PageContent: "Doc 1", Metadata: map[string]any{"id": 1}},
	}

	// Inject documents
	loader.InjectIntoState(state, "docs", docs)

	// Verify documents are stored
	stored, ok := state.Get("docs")
	if !ok {
		t.Fatal("Documents not found in state")
	}

	storedDocs, ok := stored.([]loader.Document)
	if !ok {
		t.Fatalf("Expected []Document, got %T", stored)
	}

	if len(storedDocs) != 1 {
		t.Fatalf("Expected 1 document, got %d", len(storedDocs))
	}
}

// TestMultiLoader tests combining multiple loaders
func TestMultiLoader(t *testing.T) {
	// Create two temporary files
	tmpfile1, err := os.CreateTemp("", "test1*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile1.Name())

	tmpfile2, err := os.CreateTemp("", "test2*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile2.Name())

	if _, err := tmpfile1.WriteString("Content 1"); err != nil {
		t.Fatal(err)
	}
	if _, err := tmpfile2.WriteString("Content 2"); err != nil {
		t.Fatal(err)
	}
	tmpfile1.Close()
	tmpfile2.Close()

	// Create multi-loader
	ldr1 := loader.NewTextLoader(tmpfile1.Name())
	ldr2 := loader.NewTextLoader(tmpfile2.Name())
	multi := loader.NewMultiLoader(ldr1, ldr2)

	docs, err := multi.Load(context.Background())
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(docs) != 2 {
		t.Fatalf("Expected 2 documents, got %d", len(docs))
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
