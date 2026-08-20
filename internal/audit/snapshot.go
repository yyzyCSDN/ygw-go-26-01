package audit

import "example.com/tracelink/internal/model"

// Snapshot is an immutable view of the decision ledger.
type Snapshot struct {
	Decisions map[string]*model.SamplingDecision
	Committed map[string]bool
}

// Snapshot copies the current ledger state.
func (l *Ledger) Snapshot() Snapshot {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := Snapshot{
		Decisions: make(map[string]*model.SamplingDecision, len(l.decisions)),
		Committed: make(map[string]bool, len(l.committed)),
	}
	for id, d := range l.decisions {
		copyDecision := *d
		out.Decisions[id] = &copyDecision
	}
	for id, v := range l.committed {
		out.Committed[id] = v
	}
	return out
}

// Consistent reports whether every committed decision has a recorded entry.
func (s Snapshot) Consistent() bool {
	for id := range s.Committed {
		if _, ok := s.Decisions[id]; !ok {
			return false
		}
	}
	return true
}
