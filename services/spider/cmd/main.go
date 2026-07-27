package main

import (
	"log"

	"swiftsearch/spider/config"
	"swiftsearch/spider/internal/filter"
	"swiftsearch/spider/internal/frontier"
	"swiftsearch/spider/internal/parser"
	"swiftsearch/spider/internal/store"
	"swiftsearch/spider/internal/worker"
)

func main() {
	log.Println("Initializing Distributed Moogle Spider...")

	// 1. Load Configurations (Redis URI, Mongo URI, Worker Count)
	cfg := config.LoadConfig()

	// 2. Initialize all Micro-components
	queue := frontier.NewRedisQueue(cfg.RedisURI)
	bloom := filter.NewBloomFilter(cfg.RedisURI)
	mongoDB := store.NewMongoWriter(cfg.MongoURI)

	// Enforce a 3-second delay per domain to prevent IP bans
	politeness := frontier.NewPolitenessPolicy(cfg.RedisURI, 3)

	// Respect domain rules using our custom bot identity
	robotsChecker := parser.NewRobotsChecker("SwiftBot")

	// 3. Seed the queue with starting URLs
	// (It's safe to run this multiple times; the Bloom filter will drop duplicates anyway)
	queue.Push("https://en.wikipedia.org/wiki/Distributed_computing")
	queue.Push("https://golang.org/doc/")
	queue.Push("https://pkg.go.dev/")

	// 4. Assemble and Ignite the Worker Pool
	pool := worker.NewPool(queue, bloom, mongoDB, politeness, robotsChecker)

	// This will block the main thread and keep your workers running indefinitely
	pool.Start(cfg.WorkerCount)
}
