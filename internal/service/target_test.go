package service_test

import (
	"fmt"
	"strings"
	"testing"

	"example.com/tracelink/internal/model"
	"example.com/tracelink/internal/sampler"
	"example.com/tracelink/internal/service"
)

func TestBaggageMergeNoPartialMutation(t *testing.T) {
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
	dst := model.Baggage{"keep": "me"}
	parts := make([]string, 0, 20)
	for i := 0; i < 20; i++ {
		parts = append(parts, fmt.Sprintf("k%02d=v", i))
	}
	err := svc.BaggageMerge(dst, strings.Join(parts, ","))
	if err != model.ErrInvalidBaggage {
		t.Fatalf("err = %v, want ErrInvalidBaggage", err)
	}
	if len(dst) != 1 || dst["keep"] != "me" {
		t.Fatalf("dst was partially mutated by an invalid header: %v", dst)
	}
}
