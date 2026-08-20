package model

import "fmt"

// BuildCompleted assembles a CompletedTrace from its invariants.
func BuildCompleted(traceID, rootSpanID string, spanCount int, status SpanStatus) *CompletedTrace {
	return &CompletedTrace{
		TraceID:    traceID,
		RootSpanID: rootSpanID,
		SpanCount:  spanCount,
		Status:     status,
	}
}

// ValidateCompleted checks the invariants of a completed trace.
func ValidateCompleted(t *CompletedTrace) error {
	if t == nil {
		return fmt.Errorf("%w: nil completed trace", ErrUnknownTrace)
	}
	if t.TraceID == "" {
		return fmt.Errorf("%w: empty trace id", ErrUnknownTrace)
	}
	if t.RootSpanID == "" {
		return fmt.Errorf("%w: empty root span id", ErrUnknownTrace)
	}
	if t.SpanCount <= 0 {
		return fmt.Errorf("%w: empty span count", ErrUnknownTrace)
	}
	return nil
}

// StatusLabel returns a stable human-readable status label.
func StatusLabel(s SpanStatus) string {
	return s.String()
}
