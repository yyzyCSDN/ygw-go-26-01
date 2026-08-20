package trace

// Count returns the number of stored traces.
func (s *Store) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.traces)
}

// IDs returns a copy of the stored trace ids in insertion order.
func (s *Store) IDs() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]string, len(s.order))
	copy(out, s.order)
	return out
}
