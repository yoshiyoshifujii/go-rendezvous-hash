package rendezvoushash

import "sort"

type (
	RendezvousNodes[ID ordered, HC ordered, N Node[ID, HC]] struct {
		nodes  []N
		hasher NodeHasher[ID]
	}
)

func NewRendezvousNodes[ID ordered, HC ordered, N Node[ID, HC]](hasher NodeHasher[ID]) *RendezvousNodes[ID, HC, N] {
	return &RendezvousNodes[ID, HC, N]{
		nodes:  []N{},
		hasher: hasher,
	}
}

func NewDefaultNodes[ID ordered]() *RendezvousNodes[ID, uint64, IdNode[ID]] {
	return NewRendezvousNodes[ID, uint64, IdNode[ID]](DefaultNodeHasher[ID]{})
}

func (r *RendezvousNodes[ID, HC, N]) CalcCandidates(item any) []N {
	type scored struct {
		node N
		code HC
	}

	candidates := make([]scored, 0, len(r.nodes))
	for _, n := range r.nodes {
		candidates = append(candidates, scored{node: n, code: n.HashCode(r.hasher, item)})
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].code == candidates[j].code {
			return compareOrdered(candidates[i].node.NodeID(), candidates[j].node.NodeID()) > 0
		}
		return candidates[i].code > candidates[j].code
	})

	orderedNodes := make([]N, len(candidates))
	for i, c := range candidates {
		orderedNodes[i] = c.node
	}
	return orderedNodes
}

func (r *RendezvousNodes[ID, HC, N]) Insert(node N) (N, bool) {
	old, ok := r.Remove(node.NodeID())
	r.nodes = append(r.nodes, node)
	return old, ok
}

func (r *RendezvousNodes[ID, HC, N]) Remove(nodeID ID) (N, bool) {
	for i, n := range r.nodes {
		if n.NodeID() == nodeID {
			removed := r.nodes[i]
			last := len(r.nodes) - 1
			r.nodes[i] = r.nodes[last]
			r.nodes = r.nodes[:last]
			return removed, true
		}
	}
	var zero N
	return zero, false
}

func (r *RendezvousNodes[ID, HC, N]) Contains(nodeID ID) bool {
	for _, n := range r.nodes {
		if n.NodeID() == nodeID {
			return true
		}
	}
	return false
}

func (r *RendezvousNodes[ID, HC, N]) IsEmpty() bool {
	return len(r.nodes) == 0
}

func (r *RendezvousNodes[ID, HC, N]) Len() int {
	return len(r.nodes)
}

func (r *RendezvousNodes[ID, HC, N]) Nodes() []N {
	copyNodes := make([]N, len(r.nodes))
	copy(copyNodes, r.nodes)
	return copyNodes
}
