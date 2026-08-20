package limiter

import "sync"

// Registry keeps per-tenant limiters.
type Registry struct {
	mu      sync.Mutex
	rate    float64
	burst   float64
	tenants map[string]*Limiter
}

// NewRegistry creates a registry with shared rate and burst.
func NewRegistry(rate, burst float64) *Registry {
	if burst <= 0 {
		burst = BurstFor(rate, 1)
	}
	return &Registry{
		rate:    rate,
		burst:   burst,
		tenants: make(map[string]*Limiter),
	}
}

// LimiterFor returns the tenant's limiter, creating it on first use.
func (r *Registry) LimiterFor(tenant string) *Limiter {
	r.mu.Lock()
	defer r.mu.Unlock()
	l, ok := r.tenants[tenant]
	if !ok {
		l = New(r.rate, r.burst)
		r.tenants[tenant] = l
	}
	return l
}

// TenantCount reports the number of registered tenants.
func (r *Registry) TenantCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.tenants)
}
