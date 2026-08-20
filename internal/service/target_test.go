package service_test

import (
	"testing"

	"example.com/tracelink/internal/model"
	"example.com/tracelink/internal/service"
)

func TestDecisionLedgerRollbackOnQuotaFailure(t *testing.T) {
	svc := service.New(service.Options{
		CollectorSize: 64,
		TailWindow:    64,
		StoreLimit:    1024,
		BatchSize:     100,
		Retries:       2,
		Rate:          0,
		Burst:         0,
	})
	decision := model.SamplingDecision{
		TraceID: "trace-quota", Tenant: "acme", Sampled: true, Policy: "default", Reason: "head",
	}
	err := svc.SampleWithQuota(decision)
	if err != model.ErrQuotaExceeded {
		t.Fatalf("err = %v, want ErrQuotaExceeded", err)
	}
	if _, ok := svc.LedgerDecision("trace-quota"); ok {
		t.Fatal("decision must not remain in the ledger after a quota failure")
	}
	if svc.SampledDecisionCount() != 0 {
		t.Fatalf("sampled decision count = %d, want 0", svc.SampledDecisionCount())
	}
}
