package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"moogle-go/pkg/models"
	"moogle-go/services/spider/internal/config"
	"moogle-go/services/spider/internal/crawler"
	"moogle-go/services/spider/internal/parser"
	"moogle-go/services/spider/internal/storage"
)

func main() {
	log.Println(" Starting Moogle Web Spider...")

	// 1. Load Configurations
	cfg := config.LoadConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 2. Initialize Redis Visited Store (Deduplication)
	visitedStore, err := storage.NewRedisVisitedStore(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer visitedStore.Close()
	log.Println(" Connected to Redis Visited Store")

	// 3. Initialize Document Persistence Store (MongoDB)
	docStore, err := storage.NewMongoStore(ctx, cfg.MongoURI, cfg.DBName, cfg.CollectionName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := docStore.Close(ctx); err != nil {
			log.Printf("Error closing Mongo connection: %v", err)
		}
	}()
	log.Println(" Connected to MongoDB Page Store")

	// 4. Initialize Parser Services
	normalizer := parser.NewURLNormalizer()
	htmlParser := parser.NewHTMLParser(normalizer)

	// 5. Build Engine Options & Instantiate Crawler Engine
	engineOpts := crawler.EngineOptions{
		WorkerCount:  cfg.WorkerCount,
		MaxDepth:     cfg.MaxDepth,
		CrawlDelay:   time.Duration(cfg.CrawlDelayMS) * time.Millisecond,
		UserAgent:    cfg.UserAgent,
		VisitedStore: visitedStore,
		DocStore:     docStore,
		Parser:       htmlParser,
		Normalizer:   normalizer,
	}

	engine := crawler.NewEngine(engineOpts)

	// 6. Handle Graceful Shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stopChan
		log.Println(" Received shutdown signal. Stopping crawler engine gracefully...")
		engine.Stop()
		cancel()
	}()

	// 7. Seed Initial URLs to Start Crawling
	seedJobs := []models.CrawlJob{
		{URL: "https://go.dev", Depth: 0},
		{URL: "https://wikipedia.org", Depth: 0},
		{URL: "https://news.ycombinator.com", Depth: 0},
	}

	go func() {
		for _, job := range seedJobs {
			engine.SubmitJob(job)
		}
	}()

	// 8. Start Crawling Routine (blocks until engine shuts down)
	log.Printf(" Starting engine with %d workers (Max Depth: %d)...", cfg.WorkerCount, cfg.MaxDepth)
	engine.Start(ctx)

	log.Println(" Spider service terminated cleanly.")
}