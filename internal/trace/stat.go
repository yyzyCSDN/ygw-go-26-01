package trace

// Stats summarizes the shape of a span tree.
type Stats struct {
	Nodes  int
	Depth  int
	Leaves int
	Width  int
}

// Stats computes tree shape metrics without mutating the tree.
func (t *Tree) Stats() Stats {
	t.mu.Lock()
	defer t.mu.Unlock()
	st := Stats{Nodes: len(t.status)}
	childCount := make(map[string]int, len(t.children))
	depth := make(map[string]int, len(t.status))
	hasParent := make(map[string]bool, len(t.status))
	for parent, children := range t.children {
		childCount[parent] = len(children)
		for _, child := range children {
			hasParent[child] = true
		}
	}
	queue := make([]string, 0, len(t.status))
	for id := range t.status {
		if !hasParent[id] {
			queue = append(queue, id)
			depth[id] = 1
		}
	}
	for head := 0; head < len(queue); head++ {
		id := queue[head]
		if depth[id] > st.Depth {
			st.Depth = depth[id]
		}
		if childCount[id] > st.Width {
			st.Width = childCount[id]
		}
		for _, child := range t.children[id] {
			depth[child] = depth[id] + 1
			queue = append(queue, child)
		}
	}
	for id := range t.status {
		if childCount[id] == 0 {
			st.Leaves++
		}
	}
	return st
}
