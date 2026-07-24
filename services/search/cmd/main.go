package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"moogle-go/services/search/internal/api"
	"moogle-go/services/search/internal/config"
	"moogle-go/services/search/internal/engine"
	"moogle-go/services/search/internal/store"
)

func main() {
	log.Println("Initializing Moogle Search API Service...")
	cfg := config.LoadConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Initialize MongoDB Reader
	mongoReader, err := store.NewMongoReader(ctx, cfg.MongoURI, cfg.DBName, cfg.IndexCollection, cfg.RankCollection)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoReader.Close(context.Background())
	log.Println("Connected to MongoDB Search Clusters")

	// 2. Initialize Redis Cache
	redisCache := store.NewRedisCache(cfg.RedisURI)
	log.Println("Initialized Redis Caching Layer")

	// 3. Initialize Ranking Engine
	ranker := engine.NewRanker(mongoReader, cfg.AlphaWeight, cfg.BetaWeight)

	// 4. Initialize HTTP Handler
	handler := api.NewSearchHandler(ranker, redisCache)

	// 5. Setup Gin Web Framework
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// CORS middleware for frontend access
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	})

	// Define API Routes
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Engine is operational"})
	})

	router.GET("/search", handler.HandleSearch)

	// 6. Start the Server
	log.Printf("Moogle Search API running on http://localhost:%s\n", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}