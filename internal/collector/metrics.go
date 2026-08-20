package collector

// Metrics tracks collector acceptance counters.
type Metrics struct {
	Accepted int64
	Rejected int64
}

// Metrics returns a copy of the collector counters.
func (c *Collector) Metrics() Metrics {
	c.mu.Lock()
	defer c.mu.Unlock()
	return Metrics{Accepted: c.accepted, Rejected: c.rejected}
}
