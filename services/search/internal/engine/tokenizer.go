package engine

import (
	"regexp"
	"strings"
)

var (
	nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
	stopWords            = map[string]struct{}{
		"a": {}, "and": {}, "the": {}, "is": {}, "in": {}, "at": {}, "of": {}, "on": {}, "for": {}, "to": {}, "with": {},
	}
)

func CleanAndTokenizeQuery(query string) []string {
	query = strings.ToLower(query)
	query = nonAlphanumericRegex.ReplaceAllString(query, " ")
	rawTokens := strings.Fields(query)

	var filtered []string
	for _, token := range rawTokens {
		if _, exists := stopWords[token]; !exists && len(token) > 1 {
			filtered = append(filtered, token)
		}
	}
	return filtered
}