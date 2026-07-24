package tokenizer

import (
	"regexp"
	"strings"
)

var (
	// nonAlphanumericRegex matches anything that is NOT a letter, number, or whitespace.
	// We compile it once at startup for performance.
	nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
)

// CleanAndTokenize normalizes raw text and splits it into actionable word tokens.
func CleanAndTokenize(text string) []string {
	// 1. Convert everything to lowercase to ensure "Apple" and "apple" are the same token
	text = strings.ToLower(text)

	// 2. Replace all punctuation and special characters with spaces
	text = nonAlphanumericRegex.ReplaceAllString(text, " ")

	// 3. Split the text by whitespace (this automatically handles multiple spaces/newlines)
	rawTokens := strings.Fields(text)

	// 4. Filter out extremely short or absurdly long junk tokens
	var finalTokens []string
	for _, token := range rawTokens {
		// Keep words between 2 and 45 characters long (longest English word is 45 chars)
		if len(token) > 1 && len(token) <= 45 {
			finalTokens = append(finalTokens, token)
		}
	}

	return finalTokens
}
