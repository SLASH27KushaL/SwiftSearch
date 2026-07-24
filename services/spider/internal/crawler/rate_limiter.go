package crawler

import (
	"sync"
	"time"
)

// RateLimiter ensures the crawler behaves politely and doesn't hammer domains.
type RateLimiter struct {
	domains sync.Map
	delay   time.Duration
}

// NewRateLimiter creates a new domain-based rate limiter.
func NewRateLimiter(delay time.Duration) *RateLimiter {
	return &RateLimiter{
		delay: delay,
	}
}

// Wait blocks the current goroutine until it is safe to crawl the given domain again.
func (rl *RateLimiter) Wait(domain string) {
	now := time.Now()

	// Check the last time we hit this domain
	if lastSeen, loaded := rl.domains.LoadOrStore(domain, now); loaded {
		lastTime := lastSeen.(time.Time)
		elapsed := now.Sub(lastTime)

		// If we are too fast, sleep for the remaining duration
		if elapsed < rl.delay {
			time.Sleep(rl.delay - elapsed)
		}

		// Update the timestamp to now
		rl.domains.Store(domain, time.Now())
	}
}