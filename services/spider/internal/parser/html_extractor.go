package parser

import (
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type PageData struct {
	URL   string
	HTML  string
	Links []string
}

func FetchAndParse(targetURL string) (*PageData, error) {
	// 1. Create a custom HTTP client with a 10-second timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 2. Formulate the GET request
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, err
	}

	// 3. Inject browser-like headers to bypass bot-blockers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 SwiftSearchBot/1.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	// 4. Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 5. Read the fully downloaded HTML
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 🚨 THE UTF-8 FIX: Scrub invalid characters before saving to MongoDB
	rawHTML := string(bodyBytes)
	cleanHTML := strings.ToValidUTF8(rawHTML, " ")

	links := extractLinks(cleanHTML, targetURL)

	return &PageData{
		URL:   targetURL,
		HTML:  cleanHTML, // 👈 SAVING THE CLEANED HTML
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