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

type MongoWriter struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoWriter(ctx context.Context, uri, dbName, collectionName string) (*MongoWriter, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("writer connection failed: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("writer ping failed: %w", err)
	}

	writer := &MongoWriter{
		client:     client,
		collection: client.Database(dbName).Collection(collectionName),
	}

	// Create a unique index on URL for fast searching later
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "url", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, _ = writer.collection.Indexes().CreateOne(ctx, indexModel)

	return writer, nil
}

func (m *MongoWriter) BulkUpsertScores(ctx context.Context, scores []models.PageRankScore) error {
	if len(scores) == 0 {
		return nil
	}

	var operations []mongo.WriteModel
	for _, s := range scores {
		filter := bson.M{"url": s.URL}
		update := bson.M{"$set": bson.M{"score": s.Score}}

		op := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)

		operations = append(operations, op)
	}

	opts := options.BulkWrite().SetOrdered(false)
	_, err := m.collection.BulkWrite(ctx, operations, opts)
	return err
}

func (m *MongoWriter) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}