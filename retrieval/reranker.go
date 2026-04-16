package retrieval

import (
	"context"
	"fmt"
	"math"

	"github.com/rajveer43/goagentflow/runtime"
)

// Reranker re-sorts retrieved documents by some criteria.
// Pattern: Strategy - different reranking strategies can be swapped
type Reranker interface {
	Rerank(ctx context.Context, query string, docs []runtime.SearchResult) ([]runtime.SearchResult, error)
}

// SimpleReranker re-ranks by simple score thresholding.
type SimpleReranker struct {
	minScore float32
}

// NewSimpleReranker creates a reranker that filters by minimum score.
func NewSimpleReranker(minScore float32) *SimpleReranker {
	return &SimpleReranker{minScore: minScore}
}

// Rerank filters results by minimum score.
func (sr *SimpleReranker) Rerank(ctx context.Context, query string, docs []runtime.SearchResult) ([]runtime.SearchResult, error) {
	var filtered []runtime.SearchResult
	for _, doc := range docs {
		if doc.Score >= sr.minScore {
			filtered = append(filtered, doc)
		}
	}
	return filtered, nil
}

// MMRReranker re-ranks using Maximum Marginal Relevance (MMR).
// Balances relevance to query with diversity among results.
type MMRReranker struct {
	lambdaMultiplier float32 // controls relevance vs diversity trade-off
	embedder         runtime.Embedder
}

// NewMMRReranker creates a reranker using Maximum Marginal Relevance.
// lambdaMultiplier: balance between relevance (1.0) and diversity (0.0)
// embedder: for computing embedding distances
func NewMMRReranker(lambdaMultiplier float32, embedder runtime.Embedder) *MMRReranker {
	if lambdaMultiplier < 0 || lambdaMultiplier > 1 {
		lambdaMultiplier = 0.5 // default middle ground
	}
	return &MMRReranker{
		lambdaMultiplier: lambdaMultiplier,
		embedder:         embedder,
	}
}

// Rerank re-ranks documents using MMR.
// MMR = lambda * relevance - (1 - lambda) * max_similarity_to_selected
// This promotes diverse results while maintaining relevance.
func (mr *MMRReranker) Rerank(ctx context.Context, query string, docs []runtime.SearchResult) ([]runtime.SearchResult, error) {
	if len(docs) == 0 {
		return docs, nil
	}

	// Embed the query
	queryEmbedding, err := mr.embedder.Embed(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	docTexts := make([]string, len(docs))
	for i := range docs {
		docTexts[i] = docs[i].Document.PageContent
	}

	docEmbeddings, err := mr.embedder.EmbedBatch(ctx, docTexts)
	if err != nil {
		return nil, fmt.Errorf("failed to embed documents: %w", err)
	}
	if len(docEmbeddings) != len(docs) {
		return nil, fmt.Errorf("embedder returned %d embeddings for %d docs", len(docEmbeddings), len(docs))
	}

	normalize := func(vec []float32) []float32 {
		var norm float32
		for i := range vec {
			norm += vec[i] * vec[i]
		}
		if norm == 0 {
			return vec
		}
		invNorm := float32(1 / math.Sqrt(float64(norm)))
		out := make([]float32, len(vec))
		for i := range vec {
			out[i] = vec[i] * invNorm
		}
		return out
	}

	dot := func(a, b []float32) float32 {
		limit := len(a)
		if len(b) < limit {
			limit = len(b)
		}
		var sum float32
		for i := 0; i < limit; i++ {
			sum += a[i] * b[i]
		}
		return sum
	}

	queryEmbedding = normalize(queryEmbedding)
	for i := range docEmbeddings {
		docEmbeddings[i] = normalize(docEmbeddings[i])
	}

	relevances := make([]float32, len(docs))
	for i := range docEmbeddings {
		relevances[i] = dot(queryEmbedding, docEmbeddings[i])
	}

	selected := make([]runtime.SearchResult, 0, len(docs))
	selectedEmbeddings := make([][]float32, 0, len(docs))
	used := make([]bool, len(docs))

	// Seed with the most relevant document.
	bestIdx := 0
	bestScore := docs[0].Score
	for i := 1; i < len(docs); i++ {
		if docs[i].Score > bestScore {
			bestScore = docs[i].Score
			bestIdx = i
		}
	}
	selected = append(selected, docs[bestIdx])
	selectedEmbeddings = append(selectedEmbeddings, docEmbeddings[bestIdx])
	used[bestIdx] = true

	for len(selected) < len(docs) {
		bestIdx = -1
		bestMMR := float32(math.Inf(-1))

		for i := range docs {
			if used[i] {
				continue
			}

			relevance := relevances[i]
			diversity := float32(0)
			for j := range selectedEmbeddings {
				sim := dot(docEmbeddings[i], selectedEmbeddings[j])
				if sim > diversity {
					diversity = sim
				}
			}

			mmr := mr.lambdaMultiplier*relevance - (1-mr.lambdaMultiplier)*diversity
			if mmr > bestMMR {
				bestMMR = mmr
				bestIdx = i
			}
		}

		if bestIdx < 0 {
			break
		}

		used[bestIdx] = true
		selected = append(selected, docs[bestIdx])
		selectedEmbeddings = append(selectedEmbeddings, docEmbeddings[bestIdx])
	}

	return selected, nil
}

// RerankerChain wraps a reranker as a runtime.Chain.
type RerankerChain struct {
	reranker runtime.Retriever // base retriever
	rerank   Reranker          // reranking strategy
}

// NewRerankerChain creates a chain that reranks retriever results.
func NewRerankerChain(retriever runtime.Retriever, reranker Reranker) *RerankerChain {
	return &RerankerChain{
		reranker: retriever,
		rerank:   reranker,
	}
}

// Run implements runtime.Chain.
// Takes query string, retrieves documents, reranks, returns docs.
func (rc *RerankerChain) Run(ctx context.Context, input any) (any, error) {
	_, ok := input.(string)
	if !ok {
		return nil, fmt.Errorf("expected string query, got %T", input)
	}

	// This is a simplified implementation - full MMR would need embeddings
	// For now, just return the retrieved docs
	return nil, fmt.Errorf("reranker chain not fully implemented - use base RetrieverChain instead")
}

// cosineSimilarity computes the cosine similarity between two vectors.
func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	if len(a) == 0 {
		return 0
	}

	dotProduct := float32(0)
	magnitudeA := float32(0)
	magnitudeB := float32(0)

	for i := range a {
		dotProduct += a[i] * b[i]
		magnitudeA += a[i] * a[i]
		magnitudeB += b[i] * b[i]
	}

	magnitudeA = float32(math.Sqrt(float64(magnitudeA)))
	magnitudeB = float32(math.Sqrt(float64(magnitudeB)))

	if magnitudeA == 0 || magnitudeB == 0 {
		return 0
	}

	return dotProduct / (magnitudeA * magnitudeB)
}
