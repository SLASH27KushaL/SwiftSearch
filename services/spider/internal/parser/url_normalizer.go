package parser

import (
	"fmt"
	"net/url"
	"strings"
)

// URLNormalizer handles URL cleaning, relative-to-absolute resolution, and formatting.
type URLNormalizer struct{}

// NewURLNormalizer creates a new URLNormalizer instance.
func NewURLNormalizer() *URLNormalizer {
	return &URLNormalizer{}
}

// Normalize cleans up a URL by lowercasing scheme/host, stripping fragments (#), and removing default ports.
func (n *URLNormalizer) Normalize(rawURL string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return "", fmt.Errorf("failed to parse url: %w", err)
	}

	// Validate supported scheme
	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return "", fmt.Errorf("unsupported scheme: %s", scheme)
	}

	parsed.Scheme = scheme
	parsed.Host = strings.ToLower(parsed.Host)
	parsed.Fragment = "" // Strip fragments (#heading)

	// Remove default HTTP/HTTPS ports if explicit in host string
	if scheme == "http" && strings.HasSuffix(parsed.Host, ":80") {
		parsed.Host = strings.TrimSuffix(parsed.Host, ":80")
	} else if scheme == "https" && strings.HasSuffix(parsed.Host, ":443") {
		parsed.Host = strings.TrimSuffix(parsed.Host, ":443")
	}

	// Normalize empty or trailing slashes in path
	if parsed.Path == "" {
		parsed.Path = "/"
	}

	return parsed.String(), nil
}

// Resolve converts a relative URL (e.g. "/about") into an absolute URL using the base URL.
func (n *URLNormalizer) Resolve(baseURLStr, refURLStr string) (string, error) {
	baseURL, err := url.Parse(baseURLStr)
	if err != nil {
		return "", fmt.Errorf("invalid base url: %w", err)
	}

	refURL, err := url.Parse(strings.TrimSpace(refURLStr))
	if err != nil {
		return "", fmt.Errorf("invalid reference url: %w", err)
	}

	resolved := baseURL.ResolveReference(refURL)
	return n.Normalize(resolved.String())
}