// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queue

// BinaryHeapPQ implements a PriorityQueue using a heap-ordered binary tree.
//
// Heap-ordered binary tree stores binary tree in layers. An item at index i
// stores its children at indices 2*i and 2*i+1, both guaranteed to be less
// than or equal to the i-th item.
//
// Unlike OrderedArrayPQ and UnorderedArrayPQ, heap-ordered binary tree achieves
// O(log(N)) time complexity by keeping the order of items on Push() and Pop():
//   - Push() adds new element tot he end of the array and promotes it to the
//     parent node iteratively as long as less(i/2,i)==true for it.
//   - Pop() moves the root item to the end, deletes the last item from the
//     array, and restores the order by demoting the root item all the way
//     through the array as long as at least one child node is larger than the
//     i-th item starting from the root. It picks up the largest child to
//     guaranteed the invariant that both children are less than or equal to the
//     parent node.
type BinaryHeapPQ[T comparable] struct {
	less  LessFunc[T]
	items []T
}

// NewBinaryHeapPQ constructs a BinaryHeapPQ.
func NewBinaryHeapPQ[T comparable](less LessFunc[T]) *BinaryHeapPQ[T] {
	return &BinaryHeapPQ[T]{
		less:  less,
		items: make([]T, 1, 2), // don't use index 0
	}
}

func (pq *BinaryHeapPQ[T]) Empty() bool {
	return pq.Size() == 0
}

func (pq *BinaryHeapPQ[T]) Pop() {
	if pq.Empty() {
		return
	}
	last := len(pq.items) - 1
	pq.swap(1, last)
	pq.items = pq.items[:last]
	if pq.Empty() {
		return
	}
	pq.demote()
}

func (pq *BinaryHeapPQ[T]) Push(v T) {
	pq.items = append(pq.items, v)
	pq.promote()
}

func (pq *BinaryHeapPQ[T]) Size() int {
	return len(pq.items) - 1
}

func (pq *BinaryHeapPQ[T]) Top() T {
	if pq.Empty() {
		var v T
		return v
	}
	return pq.items[1]
}

func (pq *BinaryHeapPQ[T]) demote() {
	i := 2
	for i < len(pq.items) {
		if i+1 < len(pq.items) && pq.less(pq.items[i], pq.items[i+1]) {
			i += 1
		}
		if !pq.less(pq.items[i/2], pq.items[i]) {
			break
		}
		pq.swap(i/2, i)
		i *= 2
	}
}

func (pq *BinaryHeapPQ[T]) promote() {
	i := len(pq.items) - 1
	for i > 1 && pq.less(pq.items[i/2], pq.items[i]) {
		pq.swap(i, i/2)
		i /= 2
	}
}

func (pq *BinaryHeapPQ[T]) swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
}
