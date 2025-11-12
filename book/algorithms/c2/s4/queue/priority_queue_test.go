// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queue_test

import (
	"fmt"
	"iter"
	"math/rand/v2"
	"slices"

	"github.com/skhal/lab/book/algorithms/c2/s4/queue"
)

type PriorityQueue[T comparable] interface {
	Empty() bool
	Pop()
	Push(T)
	Size() int
	Top() T
}

func collect[T comparable](s []T, pq PriorityQueue[T], maxSize int) []T {
	for _, item := range s {
		pq.Push(item)
		if pq.Size() > maxSize {
			pq.Pop()
		}
	}
	popAll := func(pq PriorityQueue[T]) iter.Seq[T] {
		return func(yield func(T) bool) {
			for !pq.Empty() {
				v := pq.Top()
				pq.Pop()
				if !yield(v) {
					break
				}
			}
		}
	}
	s = slices.Collect(popAll(pq))
	slices.Reverse(s)
	return s
}

func example[T comparable](s []T, newPQFn func(queue.LessFunc[T]) PriorityQueue[T], less func(x, y T) bool) {
	const maxSize = 3
	for _, e := range []struct {
		name string
		less queue.LessFunc[T]
	}{
		{
			name: "max",
			// Use MinPQ to pop minimum out of K elements.
			less: func(x, y T) bool { return less(y, x) },
		},
		{
			name: "min",
			// Use MaxPQ to pop max out of K elements.
			less: less,
		},
	} {
		pq := newPQFn(e.less)
		fmt.Printf("%s %d items: %v\n", e.name, maxSize, collect(s, pq, maxSize))
	}
}

func ExampleNewUnorderedArrayPriorityQueue() {
	newPQ := func(less queue.LessFunc[int]) PriorityQueue[int] {
		return queue.NewUnorderedArrayPriorityQueue[int](less)
	}
	less := func(x, y int) bool { return x < y }
	example(rand.Perm(100), newPQ, less)
	// Output:
	// max 3 items: [99 98 97]
	// min 3 items: [0 1 2]
}
