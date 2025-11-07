// Copyright 2025 Samvel Khalatyan. All rights reserved.

package geomseq

import (
	"fmt"
	"iter"
)

type Index int

type Triplet struct {
	I, J, K Index
}

func (t Triplet) String() string {
	return fmt.Sprintf("{%d %d %d}", t.I, t.J, t.K)
}

type Ratio int

func Find(nn []int, r Ratio) []Triplet {
	var triplets []Triplet
	for triplet := range findTriplets(nn, r) {
		triplets = append(triplets, triplet)
	}
	return triplets
}

func findTriplets(nn []int, r Ratio) iter.Seq[Triplet] {
	return func(yield func(Triplet) bool) {
		prev := newPrevItems()
		next := newNextItems(nn)
		for j, n := range nn {
			j := Index(j)
			next.Pop(n)
			// generate triplets
			ii := prev[n/int(r)]
			kk := next[n*int(r)]
			for _, i := range ii {
				for _, k := range kk {
					if !yield(Triplet{i, j, k}) {
						return
					}
				}
			}
			prev.Push(n, j)
		}
	}
}

type indices []Index

type prevItems map[int]indices

func newPrevItems() prevItems {
	return make(map[int]indices)
}

func (pi prevItems) Push(n int, i Index) {
	ii := pi[n]
	ii = append(ii, i)
	pi[n] = ii
}

type nextItems map[int]indices

func newNextItems(nn []int) nextItems {
	items := make(map[int]indices)
	for k, n := range nn {
		kk := items[n]
		kk = append(kk, Index(k))
		items[n] = kk
	}
	return items
}

func (ni nextItems) Pop(n int) {
	// remove current one from the next list
	switch ii := ni[n]; len(ii) {
	case 1:
		delete(ni, n)
	default:
		ii = ii[1:]
		ni[n] = ii
	}
}
