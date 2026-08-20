package service_test

import (
	"context"
	"testing"

	"example.com/tracelink/internal/model"
	"example.com/tracelink/internal/sampler"
	"example.com/tracelink/internal/service"
)

func TestIngestHonorsCancellation(t *testing.T) {
	svc := service.New(service.Options{
		CollectorSize: 64,
		TailWindow:    64,
		StoreLimit:    1024,
		BatchSize:     100,
		Retries:       2,
		Rate:          10000,
		Burst:         10000,
		Policies: map[string]sampler.Policy{
			"acme": {Name: "default", Ratio: 100, TagFilter: map[string]string{"env": "prod"}},
		},
		TailPolicy: sampler.Policy{Name: "tail", Ratio: 100},
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	span := model.Span{
		TraceID: "trace-cancel", SpanID: "s-1", Name: "call", Service: "core",
		Status: model.StatusOK, Tags: map[string]string{"tenant": "acme", "env": "prod"},
	}
	if err := svc.Ingest(ctx, span); err == nil {
		t.Fatal("expected an error when the caller context is cancelled")
	}
	if len(svc.Window()) != 0 {
		t.Fatalf("cancelled span was still ingested: window=%d", len(svc.Window()))
	}
}
