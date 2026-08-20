package propagate

import (
	"regexp"

	"example.com/tracelink/internal/model"
)

var baggageKeyPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)

// ValidBaggageKey reports whether a key is acceptable in a baggage header.
func ValidBaggageKey(key string) bool {
	return baggageKeyPattern.MatchString(key)
}

// ValidateBaggage reports whether the baggage map respects key rules and size.
func ValidateBaggage(b model.Baggage) error {
	if len(b) > maxBaggageEntries {
		return model.ErrInvalidBaggage
	}
	for key := range b {
		if !ValidBaggageKey(key) {
			return model.ErrInvalidBaggage
		}
	}
	return nil
}
