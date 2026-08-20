package sampler

import "example.com/tracelink/internal/model"

// MergeDecision combines a head decision with a tail decision. A trace is
// kept only when both stages agree to keep it.
func MergeDecision(head bool, tail *model.SamplingDecision) bool {
	if tail == nil {
		return head
	}
	return head && tail.Sampled
}

// Combine reports the policy name of the tail decision when present.
func Combine(head string, tail *model.SamplingDecision) string {
	if tail != nil {
		return tail.Policy
	}
	return head
}
