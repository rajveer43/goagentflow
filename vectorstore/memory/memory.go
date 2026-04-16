package memory

import (
	"context"
	"fmt"
	"math"
	"sync"

	"github.com/rajveer43/goagentflow/runtime"
	"github.com/rajveer43/goagentflow/types"
)

// entry represents a stored document with its embedding.
type entry struct {
	id        string
	doc       types.Document
	embedding []float32
}

// Store is an in-memory vector store using cosine similarity.
// Pattern: Repository - thread-safe document storage by semantic similarity
// Uses pure Go with no external dependencies. Suitable for small-to-medium datasets.
type Store struct {
	mu       sync.RWMutex
	entries  []entry
	idToIdx  map[string]int // maps document ID to index in entries slice
	nextID   int64
}

// New creates a new in-memory vector store.
func New() *Store {
	return &Store{
		entries: make([]entry, 0),
		idToIdx: make(map[string]int),
	}
}

// Add indexes documents with their embeddings.
// Returns the IDs assigned to each document, in the same order as input.
func (s *Store) Add(ctx context.Context, docs []types.Document, embeddings [][]float32) ([]string, error) {
	if len(docs) != len(embeddings) {
		return nil, fmt.Errorf("docs and embeddings length mismatch: %d vs %d", len(docs), len(embeddings))
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	ids := make([]string, len(docs))
	for i := range docs {
		id := fmt.Sprintf("doc_%d", s.nextID)
		s.nextID++

		idx := len(s.entries)
		s.entries = append(s.entries, entry{
			id:        id,
			doc:       docs[i],
			embedding: embeddings[i],
		})
		s.idToIdx[id] = idx
		ids[i] = id
	}

	return ids, nil
}

// Search finds the top-k documents most similar to the query embedding.
// Returns results sorted by similarity (highest first).
func (s *Store) Search(ctx context.Context, embedding []float32, k int) ([]runtime.SearchResult, error) {
	if k <= 0 {
		return nil, fmt.Errorf("k must be positive, got %d", k)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.entries) == 0 {
		return []runtime.SearchResult{}, nil
	}

	// Compute cosine similarity with all entries
	similarities := make([]struct {
		idx   int
		score float32
	}, len(s.entries))

	for i, e := range s.entries {
		sim := cosineSimilarity(embedding, e.embedding)
		similarities[i] = struct {
			idx   int
			score float32
		}{idx: i, score: sim}
	}

	// Sort by similarity (descending) using bubble sort for simplicity
	// In production, consider using sort.Slice with quicksort
	for i := 0; i < len(similarities); i++ {
		for j := i + 1; j < len(similarities); j++ {
			if similarities[j].score > similarities[i].score {
				similarities[i], similarities[j] = similarities[j], similarities[i]
			}
		}
	}

	// Return top-k
	limit := k
	if limit > len(similarities) {
		limit = len(similarities)
	}

	results := make([]runtime.SearchResult, limit)
	for i := 0; i < limit; i++ {
		idx := similarities[i].idx
		results[i] = runtime.SearchResult{
			ID:       s.entries[idx].id,
			Document: s.entries[idx].doc,
			Score:    similarities[i].score,
		}
	}

	return results, nil
}

// Delete removes documents by their IDs.
func (s *Store) Delete(ctx context.Context, ids []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Mark indices for deletion
	toDelete := make(map[int]bool)
	for _, id := range ids {
		if idx, ok := s.idToIdx[id]; ok {
			toDelete[idx] = true
			delete(s.idToIdx, id)
		}
	}

	if len(toDelete) == 0 {
		return nil
	}

	// Rebuild entries slice, skipping deleted indices
	newEntries := make([]entry, 0, len(s.entries)-len(toDelete))
	newIDToIdx := make(map[string]int)
	for i, e := range s.entries {
		if !toDelete[i] {
			newIDToIdx[e.id] = len(newEntries)
			newEntries = append(newEntries, e)
		}
	}

	s.entries = newEntries
	s.idToIdx = newIDToIdx

	return nil
}

// Clear removes all documents from the store.
func (s *Store) Clear(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.entries = make([]entry, 0)
	s.idToIdx = make(map[string]int)
	s.nextID = 0

	return nil
}

// Size returns the number of documents currently stored.
func (s *Store) Size(ctx context.Context) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.entries), nil
}

// cosineSimilarity computes the cosine similarity between two vectors.
// Assumes vectors are already normalized (for efficiency).
// Returns a value between -1 and 1, where 1 is identical.
func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	if len(a) == 0 {
		return 0
	}

	// Compute dot product
	dotProduct := float32(0)
	magnitudeA := float32(0)
	magnitudeB := float32(0)

	for i := range a {
		dotProduct += a[i] * b[i]
		magnitudeA += a[i] * a[i]
		magnitudeB += b[i] * b[i]
	}

	// Compute magnitudes
	magnitudeA = float32(math.Sqrt(float64(magnitudeA)))
	magnitudeB = float32(math.Sqrt(float64(magnitudeB)))

	if magnitudeA == 0 || magnitudeB == 0 {
		return 0
	}

	return dotProduct / (magnitudeA * magnitudeB)
}
