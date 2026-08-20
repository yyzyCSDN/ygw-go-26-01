package sampler

import (
	"sync"
	"time"

	"example.com/tracelink/internal/model"
)

// Window is a time-bounded tail sampling window.
type Window struct {
	mu    sync.Mutex
	span  time.Duration
	spans map[string][]model.Span
	seen  map[string]time.Time
}

// NewWindow creates a window that expires entries after the duration.
func NewWindow(span time.Duration) *Window {
	return &Window{
		span:  span,
		spans: make(map[string][]model.Span),
		seen:  make(map[string]time.Time),
	}
}

// Add records a span in the window.
func (w *Window) Add(span model.Span) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if _, ok := w.seen[span.TraceID]; !ok {
		w.seen[span.TraceID] = time.Now()
	}
	w.spans[span.TraceID] = append(w.spans[span.TraceID], span)
}

// Expired returns trace ids whose first span is older than the window.
func (w *Window) Expired(now time.Time) []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	out := make([]string, 0)
	for id, first := range w.seen {
		if now.Sub(first) > w.span {
			out = append(out, id)
		}
	}
	return out
}

// Remove drops a trace from the window.
func (w *Window) Remove(traceID string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.spans, traceID)
	delete(w.seen, traceID)
}

// Len returns the number of traces currently tracked.
func (w *Window) Len() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.seen)
}
