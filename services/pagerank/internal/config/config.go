package config

import (
	"os"
	"strconv"
)

type Config struct {
	MongoURI         string
	DBName           string
	SourceCollection string
	OutputCollection string
	DampingFactor    float64
	MaxIterations    int
}

func LoadConfig() *Config {
	return &Config{
		MongoURI:         getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DBName:           getEnv("DB_NAME", "swiftsearch"),
		// CHANGED: Now matches COLLECTION_NAME from your docker-compose.yml
		SourceCollection: getEnv("COLLECTION_NAME", "pages"),
		OutputCollection: getEnv("OUTPUT_COLLECTION", "pagerank_scores"),
		DampingFactor:    getEnvAsFloat("DAMPING_FACTOR", 0.85),
		MaxIterations:    getEnvAsInt("MAX_ITERATIONS", 25),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
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