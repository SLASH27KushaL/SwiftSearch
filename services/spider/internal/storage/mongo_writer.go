package store

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoWriter struct {
	collection *mongo.Collection
}

func NewMongoWriter(mongoURI string) *MongoWriter {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}

	// Storing raw HTML in the 'pages' collection
	col := client.Database("swiftsearch").Collection("pages")
	return &MongoWriter{collection: col}
}

func (m *MongoWriter) SavePage(url, htmlContent string) error {
	doc := map[string]interface{}{
		"url":        url,
		"html":       htmlContent,
		"crawled_at": time.Now(),
	}
	_, err := m.collection.InsertOne(context.Background(), doc)
	return err
}
