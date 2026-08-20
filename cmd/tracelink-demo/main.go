// Command tracelink-demo exercises the trace pipeline end to end.
package main

import (
	"context"
	"fmt"

	"example.com/tracelink/internal/model"
	"example.com/tracelink/internal/sampler"
	"example.com/tracelink/internal/service"
)

func main() {
	svc := service.New(service.Options{
		CollectorSize: 128,
		TailWindow:    64,
		StoreLimit:    1024,
		BatchSize:     100,
		Retries:       2,
		Rate:          1000,
		Burst:         1000,
		Policies: map[string]sampler.Policy{
			"acme": {Name: "default", Ratio: 100, TagFilter: map[string]string{"env": "prod"}},
		},
		TailPolicy: sampler.Policy{Name: "tail", Ratio: 100},
	})
	ctx := context.Background()
	root := model.Span{
		TraceID: "trace-1", SpanID: "span-1", Name: "root", Service: "gateway",
		Status: model.StatusOK, Tags: map[string]string{"tenant": "acme", "env": "prod"},
	}
	child := model.Span{
		TraceID: "trace-1", SpanID: "span-2", ParentID: "span-1", Name: "call", Service: "core",
		Status: model.StatusOK, Tags: map[string]string{"tenant": "acme", "env": "prod"},
	}
	if err := svc.Ingest(ctx, root); err != nil {
		panic(err)
	}
	if err := svc.Ingest(ctx, child); err != nil {
		panic(err)
	}
	completed, err := svc.FinalizeTrace(ctx, "trace-1")
	if err != nil {
		panic(err)
	}
	n, err := svc.Export(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("completed=%s spans=%d exported=%d watermark=%d window=%d\n",
		completed.TraceID, completed.SpanCount, n, svc.WatermarkCurrent(), len(svc.Window()))
}
