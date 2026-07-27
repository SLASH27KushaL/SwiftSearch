package frontier

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisQueue struct {
	client *redis.Client
	ctx    context.Context
	queue  string
}

func NewRedisQueue(redisURI string) *RedisQueue {
	client := redis.NewClient(&redis.Options{Addr: redisURI})
	return &RedisQueue{
		client: client,
		ctx:    context.Background(),
		queue:  "spider:frontier:unvisited",
	}
}

// Push adds a URL to the end of the queue
func (r *RedisQueue) Push(url string) error {
	return r.client.RPush(r.ctx, r.queue, url).Err()
}

// Pop blocks until a URL is available, then returns it
func (r *RedisQueue) Pop() (string, error) {
	result, err := r.client.BLPop(r.ctx, 0, r.queue).Result()
	if err != nil {
		return "", err
	}
	// BLPop returns [queue_name, value]
	return result[1], nil
}
