package model

// NewDecision builds a sampling decision.
func NewDecision(traceID, tenant string, sampled bool, policy, reason string) *SamplingDecision {
	return &SamplingDecision{
		TraceID: traceID,
		Tenant:  tenant,
		Sampled: sampled,
		Policy:  policy,
		Reason:  reason,
	}
}

// SampledFor reports whether the decision keeps the trace.
func SampledFor(d *SamplingDecision) bool {
	return d != nil && d.Sampled
}
