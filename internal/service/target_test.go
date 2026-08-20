package service_test

import (
	"context"
	"fmt"
	"testing"

	"example.com/tracelink/internal/model"
	"example.com/tracelink/internal/sampler"
	"example.com/tracelink/internal/service"
)

func TestSamplerShardDecisionOwnership(t *testing.T) {
	svc := service.New(service.Options{
		CollectorSize: 1024,
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
	const traces = 32
	for i := 0; i < traces; i++ {
		tid := fmt.Sprintf("trace-shard-%02d", i)
		span := model.Span{
			TraceID: tid, SpanID: fmt.Sprintf("s-%02d", i), Name: "call", Service: "core",
			Status: model.StatusOK, Tags: map[string]string{"tenant": "acme", "env": "prod"},
		}
		if err := svc.Ingest(ctx, span); err != nil {
			t.Fatal(err)
		}
	}
	for i := 0; i < traces; i++ {
		tid := fmt.Sprintf("trace-shard-%02d", i)
		decision, ok := svc.DecisionForTrace(tid)
		if !ok {
			t.Fatalf("%s: decision missing", tid)
		}
		if decision.TraceID != tid {
			t.Fatalf("%s: returned decision belongs to %s (shard slot reused)", tid, decision.TraceID)
		}
	}
}
