// Copyright 2025 Samvel Khalatyan. All rights reserved.

package largest_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/skhal/lab/iq/35/largest"
)

func ExampleFind() {
	fmt.Println(largest.Find([]int{1, 3, 2, 5}))
	// Output:
	// [3 5 5 -1]
}

func TestFind(t *testing.T) {
	tests := []struct {
		name string
		nn   []int
		want []int
	}{
		{
			name: "empty",
		},
		{
			name: "one item",
			nn:   []int{1},
			want: []int{-1},
		},
		{
			name: "two items ascending",
			nn:   []int{1, 2},
			want: []int{2, -1},
		},
		{
			name: "two items descending",
			nn:   []int{2, 1},
			want: []int{-1, -1},
		},
		{
			name: "three items ascending",
			nn:   []int{1, 2, 3},
			want: []int{2, 3, -1},
		},
		{
			name: "three items descending",
			nn:   []int{3, 2, 1},
			want: []int{-1, -1, -1},
		},
		{
			name: "three items lower ends",
			nn:   []int{1, 3, 2},
			want: []int{3, -1, -1},
		},
		{
			name: "three items up ends left max",
			nn:   []int{3, 1, 2},
			want: []int{-1, 2, -1},
		},
		{
			name: "three items up ends right max",
			nn:   []int{2, 1, 3},
			want: []int{3, 3, -1},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := largest.Find(tc.nn)

			if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("largest.Find(%v) mismatch (-want, +got):\n%s", tc.nn, diff)
			}
		})
	}
}
