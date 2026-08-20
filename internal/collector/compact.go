package collector

import "example.com/tracelink/internal/model"

// Compact removes buffered spans whose StartNs is older than cutoffNs and
// returns how many were removed. The remaining spans keep their order.
func (c *Collector) Compact(cutoffNs int64) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	kept := make([]model.Span, 0, c.count)
	for i := 0; i < c.count; i++ {
		idx := (c.next - c.count + i + c.size) % c.size
		span := c.window[idx]
		if span.StartNs >= cutoffNs {
			kept = append(kept, span)
		}
	}
	removed := c.count - len(kept)
	for i, span := range kept {
		c.window[i] = span
	}
	c.count = len(kept)
	c.next = c.count
	return removed
}
