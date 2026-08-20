// Package model defines the core span, trace and sampling types.
package model

import (
	"errors"
	"fmt"
	"regexp"
)

// SpanStatus describes the terminal state of a span.
type SpanStatus int

const (
	StatusUnset SpanStatus = iota
	StatusOK
	StatusError
	StatusCancelled
)

// String returns a stable status label.
func (s SpanStatus) String() string {
	switch s {
	case StatusOK:
		return "ok"
	case StatusError:
		return "error"
	case StatusCancelled:
		return "cancelled"
	default:
		return "unset"
	}
}

// Span is a single unit of work recorded by the pipeline.
type Span struct {
	TraceID  string
	SpanID   string
	ParentID string
	Name     string
	Service  string
	Status   SpanStatus
	StartNs  int64
	EndNs    int64
	Tags     map[string]string
}

// TraceContext carries the active trace identity across service boundaries.
type TraceContext struct {
	TraceID string
	SpanID  string
	Sampled bool
}

// Valid reports whether both identifiers are present.
func (c TraceContext) Valid() bool {
	return c.TraceID != "" && c.SpanID != ""
}

// Baggage is a set of tenant-controlled key/value pairs propagated with a trace.
type Baggage map[string]string

// SamplingDecision is the outcome of head or tail sampling for one trace.
type SamplingDecision struct {
	TraceID string
	Tenant  string
	Sampled bool
	Policy  string
	Reason  string
}

// ExportRecord is one trace's export summary written to a sink.
type ExportRecord struct {
	TraceID      string
	SpanCount    int
	Sampled      bool
	ExportedAtNs int64
}

// CompletedTrace is a fully correlated trace ready for export.
type CompletedTrace struct {
	TraceID    string
	RootSpanID string
	SpanCount  int
	Status     SpanStatus
}

var (
	ErrInvalidSpan    = errors.New("invalid span")
	ErrUnknownTrace   = errors.New("unknown trace")
	ErrNoPolicy       = errors.New("no sampling policy for tenant")
	ErrQuotaExceeded  = errors.New("quota exceeded")
	ErrSinkFailed     = errors.New("sink write failed")
	ErrInvalidBaggage = errors.New("invalid baggage header")
)

var idPattern = regexp.MustCompile(`^[A-Za-z0-9._-]{1,64}$`)

// ValidateSpan checks the required span fields.
func ValidateSpan(s Span) error {
	if !idPattern.MatchString(s.TraceID) {
		return fmt.Errorf("%w: trace id %q", ErrInvalidSpan, s.TraceID)
	}
	if !idPattern.MatchString(s.SpanID) {
		return fmt.Errorf("%w: span id %q", ErrInvalidSpan, s.SpanID)
	}
	if s.ParentID != "" && !idPattern.MatchString(s.ParentID) {
		return fmt.Errorf("%w: parent id %q", ErrInvalidSpan, s.ParentID)
	}
	if s.Name == "" {
		return fmt.Errorf("%w: empty span name", ErrInvalidSpan)
	}
	return nil
}

// CloneBaggage returns a defensive copy of the baggage map.
func CloneBaggage(in Baggage) Baggage {
	out := make(Baggage, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
