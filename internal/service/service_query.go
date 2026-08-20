package service

import (
	"example.com/tracelink/internal/collector"
	"example.com/tracelink/internal/audit"
	"example.com/tracelink/internal/model"
	"example.com/tracelink/internal/propagate"
)

// Window returns the collector's pending span window.
func (s *Service) Window() []model.Span {
	return s.collector.Window()
}

// WindowStats returns collector occupancy metrics.
func (s *Service) WindowStats() collector.Stats {
	return s.collector.Stats()
}

// ExportedRecords returns the records written to the memory sink.
func (s *Service) ExportedRecords() []model.ExportRecord {
	return s.sink.Records()
}

// WatermarkCurrent returns the current export watermark generation.
func (s *Service) WatermarkCurrent() int {
	return s.watermark.Current()
}

// TraceSpanCount returns the stored span count for a completed trace.
func (s *Service) TraceSpanCount(traceID string) int {
	t, ok := s.store.Get(traceID)
	if !ok {
		return -1
	}
	return t.SpanCount
}

// LedgerDecision returns the recorded decision for a trace.
func (s *Service) LedgerDecision(traceID string) (*model.SamplingDecision, bool) {
	return s.ledger.Decision(traceID)
}

// DecisionForTrace returns the head decision cached for a trace.
func (s *Service) DecisionForTrace(traceID string) (*model.SamplingDecision, bool) {
	return s.head.DecisionForTrace(traceID)
}

// BaggageMerge validates and merges a baggage header into dst.
func (s *Service) BaggageMerge(dst model.Baggage, header string) error {
	if err := propagate.ValidateBaggage(dst); err != nil {
		return err
	}
	return propagate.MergeBaggage(dst, header)
}

// TraceHeaderDecode parses a trace header into a trace context.
func (s *Service) TraceHeaderDecode(header string) (model.TraceContext, error) {
	return propagate.DecodeTraceHeader(header)
}

// SamplingWindowLen returns the number of traces tracked by the tail window.
func (s *Service) SamplingWindowLen() int {
	return s.window.Len()
}

// ExportMetrics returns the exporter counters.
func (s *Service) ExportMetrics() (records int64, failures int64, flushes int64) {
	return s.metrics.Snapshot()
}

// LedgerReport returns per-tenant ledger statistics.
func (s *Service) LedgerReport() []audit.Report {
	return s.ledger.ReportByTenant()
}

// CollectorMetrics returns the collector acceptance counters.
func (s *Service) CollectorMetrics() collector.Metrics {
	return s.collector.Metrics()
}
