package trace

import "strings"

// StoreByPrefix returns stored trace ids that start with the given prefix.
func (s *Store) StoreByPrefix(prefix string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]string, 0)
	for _, id := range s.order {
		if strings.HasPrefix(id, prefix) {
			out = append(out, id)
		}
	}
	return out
}
