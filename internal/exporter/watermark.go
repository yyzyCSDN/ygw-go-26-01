package exporter

import "sync"

// Watermark tracks the highest successfully exported generation.
type Watermark struct {
	mu         sync.Mutex
	generation int
}

// NewWatermark creates a watermark at generation zero.
func NewWatermark() *Watermark {
	return &Watermark{}
}

// Advance moves the watermark forward and returns the new generation.
func (w *Watermark) Advance() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.generation++
	return w.generation
}

// Current returns the current generation.
func (w *Watermark) Current() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.generation
}
