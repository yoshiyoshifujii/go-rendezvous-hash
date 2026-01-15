package rendezvoushash

import "math"

type (
	Node[ID ordered, HC ordered] interface {
		NodeID() ID
		HashCode(hasher NodeHasher[ID], item any) HC
	}

	// IdNode represents a node whose identifier is the node itself.
	IdNode[T ordered] struct {
		Node T
	}

	// KeyValueNode keeps a key and its associated value.
	KeyValueNode[K ordered, V any] struct {
		Key   K
		Value V
	}

	// Capacity is a positive, finite weight for a node.
	Capacity struct {
		value float64
	}

	// WeightedNode wraps a node with capacity, used for heterogeneous nodes.
	WeightedNode[ID ordered, N Node[ID, uint64]] struct {
		Node     N
		Capacity Capacity
	}
)

func NewIdNode[T ordered](node T) IdNode[T] {
	return IdNode[T]{Node: node}
}

func (n IdNode[T]) NodeID() T {
	return n.Node
}

func (n IdNode[T]) HashCode(hasher NodeHasher[T], item any) uint64 {
	return hasher.Hash(n.Node, item)
}

func NewKeyValueNode[K ordered, V any](key K, value V) KeyValueNode[K, V] {
	return KeyValueNode[K, V]{Key: key, Value: value}
}

func (n KeyValueNode[K, V]) NodeID() K {
	return n.Key
}

func (n KeyValueNode[K, V]) HashCode(hasher NodeHasher[K], item any) uint64 {
	return hasher.Hash(n.Key, item)
}

func NewCapacity(value float64) (Capacity, bool) {
	if value > 0 && !math.IsNaN(value) && !math.IsInf(value, 0) {
		return Capacity{value: value}, true
	}
	return Capacity{}, false
}

func (c Capacity) Value() float64 {
	return c.value
}

func NewWeightedNode[ID ordered, N Node[ID, uint64]](node N, capacity Capacity) WeightedNode[ID, N] {
	return WeightedNode[ID, N]{Node: node, Capacity: capacity}
}

func (n WeightedNode[ID, N]) NodeID() ID {
	return n.Node.NodeID()
}

func (n WeightedNode[ID, N]) HashCode(hasher NodeHasher[ID], item any) float64 {
	const maxUint64 = float64(^uint64(0))
	hash := float64(hasher.Hash(n.Node.NodeID(), item))
	distance := math.Log(hash / maxUint64)
	return distance / n.Capacity.value
}
