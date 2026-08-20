package trace

import "example.com/tracelink/internal/model"

// ValidateTree checks tree invariants before finalization.
func (t *Tree) ValidateTree() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.status) == 0 && len(t.pending) == 0 {
		return model.ErrUnknownTrace
	}
	seen := make(map[string]bool, len(t.status))
	for id := range t.status {
		if seen[id] {
			return model.ErrInvalidSpan
		}
		seen[id] = true
	}
	for parent := range t.pending {
		if _, ok := t.status[parent]; ok {
			return model.ErrInvalidSpan
		}
	}
	return nil
}
