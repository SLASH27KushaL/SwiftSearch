package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"moogle-go/services/indexer/internal/config"
	"moogle-go/services/indexer/internal/indexer"
	"moogle-go/services/indexer/internal/store"
)

func main() {
	log.Println("Initializing Moogle Inverted Indexer Service...")

	cfg := config.LoadConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reader, err := store.NewMongoReader(ctx, cfg.MongoURI, cfg.DBName, cfg.SourceCollection)
	if err != nil {
		log.Fatalf("Failed to connect to Mongo Reader: %v", err)
	}
	defer reader.Close(ctx)
	log.Println("Connected to MongoDB Source Pages Store")

	writer, err := store.NewIndexWriter(ctx, cfg.MongoURI, cfg.DBName, cfg.IndexCollection)
	if err != nil {
		log.Fatalf("Failed to connect to Index Writer: %v", err)
	}
	defer writer.Close(ctx)
	log.Println("Connected to MongoDB Inverted Index Store")

	engine := indexer.NewEngine(reader, writer, cfg.BatchSize)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stopChan
		log.Println("Received shutdown signal. Stopping indexer gracefully...")
		cancel()
	}()

	log.Println("Starting indexing pipeline...")
	startTime := time.Now()

	if err := engine.Run(ctx); err != nil {
		log.Fatalf("Indexing pipeline failed: %v", err)
	}

	log.Printf("Indexing process completed successfully in %v", time.Since(startTime))
}
