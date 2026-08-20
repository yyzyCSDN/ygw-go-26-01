package model

// IsTerminal reports whether the status ends the span lifecycle.
func (s SpanStatus) IsTerminal() bool {
	return s == StatusOK || s == StatusError || s == StatusCancelled
}

// ValidTransition reports whether moving from the old status to the new one is
// allowed by the span lifecycle.
func ValidTransition(old, new SpanStatus) bool {
	if old.IsTerminal() {
		return old == new
	}
	switch new {
	case StatusUnset, StatusOK, StatusError, StatusCancelled:
		return true
	default:
		return false
	}
}
