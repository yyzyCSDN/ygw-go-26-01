package limiter

import "math"

// BurstFor computes a burst size from a rate and a window in seconds.
func BurstFor(rate float64, seconds float64) float64 {
	if rate <= 0 || seconds <= 0 {
		return 0
	}
	return math.Ceil(rate * seconds)
}
