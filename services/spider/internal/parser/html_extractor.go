package parser

import (
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type PageData struct {
	URL   string
	HTML  string
	Links []string
}

func FetchAndParse(targetURL string) (*PageData, error) {
	resp, err := http.Get(targetURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	htmlContent := string(bodyBytes)
	links := extractLinks(htmlContent, targetURL)

	return &PageData{
		URL:   targetURL,
		HTML:  htmlContent,
		Links: links,
	}, nil
}

func extractLinks(htmlBody, base string) []string {
	var links []string
	tokenizer := html.NewTokenizer(strings.NewReader(htmlBody))

	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			break
		}
		if tt == html.StartTagToken {
			t := tokenizer.Token()
			if t.Data == "a" {
				for _, attr := range t.Attr {
					if attr.Key == "href" && strings.HasPrefix(attr.Val, "http") {
						links = append(links, attr.Val)
					}
				}
			}
		}
	}
	return links
}
