package collector

import (
	"context"
	"errors"
	"testing"
	"time"

	"example.com/tracelink/internal/model"
)

func validSpan() model.Span {
	return model.Span{
		TraceID: "trace-1", SpanID: "span-1", Name: "root", Service: "gateway",
		Status: model.StatusOK,
	}
}

// A span submitted under an already-cancelled context must fail cleanly: no
// window write, a cancellation error returned, and rejected (not accepted).
func TestAddCancelledContextFailsCleanly(t *testing.T) {
	c := NewCollector(8)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // caller cancels before initiating collection

	err := c.Add(ctx, validSpan())

	if err == nil {
		t.Fatal("expected cancellation error, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if got := len(c.Window()); got != 0 {
		t.Fatalf("window must stay empty after cancelled Add, got %d spans", got)
	}
	m := c.Metrics()
	if m.Accepted != 0 {
		t.Fatalf("accepted must stay 0 after cancelled Add, got %d", m.Accepted)
	}
	if m.Rejected != 1 {
		t.Fatalf("rejected must be 1 after cancelled Add, got %d", m.Rejected)
	}
}

// DeadlineExceeded is also a cancellation signal and must surface the same way.
func TestAddDeadlineExceededFailsCleanly(t *testing.T) {
	c := NewCollector(8)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	<-ctx.Done() // force the deadline to elapse

	err := c.Add(ctx, validSpan())
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context.DeadlineExceeded, got %v", err)
	}
	if got := len(c.Window()); got != 0 {
		t.Fatalf("window must stay empty, got %d spans", got)
	}
}

// An active context still collects normally.
func TestAddActiveContextCollects(t *testing.T) {
	c := NewCollector(8)

	if err := c.Add(context.Background(), validSpan()); err != nil {
		t.Fatalf("active Add must succeed, got %v", err)
	}
	if got := len(c.Window()); got != 1 {
		t.Fatalf("window must hold the span, got %d spans", got)
	}
	m := c.Metrics()
	if m.Accepted != 1 || m.Rejected != 0 {
		t.Fatalf("metrics mismatch: accepted=%d rejected=%d", m.Accepted, m.Rejected)
	}
}

// Cancelled Add short-circuits before validation: an invalid span under a
// cancelled context yields the cancellation error, not the validation error.
func TestAddCancelledShortCircuitsBeforeValidation(t *testing.T) {
	c := NewCollector(8)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	bad := validSpan()
	bad.Name = "" // would normally fail ValidateSpan

	err := c.Add(ctx, bad)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation must win over validation, got %v", err)
	}
	if got := len(c.Window()); got != 0 {
		t.Fatalf("window must stay empty, got %d spans", got)
	}
}
