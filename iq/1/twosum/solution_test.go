// Copyright 2025 Samvel Khalatyan. All rights reserved.

package twosum_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/iq/1/twosum"
)

func TestFind(t *testing.T) {
	tests := []struct {
		name string
		nn   []int
		n    int
		want []int
	}{
		{
			name: "empty input",
			n:    1,
		},
		{
			name: "one element input",
			nn:   []int{1},
			n:    1,
		},
		{
			name: "no match",
			nn:   []int{1, 2},
			n:    4,
		},
		{
			name: "match",
			nn:   []int{1, 2},
			n:    3,
			want: []int{0, 1},
		},
		{
			name: "first match",
			nn:   []int{1, 1, 2},
			n:    3,
			want: []int{0, 2},
		},
		{
			name: "move first index",
			nn:   []int{1, 2, 3, 4},
			n:    6,
			want: []int{1, 3},
		},
		{
			name: "move second index",
			nn:   []int{1, 2, 3, 4},
			n:    4,
			want: []int{0, 2},
		},
		{
			name: "negative values",
			nn:   []int{-1, 2, 3, 4},
			n:    2,
			want: []int{0, 2},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := twosum.Find(tc.nn, tc.n)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("twosum.Find(%v, %d) mismatch (-want +got):\n%s", tc.nn, tc.n, diff)
			}
		})
	}
}
