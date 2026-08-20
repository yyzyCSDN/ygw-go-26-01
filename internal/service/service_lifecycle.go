package service

import "context"

// OpenTraces returns the number of traces currently being built.
func (s *Service) OpenTraces() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.trees)
}

// Close finalizes all open traces and flushes pending exports.
func (s *Service) Close(ctx context.Context) (int, error) {
	s.mu.Lock()
	open := make([]string, 0, len(s.trees))
	for id := range s.trees {
		open = append(open, id)
	}
	s.mu.Unlock()
	for _, id := range open {
		if _, err := s.FinalizeTrace(ctx, id); err != nil {
			return 0, err
		}
	}
	return s.Export(ctx)
}
