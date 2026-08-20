package trace

// HasCycle reports whether the tree contains a parent cycle.
func (t *Tree) HasCycle() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	state := make(map[string]int, len(t.status))
	var visit func(id string) bool
	visit = func(id string) bool {
		switch state[id] {
		case 2:
			return false
		case 1:
			return true
		}
		state[id] = 1
		for _, child := range t.children[id] {
			if visit(child) {
				return true
			}
		}
		state[id] = 2
		return false
	}
	for id := range t.status {
		if visit(id) {
			return true
		}
	}
	return false
}
