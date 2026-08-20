package exporter

import (
	"context"
	"errors"
	"testing"

	"example.com/tracelink/internal/model"
)

// failingSink rejects every WriteBatch call with a sentinel error.
type failingSink struct{ err error }

func (f failingSink) WriteBatch([]model.ExportRecord) error { return f.err }

func TestFlushFailureKeepsRecordsPending(t *testing.T) {
	sinkErr := errors.New("sink down")
	b := NewBatch(failingSink{err: sinkErr}, 0, 2) // whole-batch, 2 retries
	b.Submit(model.ExportRecord{TraceID: "trace-1"})
	b.Submit(model.ExportRecord{TraceID: "trace-2"})

	err := b.Flush(context.Background())
	if !errors.Is(err, sinkErr) {
		t.Fatalf("expected sink error to be propagated, got %v", err)
	}
	// The records must remain queued so a later flush can retry them.
	if got := b.PendingCount(); got != 2 {
		t.Fatalf("expected 2 records to stay pending after failure, got %d", got)
	}
}

func TestFlushSuccessClearsPending(t *testing.T) {
	b := NewBatch(NewMemorySink(), 0, 1)
	b.Submit(model.ExportRecord{TraceID: "trace-1"})

	if err := b.Flush(context.Background()); err != nil {
		t.Fatalf("expected successful flush, got %v", err)
	}
	if got := b.PendingCount(); got != 0 {
		t.Fatalf("expected 0 pending after success, got %d", got)
	}
}

// TestFlushFailureRetryThenRecover verifies that records kept pending by a
// failed flush are exported once the sink recovers on a later flush.
func TestFlushFailureRetryThenRecover(t *testing.T) {
	flake := &flakySink{fail: true}
	b := NewBatch(flake, 0, 0) // no retries: first failure surfaces immediately
	b.Submit(model.ExportRecord{TraceID: "trace-1"})

	if err := b.Flush(context.Background()); err == nil {
		t.Fatal("expected first flush to fail")
	}
	if got := b.PendingCount(); got != 1 {
		t.Fatalf("expected 1 record pending after failure, got %d", got)
	}

	// Sink recovers; the previously failed record must now be exported.
	flake.fail = false
	if err := b.Flush(context.Background()); err != nil {
		t.Fatalf("expected recovery flush to succeed, got %v", err)
	}
	if got := b.PendingCount(); got != 0 {
		t.Fatalf("expected 0 pending after recovery, got %d", got)
	}
}

type flakySink struct {
	fail bool
}

func (f *flakySink) WriteBatch([]model.ExportRecord) error {
	if f.fail {
		return errors.New("sink down")
	}
	return nil
}
