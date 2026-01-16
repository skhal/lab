// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package twosum_test

import (
	"sort"
	"testing"

	"github.com/skhal/lab/iq/7/twosum"
)

func eqSlices(nn, mm []int) bool {
	if len(nn) != len(mm) {
		return false
	}
	for i, n := range nn {
		if n != mm[i] {
			return false
		}
	}
	return true
}

func valueIn(nnn [][]int, mm []int) bool {
	if len(nnn) == 0 && len(mm) == 0 {
		return true
	}
	for _, nn := range nnn {
		if eqSlices(nn, mm) {
			return true
		}
	}
	return false
}

func TestFind(t *testing.T) {
	tests := []struct {
		name string
		nn   []int
		x    int
		want [][]int // all possible combinations
	}{
		{
			name: "empty",
			x:    4,
		},
		{
			name: "one item",
			nn:   []int{4},
			x:    4,
		},
		{
			name: "two items",
			nn:   []int{1, 3},
			x:    4,
			want: [][]int{{0, 1}},
		},
		{
			name: "two items not found",
			nn:   []int{1, 4},
			x:    4,
		},
		/*
			{
				name: "three items",
				nn:   []int{3, 0, 1},
				x:    4,
				want: [][]int{{0, 2}},
			},
			{
				name: "three items not found",
				nn:   []int{3, 0, 2},
				x:    4,
			},
			{
				name: "three items multiple matches",
				nn:   []int{1, 3, 2, 4},
				x:    5,
				want: [][]int{{1, 2}, {0, 3}},
			},
		*/
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := twosum.Find(tc.nn, tc.x)

			gotSorted := got[:]
			sort.Ints(gotSorted)
			if !valueIn(tc.want, gotSorted) {
				t.Errorf("twosum.Find(%v, %d) = %v (sorted %v); want one of %v", tc.nn, tc.x, got, gotSorted, tc.want)
			}
		})
	}
}
