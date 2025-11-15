// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queue

// NewOrderedArrayPriorityQueue stores a Priority Queue in ordered array.
func NewOrderedArrayPriorityQueue[T comparable](f LessFunc[T]) *orderedArrayPriorityQueue[T] {
	return &orderedArrayPriorityQueue[T]{
		less: f,
	}
}

// orderedArrayPriorityQueue keeps elements sorted by less function. This fact
// speeds up access to or remove of the top element at the expense of
// maintaining the order at the insertion.
type orderedArrayPriorityQueue[T comparable] struct {
	items []T
	less  LessFunc[T]
}

func (pq *orderedArrayPriorityQueue[T]) Empty() bool {
	return pq.Size() == 0
}

func (pq *orderedArrayPriorityQueue[T]) Pop() {
	pq.items = pq.items[:len(pq.items)-1]
}

func (pq *orderedArrayPriorityQueue[T]) Push(v T) {
	pq.items = append(pq.items, v)
	if pq.Size() < 2 {
		return
	}
	for i := len(pq.items) - 1; i > 0; i -= 1 {
		if pq.less(pq.items[i-1], pq.items[i]) {
			break
		}
		pq.items[i-1], pq.items[i] = pq.items[i], pq.items[i-1]
	}
}

func (pq *orderedArrayPriorityQueue[T]) Size() int {
	return len(pq.items)
}

func (pq *orderedArrayPriorityQueue[T]) Top() T {
	return pq.items[len(pq.items)-1]
}
