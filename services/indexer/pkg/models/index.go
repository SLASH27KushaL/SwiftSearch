package models

type RawPage struct {
	URL         string `bson:"url"`
	Title       string `bson:"title"`
	TextContent string `bson:"text_content"`
	Indexed     bool   `bson:"indexed,omitempty"`
}

type DocumentMatch struct {
	URL   string  `bson:"url" json:"url"`
	Title string  `bson:"title" json:"title"`
	TF    float64 `bson:"tf" json:"tf"`
}

type IndexEntry struct {
	Term    string          `bson:"term" json:"term"`
	Matches []DocumentMatch `bson:"matches" json:"matches"`
}
