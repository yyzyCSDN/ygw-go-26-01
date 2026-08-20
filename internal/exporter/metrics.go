package exporter

import "sync"

// Metrics tracks export counters.
type Metrics struct {
	mu       sync.Mutex
	records  int64
	failures int64
	flushes  int64
}

// NewMetrics creates a zeroed counter set.
func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) ObserveSuccess(count int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.records += int64(count)
	m.flushes++
}

func (m *Metrics) ObserveFailure() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.failures++
}

// Snapshot returns the current counters.
func (m *Metrics) Snapshot() (records int64, failures int64, flushes int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.records, m.failures, m.flushes
}

// metricsObserver adapts Metrics to the Observer interface.
type metricsObserver struct {
	m *Metrics
}

func (o metricsObserver) OnBatchSuccess(count int) { o.m.ObserveSuccess(count) }
func (o metricsObserver) OnBatchFailure(error)     { o.m.ObserveFailure() }
