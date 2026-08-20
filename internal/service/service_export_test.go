package service

import (
	"context"
	"errors"
	"testing"

	"example.com/tracelink/internal/model"
	"example.com/tracelink/internal/sampler"
)

// failingSink rejects every WriteBatch with a sentinel error.
type failingSink struct{ err error }

func (f failingSink) WriteBatch([]model.ExportRecord) error { return f.err }

// newServiceWithSink builds a service wired to the supplied sink.
func newServiceWithSink(t *testing.T, sink exporterSink) *Service {
	t.Helper()
	return New(Options{
		CollectorSize: 128,
		TailWindow:    64,
		StoreLimit:    1024,
		BatchSize:     100,
		Retries:       1,
		Rate:          1000,
		Burst:         1000,
		Policies: map[string]sampler.Policy{
			"acme": {Name: "default", Ratio: 100, TagFilter: map[string]string{"env": "prod"}},
		},
		TailPolicy: sampler.Policy{Name: "tail", Ratio: 100},
		Sink:       sink,
	})
}

// exporterSink is the subset of exporter.Sink used by tests.
type exporterSink interface {
	WriteBatch([]model.ExportRecord) error
}

func TestExportHoldsWatermarkOnSinkFailure(t *testing.T) {
	sinkErr := errors.New("sink down")
	svc := newServiceWithSink(t, failingSink{err: sinkErr})
	ctx := context.Background()

	// Queue a trace for export.
	if err := svc.Ingest(ctx, model.Span{
		TraceID: "trace-1", SpanID: "span-1", Name: "root", Service: "gateway",
		Status: model.StatusOK, Tags: map[string]string{"tenant": "acme", "env": "prod"},
	}); err != nil {
		t.Fatalf("ingest: %v", err)
	}
	if _, err := svc.FinalizeTrace(ctx, "trace-1"); err != nil {
		t.Fatalf("finalize: %v", err)
	}

	wmBefore := svc.WatermarkCurrent()
	pendingBefore := svc.PendingExports()

	n, err := svc.Export(ctx)
	if !errors.Is(err, sinkErr) {
		t.Fatalf("expected sink error to surface from Export, got (%d, %v)", n, err)
	}
	// The failed batch must remain queued for a later retry.
	if got := svc.PendingExports(); got != pendingBefore {
		t.Fatalf("expected %d records to stay pending after failure, got %d", pendingBefore, got)
	}
	// The watermark must not move until the batch is truly written.
	if got := svc.WatermarkCurrent(); got != wmBefore {
		t.Fatalf("watermark advanced on failure: %d -> %d", wmBefore, got)
	}
}

func TestExportAdvancesWatermarkOnSuccess(t *testing.T) {
	svc := newServiceWithSink(t, nil) // default MemorySink never fails
	ctx := context.Background()

	if err := svc.Ingest(ctx, model.Span{
		TraceID: "trace-1", SpanID: "span-1", Name: "root", Service: "gateway",
		Status: model.StatusOK, Tags: map[string]string{"tenant": "acme", "env": "prod"},
	}); err != nil {
		t.Fatalf("ingest: %v", err)
	}
	if _, err := svc.FinalizeTrace(ctx, "trace-1"); err != nil {
		t.Fatalf("finalize: %v", err)
	}

	wmBefore := svc.WatermarkCurrent()
	if _, err := svc.Export(ctx); err != nil {
		t.Fatalf("export: %v", err)
	}
	if got := svc.WatermarkCurrent(); got != wmBefore+1 {
		t.Fatalf("expected watermark to advance by 1, got %d -> %d", wmBefore, got)
	}
	if got := svc.PendingExports(); got != 0 {
		t.Fatalf("expected 0 pending after success, got %d", got)
	}
}
