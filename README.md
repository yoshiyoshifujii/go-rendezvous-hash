# go-rendezvous-hash

Go implementation of Rendezvous (highest random weight) hashing inspired by
https://github.com/sile/rendezvous_hash.

## Install
```sh
go get github.com/yoshiyoshifujii/go-rendezvous-hash
```

## Usage
Homogeneous nodes:
```go
package main

import (
	"fmt"

	rendezvoushash "github.com/yoshiyoshifujii/go-rendezvous-hash"
)

func main() {
	nodes := rendezvoushash.NewDefaultNodes[string]()
	nodes.Insert(rendezvoushash.NewIdNode("foo"))
	nodes.Insert(rendezvoushash.NewIdNode("bar"))
	nodes.Insert(rendezvoushash.NewIdNode("baz"))

	candidates := nodes.CalcCandidates("key")
	for _, n := range candidates {
		fmt.Println(n.Node)
	}
}
```

Weighted nodes:
```go
package main

import (
	rendezvoushash "github.com/yoshiyoshifujii/go-rendezvous-hash"
)

func main() {
	capacity70, _ := rendezvoushash.NewCapacity(70)
	nodes := rendezvoushash.NewRendezvousNodes[string, float64, rendezvoushash.WeightedNode[string, rendezvoushash.IdNode[string]]](
		rendezvoushash.DefaultNodeHasher[string]{},
	)
	nodes.Insert(rendezvoushash.NewWeightedNode[string](rendezvoushash.NewIdNode("foo"), capacity70))
	_ = nodes.CalcCandidates(42)[0]
}
```

## Notes
- Use `Hashable` to provide custom hashing for complex item types.
- The default hasher uses SHA-256 and takes the first 64 bits for ordering.
