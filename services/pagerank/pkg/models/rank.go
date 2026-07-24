package models

// GraphNode represents a web page and all the URLs it links to.
// We assume your spider saved an "outlinks" array. If it didn't yet, you can modify the spider later!
type GraphNode struct {
	URL      string   `bson:"url"`
	Outlinks []string `bson:"outlinks"`
}

// PageRankScore holds the final computed mathematical importance of a URL.
type PageRankScore struct {
	URL   string  `bson:"url"`
	Score float64 `bson:"score"`
}