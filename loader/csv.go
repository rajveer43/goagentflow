package loader

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// CSVLoader loads CSV files with customizable column selection.
// Pattern: Strategy
// DSA: streaming row iteration with encoding/csv (O(n) memory usage)
type CSVLoader struct {
	Path    string   // file path
	Columns []string // if empty, use all columns
}

// NewCSVLoader creates a new CSV file loader.
func NewCSVLoader(path string) *CSVLoader {
	return &CSVLoader{Path: path}
}

// WithColumns specifies which columns to include in documents.
func (c *CSVLoader) WithColumns(cols ...string) *CSVLoader {
	c.Columns = cols
	return c
}

// Load reads CSV file and returns each row as a Document.
func (c *CSVLoader) Load(ctx context.Context) ([]Document, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Open file
	file, err := os.Open(c.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open csv file %s: %w", c.Path, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv header: %w", err)
	}

	// Determine columns to include
	columns := c.Columns
	if len(columns) == 0 {
		columns = header
	}

	// Find column indices
	colIndices := make(map[string]int)
	for i, col := range header {
		colIndices[col] = i
	}

	// Read rows
	var docs []Document
	rowNum := 1
	for {
		record, err := reader.Read()
		if err != nil {
			break // EOF or error
		}

		// Build document content from selected columns
		var lines []string
		for _, col := range columns {
			idx, ok := colIndices[col]
			if !ok {
				continue // column not found, skip
			}
			if idx < len(record) {
				lines = append(lines, fmt.Sprintf("%s: %s", col, record[idx]))
			}
		}

		content := strings.Join(lines, "\n")
		if content == "" {
			continue // skip empty rows
		}

		doc := Document{
			PageContent: content,
			Metadata: map[string]any{
				"source":   c.Path,
				"row_index": rowNum,
				"content_type": "text/csv",
			},
		}
		docs = append(docs, doc)
		rowNum++
	}

	return docs, nil
}
