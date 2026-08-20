package service

// PendingExports returns the number of queued export records.
func (s *Service) PendingExports() int {
	return s.exporter.PendingCount()
}

// SampledDecisionCount returns the number of committed sampled decisions.
func (s *Service) SampledDecisionCount() int {
	return s.ledger.SampledCount()
}

// StoredTraceCount returns the number of completed traces in the store.
func (s *Service) StoredTraceCount() int {
	return s.store.Count()
}
