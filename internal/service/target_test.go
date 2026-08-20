package service_test

import (
	"context"
	"fmt"
	"testing"

	"example.com/tracelink/internal/model"
	"example.com/tracelink/internal/sampler"
	"example.com/tracelink/internal/service"
)

func TestCollectorWindowBatchNoAlias(t *testing.T) {
	svc := service.New(service.Options{
		CollectorSize: 4,
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
	for i := 0; i < 4; i++ {
		span := model.Span{
			TraceID: "trace-window", SpanID: fmt.Sprintf("s-%02d", i), Name: "call",
			Service: "core", Status: model.StatusOK,
			Tags: map[string]string{"tenant": "acme", "env": "prod"},
		}
		if err := svc.Ingest(ctx, span); err != nil {
			t.Fatal(err)
		}
	}
	first := svc.Window()
	if len(first) != 4 {
		t.Fatalf("window size = %d, want 4", len(first))
	}
	for i := 0; i < 4; i++ {
		span := model.Span{
			TraceID: "trace-overflow", SpanID: fmt.Sprintf("o-%02d", i), Name: "call",
			Service: "core", Status: model.StatusOK,
			Tags: map[string]string{"tenant": "acme", "env": "prod"},
		}
		if err := svc.Ingest(ctx, span); err != nil {
			t.Fatal(err)
		}
	}
	if first[0].TraceID != "trace-window" || first[0].SpanID != "s-00" {
		t.Fatalf("previously returned window entry was mutated: %+v", first[0])
	}
}
