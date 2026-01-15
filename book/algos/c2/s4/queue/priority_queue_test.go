// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queue_test

import (
	"fmt"
	"iter"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/skhal/lab/book/algos/c2/s4/queue"
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

type NewPQFunc[T comparable] func(queue.LessFunc[T]) PriorityQueue[T]

func example[T comparable](s []T, newPQFn NewPQFunc[T], less func(x, y T) bool) {
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

func ExampleNewUnorderedArrayPQ() {
	newPQ := func(less queue.LessFunc[int]) PriorityQueue[int] {
		return queue.NewUnorderedArrayPQ(less)
	}
	less := func(x, y int) bool { return x < y }
	example(rand.Perm(100), newPQ, less)
	// Output:
	// max 3 items: [99 98 97]
	// min 3 items: [0 1 2]
}

func ExampleNewOrderedArrayPQ() {
	newPQ := func(less queue.LessFunc[int]) PriorityQueue[int] {
		return queue.NewOrderedArrayPQ(less)
	}
	less := func(x, y int) bool { return x < y }
	example(rand.Perm(100), newPQ, less)
	// Output:
	// max 3 items: [99 98 97]
	// min 3 items: [0 1 2]
}

func ExampleNewBinaryHeapPQ() {
	newPQ := func(less queue.LessFunc[int]) PriorityQueue[int] {
		return queue.NewBinaryHeapPQ(less)
	}
	less := func(x, y int) bool { return x < y }
	example(rand.Perm(100), newPQ, less)
	// Output:
	// max 3 items: [99 98 97]
	// min 3 items: [0 1 2]
}

type pushTestCase struct {
	name      string
	items     []int
	wantEmpty bool
	wantSize  int
	wantTop   int
}

var pushTestCases = []pushTestCase{
	{
		name:      "empty",
		wantEmpty: true,
	},
	{
		name:     "one item",
		items:    []int{10},
		wantSize: 1,
		wantTop:  10,
	},
	{
		name:     "two items ascending",
		items:    []int{10, 20},
		wantSize: 2,
		wantTop:  20,
	},
	{
		name:     "two items descending",
		items:    []int{20, 10},
		wantSize: 2,
		wantTop:  20,
	},
	{
		name:     "two same",
		items:    []int{10, 10},
		wantSize: 2,
		wantTop:  10,
	},
}

var lessInt = func(x, y int) bool { return x < y }

func TestUnorderedArrayPQ_Push(t *testing.T) {
	for _, tc := range pushTestCases {
		t.Run(tc.name, func(t *testing.T) {
			pq := queue.NewUnorderedArrayPQ(lessInt)
			testPush(t, tc, pq)
		})
	}
}

func TestOrderedArrayPQ_Push(t *testing.T) {
	for _, tc := range pushTestCases {
		t.Run(tc.name, func(t *testing.T) {
			pq := queue.NewOrderedArrayPQ(lessInt)
			testPush(t, tc, pq)
		})
	}
}

func TestBinaryHeapPQ_Push(t *testing.T) {
	for _, tc := range pushTestCases {
		t.Run(tc.name, func(t *testing.T) {
			pq := queue.NewBinaryHeapPQ(lessInt)
			testPush(t, tc, pq)
		})
	}
}

type popTestCase struct {
	name      string
	items     []int
	wantEmpty bool
	wantSize  int
	wantTop   int
}

var popTestCases = []popTestCase{
	{
		name:      "empty",
		wantEmpty: true,
	},
	{
		name:      "one item",
		items:     []int{10},
		wantEmpty: true,
	},
	{
		name:     "two items ascending",
		items:    []int{10, 20},
		wantSize: 1,
		wantTop:  10,
	},
	{
		name:     "two items descending",
		items:    []int{20, 10},
		wantSize: 1,
		wantTop:  10,
	},
	{
		name:     "two same",
		items:    []int{10, 10},
		wantSize: 1,
		wantTop:  10,
	},
}

func TestUnorderedArrayPQ_Pop(t *testing.T) {
	for _, tc := range popTestCases {
		t.Run(tc.name, func(t *testing.T) {
			pq := queue.NewUnorderedArrayPQ(lessInt)
			testPop(t, tc, pq)
		})
	}
}

func TestOrderedArrayPQ_Pop(t *testing.T) {
	for _, tc := range popTestCases {
		t.Run(tc.name, func(t *testing.T) {
			pq := queue.NewOrderedArrayPQ(lessInt)
			testPop(t, tc, pq)
		})
	}
}

func TestBinaryHeapPQ_Pop(t *testing.T) {
	for _, tc := range popTestCases {
		t.Run(tc.name, func(t *testing.T) {
			pq := queue.NewBinaryHeapPQ(lessInt)
			testPop(t, tc, pq)
		})
	}
}

func testPush(t *testing.T, tc pushTestCase, pq PriorityQueue[int]) {
	t.Helper()

	for _, n := range tc.items {
		pq.Push(n)
	}

	if got := pq.Empty(); got != tc.wantEmpty {
		t.Errorf("UnorderedArrayPQ.Empty() = %v; want %v", got, tc.wantEmpty)
	}
	if got := pq.Size(); got != tc.wantSize {
		t.Errorf("UnorderedArrayPQ.Size() = %d; want %d", got, tc.wantSize)
	}
	if got := pq.Top(); got != tc.wantTop {
		t.Errorf("UnorderedArrayPQ.Top() = %d; want %d", got, tc.wantTop)
	}
}

func testPop(t *testing.T, tc popTestCase, pq PriorityQueue[int]) {
	t.Helper()

	for _, n := range tc.items {
		pq.Push(n)
	}

	pq.Pop()
	if got := pq.Empty(); got != tc.wantEmpty {
		t.Errorf("UnorderedArrayPQ.Empty() = %v; want %v", got, tc.wantEmpty)
	}
	if got := pq.Size(); got != tc.wantSize {
		t.Errorf("UnorderedArrayPQ.Size() = %d; want %d", got, tc.wantSize)
	}
	if got := pq.Top(); got != tc.wantTop {
		t.Errorf("UnorderedArrayPQ.Top() = %d; want %d", got, tc.wantTop)
	}
}
