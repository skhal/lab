// Copyright 2025 Samvel Khalatyan. All rights reserved.

package chain_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/iq/mapset/chain"
)

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
			want: []int{1},
		},
		{
			name: "two items one chain",
			nn:   []int{1, 2},
			want: []int{1, 2},
		},
		{
			name: "two items reversed one chain",
			nn:   []int{2, 1},
			want: []int{1, 2},
		},
		{
			name: "two items two chains",
			nn:   []int{1, 3},
			want: []int{1},
		},
		{
			name: "two chains same size",
			nn:   []int{1, 5, 2, 4},
			want: []int{1, 2},
		},
		{
			name: "two chains different size",
			nn:   []int{1, 7, 5, 2, 6},
			want: []int{5, 6, 7},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := chain.Find(tc.nn)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("chain.Find(%v) mismatch (-want, +got):\n%s", tc.nn, diff)
			}
		})
	}
}
