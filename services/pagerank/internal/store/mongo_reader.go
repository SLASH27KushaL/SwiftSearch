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

// FetchAllNodes streams all documents, extracting only the URL and Outlinks
func (m *MongoReader) FetchAllNodes(ctx context.Context) ([]models.GraphNode, error) {
	// Project only necessary fields to save memory
	opts := options.Find().SetProjection(bson.M{"url": 1, "outlinks": 1, "_id": 0})

	cursor, err := m.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch nodes: %w", err)
	}
	defer cursor.Close(ctx)

	var nodes []models.GraphNode
	if err = cursor.All(ctx, &nodes); err != nil {
		return nil, fmt.Errorf("failed to decode nodes: %w", err)
	}

	return nodes, nil
}

func (m *MongoReader) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}