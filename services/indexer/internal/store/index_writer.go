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

type IndexWriter struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewIndexWriter(ctx context.Context, uri, dbName, collectionName string) (*IndexWriter, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("writer mongo connection failed: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("writer mongo ping failed: %w", err)
	}

	collection := client.Database(dbName).Collection(collectionName)
	writer := &IndexWriter{
		client:     client,
		collection: collection,
	}

	if err := writer.createIndexes(ctx); err != nil {
		return nil, fmt.Errorf("failed to create writer indexes: %w", err)
	}

	return writer, nil
}

func (i *IndexWriter) createIndexes(ctx context.Context) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "term", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := i.collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (i *IndexWriter) BulkUpsert(ctx context.Context, entries []models.IndexEntry) error {
	if len(entries) == 0 {
		return nil
	}

	var operations []mongo.WriteModel

	for _, entry := range entries {
		filter := bson.M{"term": entry.Term}
		// FIX: We now grab the first (and only) match from the slice to push to the DB
		update := bson.M{
			"$push": bson.M{"matches": entry.Matches[0]},
		}

		op := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)

		operations = append(operations, op)
	}

	bulkOptions := options.BulkWrite().SetOrdered(false)
	_, err := i.collection.BulkWrite(ctx, operations, bulkOptions)
	if err != nil {
		return fmt.Errorf("bulk write failed: %w", err)
	}

	return nil
}

func (i *IndexWriter) Close(ctx context.Context) error {
	return i.client.Disconnect(ctx)
}