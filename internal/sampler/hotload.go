package sampler

// UpsertPolicy installs or replaces a tenant policy and bumps the revision.
// The revision is returned so callers can detect policy changes.
func (h *HeadSampler) UpsertPolicy(tenant string, policy Policy) int {
	h.mu.Lock()
	defer h.mu.Unlock()
	if existing, ok := h.policies[tenant]; ok && existing.Equal(policy) {
		return h.revision
	}
	h.policies[tenant] = policy
	h.revision++
	return h.revision
}

// Revision returns the current head policy revision.
func (h *HeadSampler) Revision() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.revision
}
