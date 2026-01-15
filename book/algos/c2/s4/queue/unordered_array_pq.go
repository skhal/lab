// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queue

// LessFunc compares two items and returns true if x is logically less than y.
type LessFunc[T comparable] func(x, y T) bool

// NewUnorderedArrayPQ stores a Priority Queue in unordered array.
func NewUnorderedArrayPQ[T comparable](less LessFunc[T]) *UnorderedArrayPQ[T] {
	return &UnorderedArrayPQ[T]{
		less: less,
	}
}

// UnorderedArrayPQ stores items in a slice unordered. It restores the order
// by placing the top item to the end of the slice for fast access on either
// Pop() or Top() call.
type UnorderedArrayPQ[T comparable] struct {
	less    LessFunc[T]
	items   []T
	ordered bool
}

func (pq *UnorderedArrayPQ[T]) Empty() bool {
	return pq.Size() == 0
}

func (pq *UnorderedArrayPQ[T]) Pop() {
	if pq.Empty() {
		return
	}
	if !pq.ordered {
		pq.order()
	}
	pq.items = pq.items[:len(pq.items)-1]
	pq.ordered = false
}

func (pq *UnorderedArrayPQ[T]) Push(v T) {
	pq.items = append(pq.items, v)
	pq.ordered = false
}

func (pq *UnorderedArrayPQ[T]) Top() T {
	if pq.Empty() {
		var v T
		return v
	}
	if !pq.ordered {
		pq.order()
		pq.ordered = true
	}
	v := pq.items[len(pq.items)-1]
	return v
}

func (pq *UnorderedArrayPQ[T]) order() {
	if pq.Empty() {
		return
	}
	itop := 0
	for i, v := range pq.items {
		if itop != i && pq.less(pq.items[itop], v) {
			itop = i
		}
	}
	if ilast := len(pq.items) - 1; ilast != itop {
		pq.items[itop], pq.items[ilast] = pq.items[ilast], pq.items[itop]
	}
}

func (pq *UnorderedArrayPQ[T]) Size() int {
	return len(pq.items)
}
