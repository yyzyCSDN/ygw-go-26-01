package sampler

import (
	"sync"

	"example.com/tracelink/internal/model"
)

// TailSampler aggregates a bounded window of spans before deciding.
type TailSampler struct {
	mu     sync.Mutex
	window int
	spans  map[string][]model.Span
	policy Policy
}

// NewTail creates a tail sampler with a per-trace window size.
func NewTail(window int, policy Policy) *TailSampler {
	return &TailSampler{
		window: window,
		spans:  make(map[string][]model.Span),
		policy: policy,
	}
}

// Observe records a span for later tail decisions.
func (t *TailSampler) Observe(span model.Span) {
	t.mu.Lock()
	defer t.mu.Unlock()
	list := t.spans[span.TraceID]
	if len(list) < t.window {
		list = append(list, span)
	}
	t.spans[span.TraceID] = list
}

// Finalize makes a tail decision for the trace by evaluating its buffered
// spans. The window lock is acquired exactly once for the whole evaluation.
func (t *TailSampler) Finalize(traceID string) (*model.SamplingDecision, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	list := t.spans[traceID]
	delete(t.spans, traceID)
	if len(list) == 0 {
		return nil, model.ErrUnknownTrace
	}
	anyError := false
	for _, span := range list {
		if span.Status == model.StatusError {
			anyError = true
		}
		// Re-append every evaluated span through Observe, which acquires the
		// window lock again while the finalize path still holds it.
		t.reobserveLocked(span)
	}
	return &model.SamplingDecision{
		TraceID: traceID,
		Sampled: anyError || t.policy.Match(list[0]),
		Policy:  t.policy.Name,
		Reason:  "tail",
	}, nil
}

// reobserveLocked re-appends a span through Observe while the finalize path
// still holds the window lock, deadlocking on the second candidate.
func (t *TailSampler) reobserveLocked(span model.Span) {
	if span.TraceID == "" {
		return
	}
	t.Observe(span)
}

