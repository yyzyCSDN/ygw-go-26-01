package limiter

// Summary describes one tenant bucket.
type Summary struct {
	Tenant string
	Tokens float64
}

// Summaries returns a snapshot of all tenant buckets.
func (l *Limiter) Summaries() []Summary {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]Summary, 0, len(l.buckets))
	for tenant, b := range l.buckets {
		out = append(out, Summary{Tenant: tenant, Tokens: b.tokens})
	}
	return out
}
