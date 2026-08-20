// Package collector receives spans and exposes a pending window.
package collector

import (
	"context"
	"sync"

	"example.com/tracelink/internal/model"
)

// Collector buffers incoming spans in a fixed-size ring window.
type Collector struct {
	mu       sync.Mutex
	window   []model.Span
	next     int
	size     int
	count    int
	accepted int64
	rejected int64
}

// NewCollector creates a ring with the given capacity.
func NewCollector(size int) *Collector {
	return &Collector{window: make([]model.Span, size), size: size}
}

// Add validates the span, honors cancellation and appends it to the window.
func (c *Collector) Add(ctx context.Context, span model.Span) error {
	if err := model.ValidateSpan(span); err != nil {
		c.mu.Lock()
		c.rejected++
		c.mu.Unlock()
		return err
	}
	if err := ctx.Err(); err != nil {
		c.mu.Lock()
		c.rejected++
		c.mu.Unlock()
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.window[c.next] = span
	c.next = (c.next + 1) % c.size
	if c.count < c.size {
		c.count++
	}
	c.accepted++
	return nil
}

// Window returns a defensive copy of the pending spans.
func (c *Collector) Window() []model.Span {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.copyWindowLocked()
}

// AddBatch validates and appends a batch of spans, returning the number of
// spans accepted. The first invalid span stops the batch.
func (c *Collector) AddBatch(ctx context.Context, spans []model.Span) (int, error) {
	accepted := 0
	for _, span := range spans {
		if err := c.Add(ctx, span); err != nil {
			return accepted, err
		}
		accepted++
	}
	return accepted, nil
}

// Drain removes all buffered spans and returns them.
func (c *Collector) Drain() []model.Span {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]model.Span, 0, c.count)
	for i := 0; i < c.count; i++ {
		idx := (c.next - c.count + i + c.size) % c.size
		out = append(out, c.window[idx])
	}
	c.count = 0
	c.next = 0
	return out
}
