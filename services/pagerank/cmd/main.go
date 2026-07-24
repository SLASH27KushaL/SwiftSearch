package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"moogle-go/services/pagerank/internal/config"
	"moogle-go/services/pagerank/internal/pagerank"
	"moogle-go/services/pagerank/internal/store"
	"moogle-go/services/pagerank/pkg/models"
)

func main() {
	log.Println("Initializing Moogle PageRank Service...")

	cfg := config.LoadConfig()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Setup Data Stores
	reader, err := store.NewMongoReader(ctx, cfg.MongoURI, cfg.DBName, cfg.SourceCollection)
	if err != nil {
		log.Fatalf("Failed to init reader: %v", err)
	}
	defer reader.Close(ctx)

	writer, err := store.NewMongoWriter(ctx, cfg.MongoURI, cfg.DBName, cfg.OutputCollection)
	if err != nil {
		log.Fatalf("Failed to init writer: %v", err)
	}
	defer writer.Close(ctx)

	// Catch graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stopChan
		log.Println("Shutting down...")
		cancel()
	}()

	// 2. Fetch Graph Nodes
	log.Println("Fetching web graph from database...")
	nodes, err := reader.FetchAllNodes(ctx)
	if err != nil {
		log.Fatalf("Failed to fetch nodes: %v", err)
	}
	log.Printf("Fetched %d raw nodes.", len(nodes))

	// 3. Build Adjacency List
	log.Println("Building in-memory graph...")
	graph := pagerank.BuildGraph(nodes)
	log.Printf("Graph built with %d unique URLs.", len(graph))

	// 4. Calculate PageRank
	log.Printf("Calculating PageRank (Damping: %.2f, Iterations: %d)...", cfg.DampingFactor, cfg.MaxIterations)
	startTime := time.Now()
	scoresMap := pagerank.Calculate(graph, cfg.DampingFactor, cfg.MaxIterations)
	log.Printf("Calculation complete in %v", time.Since(startTime))

	// 5. Transform Map to Structs for Writing
	var finalScores []models.PageRankScore
	for url, score := range scoresMap {
		finalScores = append(finalScores, models.PageRankScore{
			URL:   url,
			Score: score,
		})
	}

	// 6. Save Scores to Database
	log.Println("Saving PageRank scores to database...")
	if err := writer.BulkUpsertScores(ctx, finalScores); err != nil {
		log.Fatalf("Failed to write scores: %v", err)
	}

	log.Println("PageRank successfully calculated and saved!")
}