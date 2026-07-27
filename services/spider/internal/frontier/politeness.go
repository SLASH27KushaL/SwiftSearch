package frontier

import (
	"context"
	"net/url"
	"time"

	"github.com/redis/go-redis/v9"
)

type PolitenessPolicy struct {
	client *redis.Client
	ctx    context.Context
	delay  time.Duration
}

// NewPolitenessPolicy initializes the rate limiter. delaySeconds dictates how long
// the entire cluster must wait before hitting the same domain again.
func NewPolitenessPolicy(redisURI string, delaySeconds int) *PolitenessPolicy {
	client := redis.NewClient(&redis.Options{Addr: redisURI})
	return &PolitenessPolicy{
		client: client,
		ctx:    context.Background(),
		delay:  time.Duration(delaySeconds) * time.Second,
	}
}

// ExtractDomain parses a raw URL and returns just the hostname (e.g., "en.wikipedia.org")
func ExtractDomain(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return parsedURL.Hostname()
}

// RequestAccess checks if a domain is safe to crawl right now.
// If safe, it locks the domain for the delay period and returns true.
func (p *PolitenessPolicy) RequestAccess(domain string) bool {
	if domain == "" {
		return false
	}

	key := "spider:politeness:" + domain

	// SetNX (Set if Not eXists) tries to set the key.
	// If the key exists (meaning we hit this domain recently), it returns false.
	// If it doesn't exist, it sets it with the TTL delay and returns true.
	acquired, err := p.client.SetNX(p.ctx, key, "locked", p.delay).Result()
	if err != nil {
		return false
	}

	return acquired
}
