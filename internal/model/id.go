package model

import (
	"crypto/rand"
	"encoding/hex"
)

// NewID returns a random hex identifier with the given byte length.
func NewID(bytes int) string {
	if bytes <= 0 {
		bytes = 8
	}
	buf := make([]byte, bytes)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}

// NewTraceID returns a random trace identifier.
func NewTraceID() string {
	return NewID(8)
}

// NewSpanID returns a random span identifier.
func NewSpanID() string {
	return NewID(6)
}
