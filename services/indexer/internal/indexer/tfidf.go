package indexer

// CalculateTF calculates the Term Frequency (TF) for a slice of filtered tokens.
// It returns a map where the key is the word (term) and the value is its TF score.
func CalculateTF(tokens []string) map[string]float64 {
	tfMap := make(map[string]float64)
	totalTokens := float64(len(tokens))

	if totalTokens == 0 {
		return tfMap
	}

	// Step 1: Count occurrences of each unique term
	termCounts := make(map[string]int)
	for _, token := range tokens {
		termCounts[token]++
	}

	// Step 2: Calculate the Term Frequency ratio
	for term, count := range termCounts {
		tfMap[term] = float64(count) / totalTokens
	}

	return tfMap
}
