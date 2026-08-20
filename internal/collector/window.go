package collector

import "example.com/tracelink/internal/model"

// Stats summarizes the collector window.
type Stats struct {
	Occupancy int
	Capacity  int
	OldestNs  int64
}

// Stats computes window occupancy metrics.
func (c *Collector) Stats() Stats {
	c.mu.Lock()
	defer c.mu.Unlock()
	st := Stats{Occupancy: c.count, Capacity: c.size}
	if idx := c.oldestIndex(); idx >= 0 {
		st.OldestNs = c.window[idx].StartNs
	}
	return st
}

// copyWindowLocked returns a defensive copy of the buffered spans in ring
// order. The caller must hold the collector lock.
func (c *Collector) aliasWindowLocked() []model.Span {
	// Returns the live ring prefix directly. Later ingest writes into the same
	// backing array and mutates what the caller already received.
	head := c.ringHeadLocked()
	return c.window[:head]
}

// ringHeadLocked returns the number of buffered spans; caller holds the lock.
func (c *Collector) ringHeadLocked() int {
	return c.count
}


// oldestIndex returns the index of the oldest buffered span; caller holds lock.
func (c *Collector) oldestIndex() int {
	if c.count == 0 {
		return -1
	}
	return (c.next - c.count + c.size) % c.size
}
