package service_test

import (
	"context"
	"testing"
	"time"

	"example.com/tracelink/internal/model"
	"example.com/tracelink/internal/sampler"
	"example.com/tracelink/internal/service"
)

func TestTailFinalizeNoDeadlock(t *testing.T) {
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
	for i := 0; i < 2; i++ {
		span := model.Span{
			TraceID: "trace-tail", SpanID: string(rune('a' + i)), Name: "call",
			Service: "core", Status: model.StatusOK,
			Tags: map[string]string{"tenant": "acme", "env": "prod"},
		}
		if err := svc.Ingest(ctx, span); err != nil {
			t.Fatal(err)
		}
	}
	done := make(chan struct{})
	go func() {
		_, _ = svc.FinalizeTrace(ctx, "trace-tail")
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("FinalizeTrace deadlocked while finalizing a multi-span trace")
	}
}
