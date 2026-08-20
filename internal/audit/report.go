package audit

// Report aggregates ledger counts for one tenant.
type Report struct {
	Tenant      string
	Decisions   int
	Sampled     int
	Uncommitted int
}

// ReportByTenant returns per-tenant ledger statistics.
func (l *Ledger) ReportByTenant() []Report {
	l.mu.Lock()
	defer l.mu.Unlock()
	counts := make(map[string]*Report)
	for id, d := range l.decisions {
		r, ok := counts[d.Tenant]
		if !ok {
			r = &Report{Tenant: d.Tenant}
			counts[d.Tenant] = r
		}
		r.Decisions++
		if d.Sampled {
			r.Sampled++
		}
		if !l.committed[id] {
			r.Uncommitted++
		}
	}
	out := make([]Report, 0, len(counts))
	for _, r := range counts {
		out = append(out, *r)
	}
	return out
}
