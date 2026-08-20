// Package limiter implements a per-tenant token bucket.
package limiter

import (
	"sync"
	"time"
)

type bucket struct {
	tokens float64
	last   time.Time
}

// Limiter rate-limits per tenant with a token bucket.
type Limiter struct {
	mu      sync.Mutex
	rate    float64
	burst   float64
	now     func() time.Time
	buckets map[string]*bucket
}

// New creates a limiter with rate tokens per second and the given burst.
func New(rate float64, burst float64) *Limiter {
	return &Limiter{
		rate:    rate,
		burst:   burst,
		now:     time.Now,
		buckets: make(map[string]*bucket),
	}
}

// Allow reports whether the tenant may take one token now.
func (l *Limiter) Allow(tenant string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	b, ok := l.buckets[tenant]
	if !ok {
		b = &bucket{tokens: l.burst, last: l.now()}
		l.buckets[tenant] = b
	}
	now := l.now()
	elapsed := now.Sub(b.last).Seconds()
	b.last = now
	b.tokens += elapsed * l.rate
	if b.tokens > l.burst {
		b.tokens = l.burst
	}
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}
