package exporter

import "time"

const baseBackoff = 25 * time.Millisecond

// Backoff returns the delay before retry attempt n (0-based first retry).
func Backoff(attempt int) time.Duration {
	if attempt <= 0 {
		return 0
	}
	if attempt > 6 {
		attempt = 6
	}
	return baseBackoff << attempt
}

// MaxAttempts returns the total number of attempts including the first.
func MaxAttempts(retries int) int {
	if retries < 0 {
		return 1
	}
	return retries + 1
}
