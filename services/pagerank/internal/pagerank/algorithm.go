package pagerank

// Calculate runs the iterative PageRank algorithm on an adjacency list graph.
func Calculate(graph map[string][]string, dampingFactor float64, maxIterations int) map[string]float64 {
	N := float64(len(graph))
	if N == 0 {
		return nil
	}

	// 1. Initialize all nodes with a base probability of 1/N
	ranks := make(map[string]float64, len(graph))
	for node := range graph {
		ranks[node] = 1.0 / N
	}

	// 2. Iterate to converge the scores
	for i := 0; i < maxIterations; i++ {
		newRanks := make(map[string]float64, len(graph))

		// The baseline probability a user randomly teleports to this page
		baseProb := (1.0 - dampingFactor) / N

		for node := range graph {
			newRanks[node] = baseProb
		}

		// Distribute current node's rank to its outbound links
		for node, outlinks := range graph {
			outDegree := float64(len(outlinks))

			if outDegree > 0 {
				// Node shares its rank equally among all its outbound links
				share := ranks[node] / outDegree
				for _, outlink := range outlinks {
					newRanks[outlink] += dampingFactor * share
				}
			} else {
				// Sink Node (No outlinks): distribute its rank to ALL nodes equally
				// (Simulating a user randomly typing a new URL when they hit a dead end)
				sinkShare := (ranks[node] * dampingFactor) / N
				for allNodes := range graph {
					newRanks[allNodes] += sinkShare
				}
			}
		}

		ranks = newRanks
	}

	return ranks
}