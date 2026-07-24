package models

// DocumentMatch represents the nested array inside the Inverted Index
type DocumentMatch struct {
	URL   string  `bson:"url"`
	Title string  `bson:"title"`
	TF    float64 `bson:"tf"`
}

// IndexEntry represents a full document from the inverted_index collection
type IndexEntry struct {
	Term    string          `bson:"term"`
	Matches []DocumentMatch `bson:"matches"`
}

// PageRankScore represents a document from the pagerank_scores collection
type PageRankScore struct {
	URL   string  `bson:"url"`
	Score float64 `bson:"score"`
}

// SearchResult is a single ranked result returned to the frontend
type SearchResult struct {
	Title string  `json:"title"`
	URL   string  `json:"url"`
	Score float64 `json:"score"`
}

// SearchResponse is the final JSON payload
type SearchResponse struct {
	Query           string         `json:"query"`
	TotalResults    int            `json:"total_results"`
	ExecutionTimeMs int64          `json:"execution_time_ms"`
	Results         []SearchResult `json:"results"`
}