package service_test

import (
	"context"
	"testing"

	"example.com/tracelink/internal/model"
	"example.com/tracelink/internal/sampler"
	"example.com/tracelink/internal/service"
)

func TestTreeBuffersOutOfOrderChild(t *testing.T) {
	svc := service.New(service.Options{
		CollectorSize: 256,
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
	ctx := context.Background()
	child := model.Span{
		TraceID: "trace-o3", SpanID: "s-child", ParentID: "s-root", Name: "call",
		Service: "core", Status: model.StatusOK,
		Tags: map[string]string{"tenant": "acme", "env": "prod"},
	}
	if err := svc.Ingest(ctx, child); err != nil {
		t.Fatal(err)
	}
	root := model.Span{
		TraceID: "trace-o3", SpanID: "s-root", Name: "root", Service: "gateway",
		Status: model.StatusOK, Tags: map[string]string{"tenant": "acme", "env": "prod"},
	}
	if err := svc.Ingest(ctx, root); err != nil {
		t.Fatal(err)
	}
	completed, err := svc.FinalizeTrace(ctx, "trace-o3")
	if err != nil {
		t.Fatal(err)
	}
	if completed.SpanCount != 2 {
		t.Fatalf("span count = %d, want 2 (out-of-order child must be buffered)", completed.SpanCount)
	}
}
