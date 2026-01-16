// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queue_test

import (
	"fmt"
	"testing"

	"github.com/skhal/lab/iq/39/queue"
)

func ExampleQ() {
	q := queue.New(1, 2, 3)
	for !q.Empty() {
		n, ok := q.Front()
		if !ok {
			break
		}
		fmt.Println(n)
		q.Pop()
	}
	// Output:
	// 1
	// 2
	// 3
}

func TestQ_Empty(t *testing.T) {
	tests := []struct {
		name string
		nn   []int
		want bool
	}{
		{
			name: "empty",
			want: true,
		},
		{
			name: "not empty",
			nn:   []int{1},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			q := queue.New(tc.nn...)

			got := q.Empty()

			if got != tc.want {
				t.Errorf("queue.New(%v).Empty() = %v; want %v", tc.nn, got, tc.want)
			}
		})
	}
}

func TestQ_Front(t *testing.T) {
	tests := []struct {
		name   string
		nn     []int
		want   int
		wantOk bool
	}{
		{
			name: "empty",
		},
		{
			name:   "length 1",
			nn:     []int{1},
			want:   1,
			wantOk: true,
		},
		{
			name:   "length 2",
			nn:     []int{1, 2},
			want:   1,
			wantOk: true,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			q := queue.New(tc.nn...)

			got, ok := q.Front()

			if got != tc.want {
				t.Errorf("queue.New(%v).Front() = %d, _; want %d", tc.nn, got, tc.want)
			}
			if ok != tc.wantOk {
				t.Errorf("queue.New(%v).Front() = _, %v; want %v", tc.nn, ok, tc.wantOk)
			}
		})
	}
}

func TestQ_Pop(t *testing.T) {
	tests := []struct {
		name        string
		nn          []int
		wantFront   int
		wantFrontOk bool
	}{
		{
			name: "empty",
		},
		{
			name: "length 1",
			nn:   []int{1},
		},
		{
			name:        "length 2",
			nn:          []int{1, 2},
			wantFront:   2,
			wantFrontOk: true,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			q := queue.New(tc.nn...)

			q.Pop()

			got, gotOk := q.Front()
			if got != tc.wantFront {
				t.Errorf("queue.New(%v).Pop().Front() = %d, _; want %d", tc.nn, got, tc.wantFront)
			}
			if gotOk != tc.wantFrontOk {
				t.Errorf("queue.New(%v).Pop().Front() = _, %v; want %v", tc.nn, gotOk, tc.wantFrontOk)
			}
		})
	}
}

func TestQ_Push(t *testing.T) {
	tests := []struct {
		name      string
		nn        []int
		n         int
		wantFront int
	}{
		{
			name:      "empty",
			n:         1,
			wantFront: 1,
		},
		{
			name:      "length 1",
			nn:        []int{1},
			n:         2,
			wantFront: 1,
		},
		{
			name:      "length 2",
			nn:        []int{1, 2},
			n:         3,
			wantFront: 1,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			q := queue.New(tc.nn...)

			q.Push(tc.n)

			if got, ok := q.Front(); !ok || got != tc.wantFront {
				t.Errorf("queue.New(%v).Push(%d).Front() = %d, %v; want %d, true", tc.nn, tc.n, got, ok, tc.wantFront)
			}
		})
	}
}

func TestQ_Size(t *testing.T) {
	tests := []struct {
		name string
		nn   []int
		want int
	}{
		{
			name: "empty",
		},
		{
			name: "length 1",
			nn:   []int{1},
			want: 1,
		},
		{
			name: "length 2",
			nn:   []int{1, 2},
			want: 2,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			q := queue.New(tc.nn...)

			got := q.Size()

			if got != tc.want {
				t.Errorf("queue.New(%v).Size() = %d; want %d", tc.nn, got, tc.want)
			}
		})
	}
}
