package parser

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
	"moogle-go/services/spider/pkg/models"
)

type HTMLParser struct {
	normalizer *URLNormalizer
}

func NewHTMLParser(normalizer *URLNormalizer) *HTMLParser {
	return &HTMLParser{
		normalizer: normalizer,
	}
}

func (hp *HTMLParser) Parse(rawURL string, statusCode int, body []byte) (*models.Page, error) {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse html content: %w", err)
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid page url: %w", err)
	}

	title := extractTitle(doc)
	description := extractMetaDescription(doc)
	textContent := extractCleanBodyText(doc)
	outlinks := ExtractLinks(doc, rawURL, hp.normalizer)
	images := ExtractImages(doc, rawURL, hp.normalizer)

	page := &models.Page{
		URL:         rawURL,
		Domain:      parsedURL.Hostname(),
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		RawHTML:     string(body),
		TextContent: strings.TrimSpace(textContent),
		Outlinks:    outlinks,
		Images:      images,
		StatusCode:  statusCode,
		CrawledAt:   time.Now().UTC(),
	}

	return page, nil
}

func extractTitle(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "title" {
		if n.FirstChild != nil {
			return n.FirstChild.Data
		}
		return ""
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if title := extractTitle(c); title != "" {
			return title
		}
	}
	return ""
}

func extractMetaDescription(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "meta" {
		var isDesc bool
		var content string
		for _, attr := range n.Attr {
			if strings.EqualFold(attr.Key, "name") && strings.EqualFold(attr.Val, "description") {
				isDesc = true
			}
			if strings.EqualFold(attr.Key, "content") {
				content = attr.Val
			}
		}
		if isDesc {
			return content
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if desc := extractMetaDescription(c); desc != "" {
			return desc
		}
	}
	return ""
}

func extractCleanBodyText(n *html.Node) string {
	var builder strings.Builder

	var f func(*html.Node)
	f = func(node *html.Node) {
		if node.Type == html.ElementNode {
			if node.Data == "script" || node.Data == "style" || node.Data == "noscript" || node.Data == "head" {
				return
			}
		}

		if node.Type == html.TextNode {
			text := strings.TrimSpace(node.Data)
			if text != "" {
				builder.WriteString(text)
				builder.WriteString(" ")
			}
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(n)
	return builder.String()
}