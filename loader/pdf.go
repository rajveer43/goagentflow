package loader

import (
	"context"
	"fmt"
	"os"

	"github.com/ledongthuc/pdf"
)

// PDFLoader loads text from PDF files.
// Pattern: Strategy
// Each page becomes one Document with page metadata.
type PDFLoader struct {
	Path string
}

// NewPDFLoader creates a new PDF file loader.
func NewPDFLoader(path string) *PDFLoader {
	return &PDFLoader{Path: path}
}

// Load reads PDF file and returns text content per page.
func (p *PDFLoader) Load(ctx context.Context) ([]Document, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Open PDF file
	file, err := os.Open(p.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open pdf file %s: %w", p.Path, err)
	}
	defer file.Close()

	// Get file size for metadata
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat pdf file: %w", err)
	}

	// Read PDF
	reader, err := pdf.NewReader(file, stat.Size())
	if err != nil {
		return nil, fmt.Errorf("failed to read pdf: %w", err)
	}

	totalPages := reader.NumPage()
	var docs []Document

	// Extract text from each page
	for pageNum := 1; pageNum <= totalPages; pageNum++ {
		page := reader.Page(pageNum)

		// Extract text from page
		text, err := page.GetPlainText(nil)
		if err != nil {
			// Continue on error, some pages may be unreadable
			text = fmt.Sprintf("[Error reading page %d: %v]", pageNum, err)
		}

		if text == "" {
			text = "[Empty page]"
		}

		doc := Document{
			PageContent: text,
			Metadata: map[string]any{
				"source":       p.Path,
				"page":         pageNum,
				"total_pages":  totalPages,
				"content_type": "application/pdf",
			},
		}
		docs = append(docs, doc)
	}

	return docs, nil
}
