package config

import (
	"os"
	"strconv"
)

type Config struct {
	WorkerCount    int
	MaxDepth       int
	CrawlDelayMS   int
	UserAgent      string
	RedisAddr      string
	RedisPassword  string
	RedisDB        int
	MongoURI       string
	DBName         string
	CollectionName string
}

func LoadConfig() Config {
	return Config{
		WorkerCount:    getEnvAsInt("WORKER_COUNT", 10),
		MaxDepth:       getEnvAsInt("MAX_DEPTH", 3),
		CrawlDelayMS:   getEnvAsInt("CRAWL_DELAY_MS", 500),
		UserAgent:      getEnv("USER_AGENT", "MoogleSpiderBot/1.0 (+https://github.com/yourusername/moogle)"),
		RedisAddr:      getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
		RedisDB:        getEnvAsInt("REDIS_DB", 0),
		MongoURI:       getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DBName:         getEnv("DB_NAME", "moogle_db"),
		CollectionName: getEnv("COLLECTION_NAME", "raw_pages"),
	}
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists && val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}