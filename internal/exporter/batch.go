package exporter

import (
	"context"
	"errors"
	"time"

	"example.com/tracelink/internal/model"
)

var errSinkExists = errors.New("sink name already registered")

// BatchExporter flushes traces to a sink with bounded retries.
type BatchExporter struct {
	sink      Sink
	batchSize int
	retries   int
	pending   map[string]model.ExportRecord
	observer  Observer
}

// NewBatch creates a batch exporter.
func NewBatch(sink Sink, batchSize int, retries int) *BatchExporter {
	return &BatchExporter{
		sink:      sink,
		batchSize: batchSize,
		retries:   retries,
		pending:   make(map[string]model.ExportRecord),
	}
}

// Submit queues a trace for export.
func (b *BatchExporter) Submit(record model.ExportRecord) {
	b.pending[record.TraceID] = record
}

// PendingCount reports how many records await a successful flush.
func (b *BatchExporter) PendingCount() int {
	return len(b.pending)
}

// Flush writes all pending records. On failure the records stay pending and
// the error is propagated so a later flush can retry them.
func (b *BatchExporter) Flush(ctx context.Context) error {
	if len(b.pending) == 0 {
		return nil
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	pending := make([]model.ExportRecord, 0, len(b.pending))
	for _, r := range b.pending {
		pending = append(pending, r)
	}
	var lastErr error
	for attempt := 0; attempt <= b.retries; attempt++ {
		ok := true
		step := b.batchSize
		if step <= 0 {
			step = len(pending)
		}
		for start := 0; start < len(pending); start += step {
			end := start + step
			if end > len(pending) {
				end = len(pending)
			}
			if err := b.sink.WriteBatch(pending[start:end]); err != nil {
				lastErr = err
				ok = false
				break
			}
		}
		if ok {
			b.pending = make(map[string]model.ExportRecord)
			b.notifySuccess(len(pending))
			return nil
		}
		time.Sleep(Backoff(attempt))
	}
	b.notifyFailure(lastErr)
	return lastErr
}
