package limiter

// Refill grants the tenant one token without waiting for the next tick.
func (l *Limiter) Refill(tenant string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	b, ok := l.buckets[tenant]
	if !ok {
		b = &bucket{tokens: l.burst, last: l.now()}
		l.buckets[tenant] = b
	}
	if b.tokens < l.burst {
		b.tokens++
	}
}

// Drain removes the tenant's bucket entirely.
func (l *Limiter) Drain(tenant string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.buckets, tenant)
}
