package service_test

import (
	"context"
	"testing"

	"example.com/tracelink/internal/model"
	"example.com/tracelink/internal/sampler"
	"example.com/tracelink/internal/service"
)

func TestHeadSamplerUnknownTenantNoPanic(t *testing.T) {
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
	span := model.Span{
		TraceID: "trace-ghost", SpanID: "s-1", Name: "call", Service: "core",
		Status: model.StatusOK, Tags: map[string]string{"tenant": "ghost", "env": "prod"},
	}
	err := svc.Ingest(context.Background(), span)
	if err != model.ErrNoPolicy {
		t.Fatalf("err = %v, want ErrNoPolicy", err)
	}
}
