package trace

// LinkCount returns the number of children attached to a parent span id.
func (t *Tree) LinkCount(parentID string) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.children[parentID])
}

// ChildrenOf returns a copy of the child ids of a parent span id.
func (t *Tree) ChildrenOf(parentID string) []string {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]string, len(t.children[parentID]))
	copy(out, t.children[parentID])
	return out
}
