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

// IndexWriter handles saving the inverted index to MongoDB.
type IndexWriter struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewIndexWriter connects to MongoDB, ensures indexes, and returns a writer instance.
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

	// Create a unique index on the "term" field for fast lookups
	if err := writer.createIndexes(ctx); err != nil {
		return nil, fmt.Errorf("failed to create writer indexes: %w", err)
	}

	return writer, nil
}

// createIndexes ensures that each word (term) is completely unique in the database.
func (i *IndexWriter) createIndexes(ctx context.Context) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "term", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := i.collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

// BulkUpsert takes a batch of IndexEntries and pushes them into the inverted index.
func (i *IndexWriter) BulkUpsert(ctx context.Context, entries []models.IndexEntry) error {
	if len(entries) == 0 {
		return nil
	}

	var operations []mongo.WriteModel

	for _, entry := range entries {
		filter := bson.M{"term": entry.Term}
		// $push appends the document match (URL, Title, TF) into the term's "matches" array
		update := bson.M{
			"$push": bson.M{"matches": entry.DocumentMatch},
		}

		op := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)

		operations = append(operations, op)
	}

	// Execute all writes in a single network round-trip for massive performance gains
	bulkOptions := options.BulkWrite().SetOrdered(false)
	_, err := i.collection.BulkWrite(ctx, operations, bulkOptions)
	if err != nil {
		return fmt.Errorf("bulk write failed: %w", err)
	}

	return nil
}

// Close disconnects the MongoDB client cleanly.
func (i *IndexWriter) Close(ctx context.Context) error {
	return i.client.Disconnect(ctx)
}
