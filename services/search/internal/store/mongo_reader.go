package store

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"moogle-go/services/search/pkg/models"
)

type MongoReader struct {
	client    *mongo.Client
	db        *mongo.Database
	indexColl *mongo.Collection
	rankColl  *mongo.Collection
}

func NewMongoReader(ctx context.Context, uri, dbName, indexColl, rankColl string) (*MongoReader, error) {
	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	return &MongoReader{
		client:    client,
		db:        db,
		indexColl: db.Collection(indexColl),
		rankColl:  db.Collection(rankColl),
	}, nil
}

// FetchTFIDF matches tokens to URLs and returns a map of URL -> combined TF score
func (m *MongoReader) FetchTFIDF(ctx context.Context, tokens []string) (map[string]float64, map[string]string, error) {
	filter := bson.M{"term": bson.M{"$in": tokens}}
	cursor, err := m.indexColl.Find(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	var entries []models.IndexEntry
	if err = cursor.All(ctx, &entries); err != nil {
		return nil, nil, err
	}

	urlScores := make(map[string]float64)
	urlTitles := make(map[string]string)

	for _, entry := range entries {
		for _, match := range entry.Matches {
			urlScores[match.URL] += match.TF // Additive TF for multi-word queries
			urlTitles[match.URL] = match.Title
		}
	}
	return urlScores, urlTitles, nil
}

// FetchPageRanks fetches the graph importance for a specific list of URLs
func (m *MongoReader) FetchPageRanks(ctx context.Context, urls []string) (map[string]float64, error) {
	if len(urls) == 0 {
		return make(map[string]float64), nil
	}

	filter := bson.M{"url": bson.M{"$in": urls}}
	cursor, err := m.rankColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var ranks []models.PageRankScore
	if err = cursor.All(ctx, &ranks); err != nil {
		return nil, err
	}

	prMap := make(map[string]float64)
	for _, r := range ranks {
		prMap[r.URL] = r.Score
	}
	return prMap, nil
}

func (m *MongoReader) Close(ctx context.Context) {
	m.client.Disconnect(ctx)
}