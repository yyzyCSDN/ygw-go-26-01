package sampler

import (
	"github.com/google/go-cmp/cmp"

	"example.com/tracelink/internal/model"
)

// Policy is a per-tenant sampling rule.
type Policy struct {
	Name      string
	Ratio     int // 0-100; 100 keeps everything
	TagFilter map[string]string
}

// Match reports whether the span satisfies the policy tag filter.
func (p Policy) Match(span model.Span) bool {
	for k, v := range p.TagFilter {
		if span.Tags[k] != v {
			return false
		}
	}
	return true
}

// Equal reports whether two policies are identical.
func (p Policy) Equal(other Policy) bool {
	return cmp.Equal(p, other)
}
