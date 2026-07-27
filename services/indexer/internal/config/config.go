package config

import (
	"os"
	"strconv"
)

type Config struct {
	MongoURI         string
	DBName           string
	SourceCollection string
	IndexCollection  string
	BatchSize        int
}

func LoadConfig() *Config {
	return &Config{
		MongoURI:         getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DBName:           getEnv("DB_NAME", "swiftsearch"),
		// CHANGED: Now perfectly matches your docker-compose.yml
		SourceCollection: getEnv("RAW_COLLECTION", "pages"),
		IndexCollection:  getEnv("INDEX_COLLECTION", "index"),
		BatchSize:        getEnvAsInt("BATCH_SIZE", 100),
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