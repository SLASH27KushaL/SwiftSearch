package store

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"moogle-go/services/indexer/pkg/models"
)

type MongoReader struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoReader(ctx context.Context, uri, dbName, collectionName string) (*MongoReader, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("reader mongo connection failed: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("reader mongo ping failed: %w", err)
	}

	collection := client.Database(dbName).Collection(collectionName)

	return &MongoReader{
		client:     client,
		collection: collection,
	}, nil
}

func (m *MongoReader) GetNextBatch(ctx context.Context, batchSize int) ([]models.RawPage, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"indexed": bson.M{"$exists": false}},
			{"indexed": false},
		},
	}

	findOptions := options.Find().SetLimit(int64(batchSize))

	cursor, err := m.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page batch: %w", err)
	}
	defer cursor.Close(ctx)

	var pages []models.RawPage
	if err := cursor.All(ctx, &pages); err != nil {
		return nil, fmt.Errorf("failed to decode page batch: %w", err)
	}

	return pages, nil
}

func (m *MongoReader) MarkAsProcessed(ctx context.Context, url string) error {
	filter := bson.M{"url": url}
	update := bson.M{"$set": bson.M{"indexed": true}}

	_, err := m.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to mark document as indexed: %w", err)
	}
	return nil
}

func (m *MongoReader) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}