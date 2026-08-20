// Package exporter batches completed traces to a sink.
package exporter

import (
	"sync"

	"example.com/tracelink/internal/model"
)

// Sink receives export records.
type Sink interface {
	WriteBatch(records []model.ExportRecord) error
}

// MemorySink records exports in memory for inspection.
type MemorySink struct {
	mu      sync.Mutex
	records []model.ExportRecord
}

// NewMemorySink creates an empty memory sink.
func NewMemorySink() *MemorySink {
	return &MemorySink{}
}

// WriteBatch appends the records to the sink.
func (m *MemorySink) WriteBatch(records []model.ExportRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.records = append(m.records, records...)
	return nil
}

// Records returns a copy of the exported records.
func (m *MemorySink) Records() []model.ExportRecord {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]model.ExportRecord, len(m.records))
	copy(out, m.records)
	return out
}
