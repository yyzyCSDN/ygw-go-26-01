package trace

// evictOldest removes the oldest inserted trace id and returns it.
// The caller must hold the store lock.
func (s *Store) evictOldest() string {
	old := s.order[0]
	s.order = s.order[1:]
	delete(s.traces, old)
	return old
}
