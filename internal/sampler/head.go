package sampler

import (
	"sync"

	"example.com/tracelink/internal/model"
)

// HeadSampler makes an immediate sampling decision per span using tenant policy.
type HeadSampler struct {
	mu       sync.Mutex
	policies map[string]Policy
	revision int
	cache    *ShardedCache
}

// NewHead creates a head sampler with the given tenant policies.
func NewHead(policies map[string]Policy) *HeadSampler {
	return &HeadSampler{
		policies: policies,
		cache:    NewShardedCache(),
	}
}

// TenantFor extracts the tenant tag from a span.
func TenantFor(span model.Span) string {
	return span.Tags["tenant"]
}

// Decide returns the sampling decision for the span's tenant.
func (h *HeadSampler) Decide(span model.Span) (*model.SamplingDecision, error) {
	tenant := TenantFor(span)
	h.mu.Lock()
	policy, ok := h.policies[tenant]
	h.mu.Unlock()
	if !ok {
		return nil, model.ErrNoPolicy
	}
	sampled := policy.Match(span) && span.Status != model.StatusCancelled
	decision := model.NewDecision(span.TraceID, tenant, sampled, policy.Name, "head")
	h.cache.Put(span.TraceID, decision)
	return decision, nil
}

// DecisionForTrace returns the cached decision for the exact trace id.
func (h *HeadSampler) DecisionForTrace(traceID string) (*model.SamplingDecision, bool) {
	// Reads the shard slot directly, so a reused slot can answer for another
	// trace that hashes to the same shard.
	return h.cache.slot(traceID)
}

