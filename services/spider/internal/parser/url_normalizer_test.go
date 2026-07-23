package parser

import (
	"testing"
)

func TestNormalize(t *testing.T) {
	normalizer := NewURLNormalizer()

	tests := []struct {
		name     string
		input    string
		expected string
		hasErr   bool
	}{
		{
			name:     "Strip fragment and lowercase host",
			input:    "https://EXAMPLE.com/path#section",
			expected: "https://example.com/path",
			hasErr:   false,
		},
		{
			name:     "Remove default HTTP port 80",
			input:    "HTTP://example.com:80/about",
			expected: "http://example.com/about",
			hasErr:   false,
		},
		{
			name:     "Remove default HTTPS port 443",
			input:    "https://example.com:443/docs/",
			expected: "https://example.com/docs/",
			hasErr:   false,
		},
		{
			name:     "Add root slash if path is empty",
			input:    "https://example.com",
			expected: "https://example.com/",
			hasErr:   false,
		},
		{
			name:     "Reject unsupported scheme (javascript)",
			input:    "javascript:void(0)",
			expected: "",
			hasErr:   true,
		},
		{
			name:     "Reject unsupported scheme (mailto)",
			input:    "mailto:test@example.com",
			expected: "",
			hasErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizer.Normalize(tt.input)
			if (err != nil) != tt.hasErr {
				t.Fatalf("Normalize(%q) error = %v, wantErr %v", tt.input, err, tt.hasErr)
			}
			if got != tt.expected {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestResolve(t *testing.T) {
	normalizer := NewURLNormalizer()
	baseURL := "https://example.com/blog/posts"

	tests := []struct {
		name     string
		ref      string
		expected string
	}{
		{
			name:     "Relative path with root slash",
			ref:      "/about",
			expected: "https://example.com/about",
		},
		{
			name:     "Relative path moving up a directory",
			ref:      "../contact",
			expected: "https://example.com/contact",
		},
		{
			name:     "Absolute external URL",
			ref:      "https://otherdomain.com/page",
			expected: "https://otherdomain.com/page",
		},
		{
			name:     "Query parameters relative reference",
			ref:      "?page=2",
			expected: "https://example.com/blog/posts?page=2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizer.Resolve(baseURL, tt.ref)
			if err != nil {
				t.Fatalf("Resolve(%q, %q) unexpected error: %v", baseURL, tt.ref, err)
			}
			if got != tt.expected {
				t.Errorf("Resolve(%q, %q) = %q, want %q", baseURL, tt.ref, got, tt.expected)
			}
		})
	}
}