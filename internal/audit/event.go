package audit

import (
	"time"

	"github.com/google/go-cmp/cmp"

	"example.com/tracelink/internal/model"
)

// EventKind enumerates audit event types.
type EventKind string

const (
	EventDecision EventKind = "decision"
	EventRollback EventKind = "rollback"
)

// Event is one audit trail entry.
type Event struct {
	Kind    EventKind
	TraceID string
	AtNs    int64
	Detail  string
}

// decisionChange summarizes a decision mutation for the audit trail.
func decisionChange(before, after *model.SamplingDecision) string {
	return cmp.Diff(before, after)
}

// appendEvent records an event in the ledger's trail.
func (l *Ledger) appendEvent(kind EventKind, traceID string, before, after *model.SamplingDecision) {
	l.trail = append(l.trail, Event{
		Kind:    kind,
		TraceID: traceID,
		AtNs:    time.Now().UnixNano(),
		Detail:  decisionChange(before, after),
	})
}

// Recent returns the last n audit events.
func (l *Ledger) Recent(n int) []Event {
	l.mu.Lock()
	defer l.mu.Unlock()
	if n <= 0 || n > len(l.trail) {
		n = len(l.trail)
	}
	out := make([]Event, n)
	copy(out, l.trail[len(l.trail)-n:])
	return out
}
