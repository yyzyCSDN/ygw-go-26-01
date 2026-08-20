package service

import (
	"context"

	"example.com/tracelink/internal/exporter"
	"example.com/tracelink/internal/model"
)

// ReplayTrace re-exports a previously completed trace from the store.
func (s *Service) ReplayTrace(ctx context.Context, traceID string) error {
	t, ok := s.store.Get(traceID)
	if !ok {
		return model.ErrUnknownTrace
	}
	s.exporter.Submit(exporter.BuildRecord(t, s.ledger.Sampled(traceID)))
	_, err := s.Export(ctx)
	return err
}

// ReplayAll re-exports every stored trace.
func (s *Service) ReplayAll(ctx context.Context) (int, error) {
	ids := s.store.IDs()
	for _, id := range ids {
		if err := s.ReplayTrace(ctx, id); err != nil {
			return 0, err
		}
	}
	return len(ids), nil
}
