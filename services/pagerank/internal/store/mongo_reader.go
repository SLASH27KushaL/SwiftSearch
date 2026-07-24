package store

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"moogle-go/services/pagerank/pkg/models"
)

type MongoReader struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoReader(ctx context.Context, uri, dbName, collectionName string) (*MongoReader, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("reader connection failed: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("reader ping failed: %w", err)
	}

	return &MongoReader{
		client:     client,
		collection: client.Database(dbName).Collection(collectionName),
	}, nil
}

// FetchAllNodes streams all documents, handling flexible outlink schemas
func (m *MongoReader) FetchAllNodes(ctx context.Context) ([]models.GraphNode, error) {
	opts := options.Find().SetProjection(bson.M{"url": 1, "outlinks": 1, "_id": 0})

	cursor, err := m.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch nodes: %w", err)
	}
	defer cursor.Close(ctx)

	// Read as flexible BSON maps instead of strictly typed structs
	var rawNodes []bson.M
	if err = cursor.All(ctx, &rawNodes); err != nil {
		return nil, fmt.Errorf("failed to decode nodes: %w", err)
	}

	var nodes []models.GraphNode
	for _, raw := range rawNodes {
		url, _ := raw["url"].(string)
		var outlinks []string

		// Safely parse the outlinks array, handling both strings and embedded objects
		if rawOutlinks, ok := raw["outlinks"].(bson.A); ok {
			for _, item := range rawOutlinks {
				// Scenario A: Outlink is an embedded document (e.g., {"url": "https..."})
				if doc, ok := item.(bson.M); ok {
					if link, ok := doc["url"].(string); ok {
						outlinks = append(outlinks, link)
					} else if link, ok := doc["href"].(string); ok {
						outlinks = append(outlinks, link)
					}
				} else if str, ok := item.(string); ok {
					// Scenario B: Outlink is already just a string
					outlinks = append(outlinks, str)
				}
			}
		}

		nodes = append(nodes, models.GraphNode{
			URL:      url,
			Outlinks: outlinks,
		})
	}

	return nodes, nil
}

func (m *MongoReader) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}