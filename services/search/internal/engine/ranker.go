package engine

import (
	"context"
	"sort"

	"moogle-go/services/search/internal/store"
	"moogle-go/services/search/pkg/models"
)

type Ranker struct {
	reader *store.MongoReader
	alpha  float64
	beta   float64
}

func NewRanker(reader *store.MongoReader, alpha, beta float64) *Ranker {
	return &Ranker{
		reader: reader,
		alpha:  alpha,
		beta:   beta,
	}
}

func (r *Ranker) ExecuteSearch(ctx context.Context, query string) ([]models.SearchResult, error) {
	tokens := CleanAndTokenizeQuery(query)
	if len(tokens) == 0 {
		return []models.SearchResult{}, nil
	}
log.Printf("Query: %s | Tokens: %v", query, tokens)
	// 1. Get Relevance (TF) and Titles
	tfScores, titles, err := r.reader.FetchTFIDF(ctx, tokens)
	if err != nil {
		return nil, err
	}

	// Extract unique URLs to query PageRank
	var urls []string
	for url := range tfScores {
		urls = append(urls, url)
	}

	// 2. Get Authority (PageRank)
	prScores, err := r.reader.FetchPageRanks(ctx, urls)
	if err != nil {
		return nil, err
	}

	// 3. Mathematical Fusion: Score = (Alpha * TF) + (Beta * PR)
	var results []models.SearchResult
	for url, tfScore := range tfScores {
		prScore := prScores[url] // Defaults to 0 if not found

		finalScore := (r.alpha * tfScore) + (r.beta * prScore)

		results = append(results, models.SearchResult{
			URL:   url,
			Title: titles[url],
			Score: finalScore,
		})
	}

	// 4. Sort results descending by Final Score
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results, nil
}