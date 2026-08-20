package collector

import (
	"context"
	"testing"

	"example.com/tracelink/internal/model"
)

// validSpan builds a minimal span that passes model.ValidateSpan.
func validSpan(name, id string) model.Span {
	return model.Span{TraceID: "t1", SpanID: id, Name: name, StartNs: 1, EndNs: 2}
}

// TestWindowReturnsIndependentCopy reproduces the reported aliasing bug: a
// snapshot taken via Window must not change when subsequent spans are ingested
// into the ring. Before the fix, Window() returned a slice over the live ring
// backing array, so a later Add() that recycled the same slot would overwrite
// an entry the caller already held.
func TestWindowReturnsIndependentCopy(t *testing.T) {
	// A window small enough to force ring wrap-around after a few inserts.
	c := NewCollector(2)

	if err := c.Add(context.Background(), validSpan("a", "1")); err != nil {
		t.Fatalf("add a: %v", err)
	}
	if err := c.Add(context.Background(), validSpan("b", "2")); err != nil {
		t.Fatalf("add b: %v", err)
	}

	snap := c.Window()
	if len(snap) != 2 || snap[0].Name != "a" || snap[1].Name != "b" {
		t.Fatalf("snapshot before wrap = %+v, want [a, b]", snap)
	}

	// Third write wraps the ring and overwrites slot 0 ("a") with "c".
	if err := c.Add(context.Background(), validSpan("c", "3")); err != nil {
		t.Fatalf("add c: %v", err)
	}

	// The snapshot taken earlier must still reflect the buffer as it was.
	if snap[0].Name != "a" || snap[1].Name != "b" {
		t.Fatalf("snapshot mutated by later ingest: %+v, want [a, b]", snap)
	}

	// Backing arrays must not be shared.
	if len(snap) > 0 && len(c.window) > 0 && &snap[0] == &c.window[0] {
		t.Fatalf("snapshot shares backing array with internal ring")
	}
}

// TestWindowRingOrder confirms spans are returned in insertion order even
// after the ring has wrapped, matching Drain()'s ordering.
func TestWindowRingOrder(t *testing.T) {
	c := NewCollector(3)
	for _, s := range []model.Span{
		validSpan("a", "1"),
		validSpan("b", "2"),
		validSpan("c", "3"),
		validSpan("d", "4"), // wraps, evicting "a"
	} {
		if err := c.Add(context.Background(), s); err != nil {
			t.Fatalf("add %s: %v", s.Name, err)
		}
	}

	got := c.Window()
	want := []string{"b", "c", "d"}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i, w := range want {
		if got[i].Name != w {
			t.Fatalf("got[%d].Name = %q, want %q", i, got[i].Name, w)
		}
	}

	// Drain must agree with the snapshot's contents/order.
	drained := c.Drain()
	if len(drained) != len(got) {
		t.Fatalf("drain len = %d, want %d", len(drained), len(got))
	}
	for i, s := range got {
		if drained[i].SpanID != s.SpanID {
			t.Fatalf("drain[%d] = %q, snapshot[%d] = %q", i, drained[i].SpanID, i, s.SpanID)
		}
	}
}

// TestWindowEmpty returns a non-nil empty slice for an empty window and never
// exposes the internal array.
func TestWindowEmpty(t *testing.T) {
	c := NewCollector(4)
	snap := c.Window()
	if len(snap) != 0 {
		t.Fatalf("len = %d, want 0", len(snap))
	}
}
