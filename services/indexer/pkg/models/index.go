package models

// RawPage represents the document structure that your Web Spider saved.
type RawPage struct {
	URL         string `bson:"url"`
	Title       string `bson:"title"`
	TextContent string `bson:"text_content"`
	Indexed     bool   `bson:"indexed,omitempty"` // Flag so we don't index the same page twice
}

// DocumentMatch represents a single web page that contains a specific search term.
type DocumentMatch struct {
	URL   string  `bson:"url" json:"url"`
	Title string  `bson:"title" json:"title"`
	TF    float64 `bson:"tf" json:"tf"` // Term Frequency: How relevant this page is to the term
}

// IndexEntry represents a single word in the inverted index.
// It holds the word (Term) and an array of all the pages (Matches) that contain it.
type IndexEntry struct {
	Term    string          `bson:"term" json:"term"`
	Matches []DocumentMatch `bson:"matches" json:"matches"`
}
