// Copyright 2025 Samvel Khalatyan. All rights reserved.

package chain_test

import (
	"slices"
	"testing"

	"github.com/skhal/lab/iq/10/chain"
)

func contains(t *testing.T, want []chain.C, got chain.C) bool {
	t.Helper()
	if len(want) == 0 && len(got) == 0 {
		return true
	}
	return slices.ContainsFunc(want, func(nn chain.C) bool {
		if len(nn) != len(got) {
			return false
		}
		for i, n := range nn {
			if got[i] != n {
				return false
			}
		}
		return true
	})
}

func TestFind(t *testing.T) {
	tests := []struct {
		name string
		nn   []int
		want []chain.C // different opions
	}{
		{
			name: "empty",
		},
		{
			name: "one item",
			nn:   []int{1},
			want: []chain.C{{1}},
		},
		{
			name: "two items one chain",
			nn:   []int{1, 2},
			want: []chain.C{{1, 2}},
		},
		{
			name: "two items reversed one chain",
			nn:   []int{2, 1},
			want: []chain.C{{1, 2}},
		},
		{
			name: "two items two chains",
			nn:   []int{1, 3},
			want: []chain.C{{1}, {3}},
		},
		{
			name: "two chains same size",
			nn:   []int{1, 5, 2, 4},
			want: []chain.C{{1, 2}, {4, 5}},
		},
		{
			name: "two chains different size",
			nn:   []int{1, 7, 5, 2, 6},
			want: []chain.C{{5, 6, 7}},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := chain.Find(tc.nn)

			if ok := contains(t, tc.want, got); !ok {
				t.Errorf("chain.Find(%v) = %v; want one of %v", tc.nn, got, tc.want)
			}
		})
	}
}
