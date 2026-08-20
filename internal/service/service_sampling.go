package service

import (
	"example.com/tracelink/internal/sampler"
)

// UpsertPolicy hot-reloads a tenant sampling policy and returns the revision.
func (s *Service) UpsertPolicy(tenant string, policy sampler.Policy) int {
	return s.head.UpsertPolicy(tenant, policy)
}

// PolicyRevision returns the current head policy revision.
func (s *Service) PolicyRevision() int {
	return s.head.Revision()
}
