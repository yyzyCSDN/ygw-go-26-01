package model

// WithTag returns a copy of the span with one tag applied.
func (s Span) WithTag(key, value string) Span {
	out := s
	out.Tags = CloneTags(s.Tags)
	out.Tags[key] = value
	return out
}

// CloneTags returns a defensive copy of the tag map.
func CloneTags(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

// DurationNs returns the span duration in nanoseconds.
func (s Span) DurationNs() int64 {
	if s.EndNs <= s.StartNs {
		return 0
	}
	return s.EndNs - s.StartNs
}

// Finished reports whether the span has a terminal status.
func (s Span) Finished() bool {
	return s.Status == StatusOK || s.Status == StatusError || s.Status == StatusCancelled
}
