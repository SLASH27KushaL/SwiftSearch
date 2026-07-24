package storage

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const VisitedSetKey = "spider:visited_urls"

// RedisVisitedStore implements deduplication using Redis Sets.
type RedisVisitedStore struct {
	client *redis.Client
}

// NewRedisVisitedStore connects to Redis and returns a new store instance.
func NewRedisVisitedStore(addr, password string, db int) (*RedisVisitedStore, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &RedisVisitedStore{client: rdb}, nil
}

// IsVisited checks if a normalized URL exists in the Redis visited set.
func (r *RedisVisitedStore) IsVisited(ctx context.Context, rawURL string) (bool, error) {
	exists, err := r.client.SIsMember(ctx, VisitedSetKey, rawURL).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check visited state in redis: %w", err)
	}
	return exists, nil
}

// MarkVisited adds a normalized URL into the Redis visited set.
func (r *RedisVisitedStore) MarkVisited(ctx context.Context, rawURL string) error {
	if err := r.client.SAdd(ctx, VisitedSetKey, rawURL).Err(); err != nil {
		return fmt.Errorf("failed to mark url as visited in redis: %w", err)
	}
	return nil
}

// Close closes the underlying Redis client connection.
func (r *RedisVisitedStore) Close() error {
	return r.client.Close()
}