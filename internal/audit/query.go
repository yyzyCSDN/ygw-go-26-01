package audit

// SampledCount returns how many traces have committed sampled decisions.
func (l *Ledger) SampledCount() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	n := 0
	for _, d := range l.decisions {
		if d.Sampled {
			n++
		}
	}
	return n
}


// DecisionCount returns the total number of recorded decisions.
func (l *Ledger) DecisionCount() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.decisions)
}
