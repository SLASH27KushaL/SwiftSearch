package pagerank

import (
	"moogle-go/services/pagerank/pkg/models"
)

// BuildGraph transforms database models into an in-memory adjacency list.
// Returns a map where Key = URL, Value = Slice of URLs it links to.
func BuildGraph(nodes []models.GraphNode) map[string][]string {
	graph := make(map[string][]string, len(nodes))

	for _, node := range nodes {
		// Ensure the node exists in the graph even if it has no outlinks
		if _, exists := graph[node.URL]; !exists {
			graph[node.URL] = []string{}
		}

		for _, outlink := range node.Outlinks {
			graph[node.URL] = append(graph[node.URL], outlink)

			// Initialize target nodes in the graph to avoid dead ends dropping out
			if _, exists := graph[outlink]; !exists {
				graph[outlink] = []string{}
			}
		}
	}

	return graph
}