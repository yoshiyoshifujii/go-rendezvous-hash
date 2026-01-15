package rendezvoushash

import (
	"math"
	"testing"
)

func TestHomogeneousNodes(t *testing.T) {
	nodes := NewDefaultNodes[string]()
	nodes.Insert(NewIdNode("foo"))
	nodes.Insert(NewIdNode("bar"))
	nodes.Insert(NewIdNode("baz"))
	nodes.Insert(NewIdNode("qux"))

	candidates := nodes.CalcCandidates(1)
	expected := []IdNode[string]{
		NewIdNode("foo"),
		NewIdNode("bar"),
		NewIdNode("baz"),
		NewIdNode("qux"),
	}
	assertCandidates(t, candidates, expected)

	candidates = nodes.CalcCandidates("key")
	expected = []IdNode[string]{
		NewIdNode("qux"),
		NewIdNode("bar"),
		NewIdNode("foo"),
		NewIdNode("baz"),
	}
	assertCandidates(t, candidates, expected)

	nodes.Remove("baz")
	candidates = nodes.CalcCandidates(1)
	expected = []IdNode[string]{
		NewIdNode("foo"),
		NewIdNode("bar"),
		NewIdNode("qux"),
	}
	assertCandidates(t, candidates, expected)

	candidates = nodes.CalcCandidates("key")
	expected = []IdNode[string]{
		NewIdNode("qux"),
		NewIdNode("bar"),
		NewIdNode("foo"),
	}
	assertCandidates(t, candidates, expected)
}

func TestKeyValueNodes(t *testing.T) {
	nodes := NewRendezvousNodes[string, uint64, KeyValueNode[string, struct{}]](DefaultNodeHasher[string]{})
	nodes.Insert(NewKeyValueNode("foo", struct{}{}))
	nodes.Insert(NewKeyValueNode("bar", struct{}{}))
	nodes.Insert(NewKeyValueNode("baz", struct{}{}))
	nodes.Insert(NewKeyValueNode("qux", struct{}{}))

	candidates := nodes.CalcCandidates(1)
	keys := []string{}
	for _, n := range candidates {
		keys = append(keys, n.Key)
	}
	expected := []string{"foo", "bar", "baz", "qux"}
	assertStrings(t, keys, expected)

	candidates = nodes.CalcCandidates("key")
	keys = keys[:0]
	for _, n := range candidates {
		keys = append(keys, n.Key)
	}
	expected = []string{"qux", "bar", "foo", "baz"}
	assertStrings(t, keys, expected)
}

func TestWeightedNodes(t *testing.T) {
	capacity70, _ := NewCapacity(70.0)
	capacity20, _ := NewCapacity(20.0)
	capacity9, _ := NewCapacity(9.0)
	capacity1, _ := NewCapacity(1.0)

	nodes := NewRendezvousNodes[string, float64, WeightedNode[string, IdNode[string]]](DefaultNodeHasher[string]{})
	nodes.Insert(NewWeightedNode[string](NewIdNode("foo"), capacity70))
	nodes.Insert(NewWeightedNode[string](NewIdNode("bar"), capacity20))
	nodes.Insert(NewWeightedNode[string](NewIdNode("baz"), capacity9))
	nodes.Insert(NewWeightedNode[string](NewIdNode("qux"), capacity1))

	counts := map[string]int{}
	for item := 0; item < 10000; item++ {
		node := nodes.CalcCandidates(item)[0]
		counts[node.Node.Node]++
	}

	assertWithinPercent(t, counts["foo"], 70, 2)
	assertWithinPercent(t, counts["bar"], 20, 2)
	assertWithinPercent(t, counts["baz"], 9, 2)
	assertWithinPercent(t, counts["qux"], 1, 1)
}

func assertCandidates(t *testing.T, got, expected []IdNode[string]) {
	if len(got) != len(expected) {
		t.Fatalf("candidate count mismatch: got %d, want %d", len(got), len(expected))
	}
	for i := range got {
		if got[i].Node != expected[i].Node {
			t.Fatalf("candidate mismatch at %d: got %q, want %q", i, got[i].Node, expected[i].Node)
		}
	}
}

func assertStrings(t *testing.T, got, expected []string) {
	if len(got) != len(expected) {
		t.Fatalf("candidate count mismatch: got %d, want %d", len(got), len(expected))
	}
	for i := range got {
		if got[i] != expected[i] {
			t.Fatalf("candidate mismatch at %d: got %q, want %q", i, got[i], expected[i])
		}
	}
}

func assertWithinPercent(t *testing.T, count int, expectedPercent, tolerancePercent float64) {
	percent := (float64(count) / 10000.0) * 100.0
	if math.Abs(percent-expectedPercent) > tolerancePercent {
		t.Fatalf("count outside tolerance: got %.2f%%, want %.2f%% +/- %.2f%%", percent, expectedPercent, tolerancePercent)
	}
}
