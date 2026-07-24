package storage

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"moogle-go/services/spider/pkg/models"
)

// MongoStore handles page document persistence in MongoDB.
type MongoStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewMongoStore initializes a MongoDB connection, sets up indexes, and returns a MongoStore instance.
func NewMongoStore(ctx context.Context, uri, dbName, collectionName string) (*MongoStore, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("mongo connection failed: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("mongo ping failed: %w", err)
	}

	collection := client.Database(dbName).Collection(collectionName)

	store := &MongoStore{
		client:     client,
		collection: collection,
	}

	// Ensure unique index on URL field
	if err := store.createIndexes(ctx); err != nil {
		return nil, fmt.Errorf("failed to create mongo indexes: %w", err)
	}

	return store, nil
}

// createIndexes creates database indexes for unique URLs.
func (m *MongoStore) createIndexes(ctx context.Context) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "url", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := m.collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

// SavePage upserts a crawled page document into MongoDB based on its unique URL.
func (m *MongoStore) SavePage(ctx context.Context, page *models.Page) error {
	filter := bson.D{{Key: "url", Value: page.URL}}
	update := bson.D{{Key: "$set", Value: page}}

	opts := options.UpdateOne().SetUpsert(true)

	_, err := m.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to upsert page document: %w", err)
	}

	return nil
}

// Close disconnects the MongoDB client cleanly.
func (m *MongoStore) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}