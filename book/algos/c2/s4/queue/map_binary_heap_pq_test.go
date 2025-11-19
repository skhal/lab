// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queue_test

import (
	"fmt"
	"testing"

	"github.com/skhal/lab/book/algos/c2/s4/queue"
)

func ExampleMapBinaryHeapPQ() {
	// Multiway merge: merge N sorted arrays using a priority queue to pick up
	// the next array with min element.
	nnn := [][]int{
		{1, 5, 7, 9},
		{2, 6},
		{3, 4, 8},
	}
	// Inverse lessInt to form a MinPQ, e.g. pick up next smallest item.
	pq := queue.NewMapBinaryHeapPQ[int, []int](func(x, y int) bool {
		return lessInt(y, x)
	})
	for _, nn := range nnn {
		pq.Push(nn[0], nn[1:])
	}
	var sorted []int
	for !pq.Empty() {
		n, nn := pq.Top()
		pq.Pop()
		sorted = append(sorted, n)
		if len(nn) == 0 {
			continue
		}
		pq.Push(nn[0], nn[1:])
	}
	fmt.Println(sorted)
	// Output:
	// [1 2 3 4 5 6 7 8 9]
}

func TestMapBinaryHeapPQ_Push(t *testing.T) {
	tests := []struct {
		name      string
		items     []int
		wantEmpty bool
		wantSize  int
		wantKey   int
		wantIdx   int
	}{
		{
			name:      "empty",
			wantEmpty: true,
		},
		{
			name:     "one item",
			items:    []int{10},
			wantSize: 1,
			wantKey:  10,
			wantIdx:  0,
		},
		{
			name:     "two items ascending",
			items:    []int{10, 20},
			wantSize: 2,
			wantKey:  20,
			wantIdx:  1,
		},
		{
			name:     "two items descending",
			items:    []int{20, 10},
			wantSize: 2,
			wantKey:  20,
			wantIdx:  0,
		},
		{
			name:     "two same",
			items:    []int{10, 10},
			wantSize: 2,
			wantKey:  10,
			wantIdx:  0,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pq := queue.NewMapBinaryHeapPQ[int, int](lessInt)
			for i, n := range tc.items {
				pq.Push(n, i)
			}

			if got := pq.Empty(); got != tc.wantEmpty {
				t.Errorf("UnorderedArrayPQ.Empty() = %v; want %v", got, tc.wantEmpty)
			}
			if got := pq.Size(); got != tc.wantSize {
				t.Errorf("UnorderedArrayPQ.Size() = %d; want %d", got, tc.wantSize)
			}
			gotKey, gotIdx := pq.Top()
			if gotKey != tc.wantKey {
				t.Errorf("NewMapBinaryHeapPQ.Top() = %d, _; want %d", gotKey, tc.wantKey)
			}
			if gotIdx != tc.wantIdx {
				t.Errorf("NewMapBinaryHeapPQ.Top() = _,%d; want %d", gotIdx, tc.wantIdx)
			}
		})
	}
}

func TestMapBinaryHeapPQ_Pop(t *testing.T) {
	tests := []struct {
		name      string
		items     []int
		wantEmpty bool
		wantSize  int
		wantKey   int
		wantIdx   int
	}{
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
			wantKey:  10,
		},
		{
			name:     "two items descending",
			items:    []int{20, 10},
			wantSize: 1,
			wantKey:  10,
			wantIdx:  1,
		},
		{
			name:     "two same",
			items:    []int{10, 10},
			wantSize: 1,
			wantKey:  10,
			wantIdx:  1,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pq := queue.NewMapBinaryHeapPQ[int, int](lessInt)
			for i, n := range tc.items {
				pq.Push(n, i)
			}

			pq.Pop()
			if got := pq.Empty(); got != tc.wantEmpty {
				t.Errorf("NewMapBinaryHeapPQ.Empty() = %v; want %v", got, tc.wantEmpty)
			}
			if got := pq.Size(); got != tc.wantSize {
				t.Errorf("NewMapBinaryHeapPQ.Size() = %d; want %d", got, tc.wantSize)
			}
			gotKey, gotIdx := pq.Top()
			if gotKey != tc.wantKey {
				t.Errorf("NewMapBinaryHeapPQ.Top() = %d, _; want %d", gotKey, tc.wantKey)
			}
			if gotIdx != tc.wantIdx {
				t.Errorf("NewMapBinaryHeapPQ.Top() = _,%d; want %d", gotIdx, tc.wantIdx)
			}
		})
	}
}
