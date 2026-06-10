package ratelimiter

import (
	"golang.org/x/time/rate"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	clients  map[string]*clientLimiter
	rate     rate.Limit
	burst    int
	ttl      time.Duration
	cleanup  *time.Ticker
	doneChan chan struct{}
}

type clientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// NewRateLimiter creates a new rateLimiter
func NewRateLimiter(cfg *Config) *RateLimiter {
	rl := &RateLimiter{
		clients:  make(map[string]*clientLimiter),
		rate:     cfg.Limit,
		burst:    cfg.Burst,
		ttl:      cfg.TTL,
		cleanup:  time.NewTicker(cfg.TTL),
		doneChan: make(chan struct{}),
	}
	// Start background cleanup
	go rl.cleanupExpiredClients()
	return rl
}

// getLimiter returns the limiter for the given key, creating a new one if necessary
func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if client, exists := rl.clients[key]; exists {
		client.lastSeen = time.Now()
		return client.limiter
	}

	limiter := rate.NewLimiter(rl.rate, rl.burst)
	rl.clients[key] = &clientLimiter{limiter: limiter, lastSeen: time.Now()}
	return limiter
}

// cleanupExpiredClients removes clients that haven't been used recently
func (rl *RateLimiter) cleanupExpiredClients() {
	for {
		select {
		case <-rl.cleanup.C:
			rl.mu.Lock()
			now := time.Now()
			for key, client := range rl.clients {
				if now.Sub(client.lastSeen) > rl.ttl {
					delete(rl.clients, key)
				}
			}
			rl.mu.Unlock()
		case <-rl.doneChan:
			rl.cleanup.Stop()
			return
		}
	}
}

// stopCleanup stops the background cleanup
func (rl *RateLimiter) stopCleanup() {
	close(rl.doneChan)
}

// Middleware returns a rate-limiting middleware that accepts a key-extraction function
func (rl *RateLimiter) Middleware(keyFunc func(*http.Request) string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := keyFunc(r)
			if key == "" {
				http.Error(w, "Rate limiting key not provided", http.StatusForbidden)
				return
			}

			limiter := rl.getLimiter(key)

			if !limiter.Allow() {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
