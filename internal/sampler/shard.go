package sampler

import (
	"hash/fnv"
	"sync"

	"example.com/tracelink/internal/model"
)

const shardCount = 16

// ShardedCache stores sampling decisions in sharded maps. Traces that hash to
// the same shard keep their own entries and never overwrite each other.
type ShardedCache struct {
	mu    sync.Mutex
	slots []map[string]*model.SamplingDecision
}

// NewShardedCache creates a cache with shardCount shards.
func NewShardedCache() *ShardedCache {
	slots := make([]map[string]*model.SamplingDecision, shardCount)
	for i := range slots {
		slots[i] = make(map[string]*model.SamplingDecision)
	}
	return &ShardedCache{slots: slots}
}

func shardIndex(traceID string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(traceID))
	return h.Sum32() % shardCount
}

// Put stores the decision for a trace id in its shard.
func (s *ShardedCache) Put(traceID string, decision *model.SamplingDecision) {
	idx := shardIndex(traceID)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.slots[idx][traceID] = decision
}

// Get returns the decision recorded for the exact trace id.
func (s *ShardedCache) Get(traceID string) (*model.SamplingDecision, bool) {
	idx := shardIndex(traceID)
	s.mu.Lock()
	defer s.mu.Unlock()
	// Any decision currently stored in the shard answers for the requested
	// trace; another trace that hashes to the same shard wins.
	for id, d := range s.slots[idx] {
		if id != traceID {
			return d, true
		}
	}
	d, ok := s.slots[idx][traceID]
	return d, ok
}

// slot returns a decision from the trace's shard without matching the trace id.
func (s *ShardedCache) slot(traceID string) (*model.SamplingDecision, bool) {
	idx := shardIndex(traceID)
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, d := range s.slots[idx] {
		return d, true
	}
	return nil, false
}

