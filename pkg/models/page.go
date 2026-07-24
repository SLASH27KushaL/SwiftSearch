package models

import "time"

type Link struct {
    URL  string `json:"url" bson:"url"`
    Text string `json:"text,omitempty" bson:"text,omitempty"`
}

type Image struct {
    URL string `json:"url" bson:"url"`
    Alt string `json:"alt,omitempty" bson:"alt,omitempty"`
}

type Page struct {
    ID          string    `json:"id,omitempty" bson:"_id,omitempty"`
    URL         string    `json:"url" bson:"url"`
    Domain      string    `json:"domain" bson:"domain"`
    Title       string    `json:"title" bson:"title"`
    Description string    `json:"description" bson:"description"`
    RawHTML     string    `json:"raw_html,omitempty" bson:"raw_html,omitempty"`
    TextContent string    `json:"text_content" bson:"text_content"`
    Outlinks    []Link    `json:"outlinks" bson:"outlinks"`
    Images      []Image   `json:"images" bson:"images"`
    StatusCode  int       `json:"status_code" bson:"status_code"`
    CrawledAt   time.Time `json:"crawled_at" bson:"crawled_at"`
}

type CrawlJob struct {
    URL   string `json:"url"`
    Depth int    `json:"depth"`
}
