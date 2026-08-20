package service_test

import (
	"context"
	"testing"

	"example.com/tracelink/internal/model"
	"example.com/tracelink/internal/sampler"
	"example.com/tracelink/internal/service"
)

type failSink struct {
	err error
}

func (f failSink) WriteBatch(_ []model.ExportRecord) error {
	return f.err
}

func TestExporterFlushPropagatesError(t *testing.T) {
	svc := service.New(service.Options{
		CollectorSize: 256,
		TailWindow:    64,
		StoreLimit:    1024,
		BatchSize:     100,
		Retries:       0,
		Rate:          10000,
		Burst:         10000,
		Sink:          failSink{err: model.ErrSinkFailed},
		Policies: map[string]sampler.Policy{
			"acme": {Name: "default", Ratio: 100, TagFilter: map[string]string{"env": "prod"}},
		},
		TailPolicy: sampler.Policy{Name: "tail", Ratio: 100},
	})
	ctx := context.Background()
	root := model.Span{
		TraceID: "trace-export", SpanID: "root", Name: "root", Service: "gateway",
		Status: model.StatusOK, Tags: map[string]string{"tenant": "acme", "env": "prod"},
	}
	if err := svc.Ingest(ctx, root); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.FinalizeTrace(ctx, "trace-export"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Export(ctx); err == nil {
		t.Fatal("expected export error when the sink fails")
	}
	if svc.WatermarkCurrent() != 0 {
		t.Fatalf("watermark advanced despite failure: %d", svc.WatermarkCurrent())
	}
	if svc.PendingExports() != 1 {
		t.Fatalf("pending exports = %d, want 1 (failed batch must stay retryable)", svc.PendingExports())
	}
}
