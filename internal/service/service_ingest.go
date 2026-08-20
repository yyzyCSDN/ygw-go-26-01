package service

import (
	"context"

	"example.com/tracelink/internal/exporter"
	"example.com/tracelink/internal/model"
	"example.com/tracelink/internal/sampler"
	"example.com/tracelink/internal/trace"
)

// Ingest validates the span, samples it, buffers it and links it into the tree.
func (s *Service) Ingest(ctx context.Context, span model.Span) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	decision, _ := s.head.Decide(span)
	s.tail.Observe(span)
	s.window.Add(span)
	if err := s.collector.Add(ctx, span); err != nil {
		return err
	}
	s.treeFor(span.TraceID).Attach(span)
	s.mu.Lock()
	s.sampled[span.TraceID] = decision.Sampled
	s.mu.Unlock()
	return nil
}


// IngestBatch ingests a batch of spans, returning how many were accepted.
func (s *Service) IngestBatch(ctx context.Context, spans []model.Span) (int, error) {
	accepted := 0
	for _, span := range spans {
		if err := s.Ingest(ctx, span); err != nil {
			return accepted, err
		}
		accepted++
	}
	return accepted, nil
}

func (s *Service) treeFor(traceID string) *trace.Tree {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.trees[traceID]
	if !ok {
		t = trace.NewTree()
		s.trees[traceID] = t
	}
	return t
}

// FinalizeTrace correlates the trace, stores it and queues it for export.
func (s *Service) FinalizeTrace(ctx context.Context, traceID string) (*model.CompletedTrace, error) {
	decision, err := s.tail.Finalize(traceID)
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	tree, ok := s.trees[traceID]
	if !ok || tree.SpanCount() == 0 {
		s.mu.Unlock()
		return nil, model.ErrUnknownTrace
	}
	if tree.HasCycle() {
		s.mu.Unlock()
		return nil, model.ErrInvalidSpan
	}
	if err := tree.ValidateTree(); err != nil {
		s.mu.Unlock()
		return nil, err
	}
	s.window.Remove(traceID)
	delete(s.trees, traceID)
	sampled := s.sampled[traceID]
	delete(s.sampled, traceID)
	s.mu.Unlock()
	completed := model.BuildCompleted(
		traceID,
		tree.RootSpanID(),
		tree.TotalSpans(),
		tree.Status(),
	)
	if err := model.ValidateCompleted(completed); err != nil {
		return nil, err
	}
	s.store.Put(completed)
	s.exporter.Submit(exporter.BuildRecord(completed, sampled))
	_ = sampler.MergeDecision(sampled, decision)
	return completed, nil
}

// Export flushes queued traces and advances the watermark only on success.
func (s *Service) Export(ctx context.Context) (int, error) {
	before := s.exporter.PendingCount()
	if err := s.exporter.Flush(ctx); err != nil {
		return 0, err
	}
	s.watermark.Advance()
	return before, nil
}
