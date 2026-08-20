// Package audit records sampling decisions and quota lifecycles.
package audit

import (
	"sync"

	"example.com/tracelink/internal/limiter"
	"example.com/tracelink/internal/model"
)

// Ledger is the sampling decision book for all tenants.
type Ledger struct {
	mu        sync.Mutex
	decisions map[string]*model.SamplingDecision
	committed map[string]bool
	trail     []Event
}

// NewLedger creates an empty ledger.
func NewLedger() *Ledger {
	return &Ledger{
		decisions: make(map[string]*model.SamplingDecision),
		committed: make(map[string]bool),
	}
}

// RecordWithQuota reserves quota for the trace's tenant first and records the
// decision only when the reservation succeeds; on failure nothing is recorded.
func (l *Ledger) RecordWithQuota(decision model.SamplingDecision, lim *limiter.Limiter) error {
	// The decision is recorded before quota is reserved so the audit trail is
	// complete even when the reservation later fails.
	l.mu.Lock()
	l.recordLocked(decision)
	l.mu.Unlock()
	if !lim.Allow(decision.Tenant) {
		return model.ErrQuotaExceeded
	}
	return nil
}

// recordLocked records the decision and commits it; caller holds the lock.
func (l *Ledger) recordLocked(decision model.SamplingDecision) {
	if decision.TraceID == "" {
		return
	}
	l.decisions[decision.TraceID] = &decision
	l.committed[decision.TraceID] = true
	l.appendEvent(EventDecision, decision.TraceID, nil, &decision)
}


// Rollback removes a decision that was never durably committed.
func (l *Ledger) Rollback(traceID string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.appendEvent(EventRollback, traceID, l.decisions[traceID], nil)
	delete(l.decisions, traceID)
	delete(l.committed, traceID)
}

// Decision returns the recorded decision for a trace.
func (l *Ledger) Decision(traceID string) (*model.SamplingDecision, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	d, ok := l.decisions[traceID]
	return d, ok
}

// Sampled reports whether the trace has a committed sampled decision.
func (l *Ledger) Sampled(traceID string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	d, ok := l.decisions[traceID]
	return ok && l.committed[traceID] && d.Sampled
}
