package parser

import (
	"strings"

	"golang.org/x/net/html"
	"moogle-go/services/spider/pkg/models"
)

func ExtractLinks(node *html.Node, baseURL string, normalizer *URLNormalizer) []models.Link {
	var links []models.Link
	var crawler func(*html.Node)

	crawler = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href := strings.TrimSpace(attr.Val)
					if href == "" || strings.HasPrefix(href, "#") || strings.HasPrefix(href, "javascript:") || strings.HasPrefix(href, "mailto:") {
						continue
					}

					resolvedURL, err := normalizer.Resolve(baseURL, href)
					if err == nil && resolvedURL != "" {
						anchorText := extractNodeText(n)
						links = append(links, models.Link{
							URL:  resolvedURL,
							Text: strings.TrimSpace(anchorText),
						})
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawler(c)
		}
	}

	crawler(node)
	return links
}

func ExtractImages(node *html.Node, baseURL string, normalizer *URLNormalizer) []models.Image {
	var images []models.Image
	var crawler func(*html.Node)

	crawler = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			var src, alt string
			for _, attr := range n.Attr {
				if attr.Key == "src" {
					src = strings.TrimSpace(attr.Val)
				}
				if attr.Key == "alt" {
					alt = strings.TrimSpace(attr.Val)
				}
			}

			if src != "" && !strings.HasPrefix(src, "data:") {
				resolvedURL, err := normalizer.Resolve(baseURL, src)
				if err == nil && resolvedURL != "" {
					images = append(images, models.Image{
						URL: resolvedURL,
						Alt: alt,
					})
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawler(c)
		}
	}

	crawler(node)
	return images
}

func extractNodeText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var sb strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		sb.WriteString(extractNodeText(c))
		sb.WriteString(" ")
	}
	return sb.String()
}