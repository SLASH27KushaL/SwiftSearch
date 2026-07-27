package config

import "os"

type Config struct {
	RedisURI    string
	MongoURI    string
	WorkerCount int
}

func LoadConfig() *Config {
	return &Config{
		RedisURI:    getEnv("REDIS_URI", "localhost:6379"),
		MongoURI:    getEnv("MONGO_URI", "mongodb://localhost:27017"),
		WorkerCount: 10, // Number of concurrent goroutines per node
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
