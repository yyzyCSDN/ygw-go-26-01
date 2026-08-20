package trace

import (
	"sync"

	"example.com/tracelink/internal/model"
)

// Tree builds the parent-child span tree of one trace.
type Tree struct {
	mu       sync.Mutex
	root     string
	parents  map[string]string
	children map[string][]string
	status   map[string]model.SpanStatus
	pending  map[string][]string
}

// NewTree creates an empty correlation tree.
func NewTree() *Tree {
	return &Tree{
		parents:  make(map[string]string),
		children: make(map[string][]string),
		status:   make(map[string]model.SpanStatus),
		pending:  make(map[string][]string),
	}
}

// Attach links a span into the tree. Children whose parent has not arrived
// yet are buffered and linked once the parent shows up.
func (t *Tree) Attach(span model.Span) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.status[span.SpanID] = span.Status
	if span.ParentID == "" {
		if t.root == "" {
			t.root = span.SpanID
		}
	} else {
		t.parents[span.SpanID] = span.ParentID
		if _, ok := t.status[span.ParentID]; ok {
			t.children[span.ParentID] = append(t.children[span.ParentID], span.SpanID)
		} else {
			t.pending[span.ParentID] = append(t.pending[span.ParentID], span.SpanID)
		}
	}
	for _, child := range t.pending[span.SpanID] {
		t.children[span.SpanID] = append(t.children[span.SpanID], child)
	}
	delete(t.pending, span.SpanID)
}

// RootSpanID returns the tree root, or the first span when no root arrived.
func (t *Tree) RootSpanID() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.root != "" {
		return t.root
	}
	for id := range t.status {
		return id
	}
	return ""
}

// SpanCount returns the number of directly attached spans.
func (t *Tree) SpanCount() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.status)
}

// TotalSpans includes spans buffered while waiting for their parent.
func (t *Tree) TotalSpans() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	total := len(t.status)
	for _, list := range t.pending {
		total += len(list)
	}
	return total
}

// Status returns the worst status observed in the tree.
func (t *Tree) Status() model.SpanStatus {
	t.mu.Lock()
	defer t.mu.Unlock()
	for id := range t.status {
		if t.status[id] == model.StatusError {
			return model.StatusError
		}
	}
	return model.StatusOK
}
