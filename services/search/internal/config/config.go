package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port            string
	MongoURI        string
	DBName          string
	IndexCollection string
	RankCollection  string
	RedisURI        string
	AlphaWeight     float64 // Importance of Term Frequency
	BetaWeight      float64 // Importance of PageRank
}

func LoadConfig() *Config {
	return &Config{
		Port:            getEnv("PORT", "8080"),
		// CHANGED: 127.0.0.1 forces IPv4, 27017 matches your docker container
		MongoURI:        getEnv("MONGO_URI", "mongodb://127.0.0.1:27017"),
		// CHANGED: Matches the DB where spider, indexer, and pagerank saved data
		DBName:          getEnv("DB_NAME", "swiftsearch"),
		// CHANGED: Matches the Indexer's output collection
		IndexCollection: getEnv("INDEX_COLLECTION", "index"),
		RankCollection:  getEnv("RANK_COLLECTION", "pagerank_scores"),
		// CHANGED: 127.0.0.1 bypasses the Windows [::1] IPv6 bug
		RedisURI:        getEnv("REDIS_URI", "127.0.0.1:6380"),
		AlphaWeight:     getEnvAsFloat("ALPHA_WEIGHT", 0.6), // Default 60% text relevance
		BetaWeight:      getEnvAsFloat("BETA_WEIGHT", 0.4),  // Default 40% graph authority
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsFloat(key string, fallback float64) float64 {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
			return value
		}
	}
	return fallback
}