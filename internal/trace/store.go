package trace

import (
	"sync"

	"example.com/tracelink/internal/model"
)

// Store keeps completed traces with a bounded retention.
type Store struct {
	mu         sync.Mutex
	limit      int
	traces     map[string]*model.CompletedTrace
	order      []string
	generation int
}

// NewStore creates a store retaining at most limit traces.
func NewStore(limit int) *Store {
	return &Store{
		limit:  limit,
		traces: make(map[string]*model.CompletedTrace),
	}
}

// Put stores a completed trace, evicting the oldest when over the limit.
func (s *Store) Put(completed *model.CompletedTrace) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.traces[completed.TraceID] = completed
	s.order = append(s.order, completed.TraceID)
	s.generation++
	for len(s.order) > s.limit {
		s.evictOldest()
	}
}

// Get returns a completed trace; evicted traces are unknown.
func (s *Store) Get(traceID string) (*model.CompletedTrace, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.traces[traceID]
	return t, ok
}

// Generation returns the number of completed puts so far.
func (s *Store) Generation() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.generation
}
