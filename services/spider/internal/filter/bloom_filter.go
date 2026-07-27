package filter

import (
	"context"
	"hash/fnv"

	"github.com/redis/go-redis/v9"
)

type BloomFilter struct {
	client *redis.Client
	ctx    context.Context
	key    string
}

func NewBloomFilter(redisURI string) *BloomFilter {
	client := redis.NewClient(&redis.Options{Addr: redisURI})
	return &BloomFilter{
		client: client,
		ctx:    context.Background(),
		key:    "spider:bloom:visited",
	}
}

// getBitOffset generates a pseudo-random bit index for a URL
func (b *BloomFilter) getBitOffset(url string) int64 {
	h := fnv.New32a()
	h.Write([]byte(url))
	// Modulo to keep it within Redis bitmap size limits
	return int64(h.Sum32()) % 100000000
}

// CheckAndAdd returns true if the URL was ALREADY visited.
// If not, it marks it as visited in O(1) time.
func (b *BloomFilter) CheckAndAdd(url string) (bool, error) {
	offset := b.getBitOffset(url)
	// SETBIT returns 1 if the bit was already set, 0 if it was false
	wasSet, err := b.client.SetBit(b.ctx, b.key, offset, 1).Result()
	return wasSet == 1, err
}
