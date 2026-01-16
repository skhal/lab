// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package maxima_test

import (
	"testing"

	"github.com/skhal/lab/iq/32/maxima"
)

func TestFind(t *testing.T) {
	tests := []struct {
		name      string
		nn        []int
		wantAnyOf []int
	}{
		{
			name:      "empty",
			wantAnyOf: []int{0},
		},
		{
			name:      "size one",
			nn:        []int{1},
			wantAnyOf: []int{1},
		},
		// size 2
		{
			name:      "size two ascending",
			nn:        []int{1, 2},
			wantAnyOf: []int{2},
		},
		{
			name:      "size two descending",
			nn:        []int{2, 1},
			wantAnyOf: []int{2},
		},
		// size 3
		{
			name:      "size three ascending",
			nn:        []int{1, 2, 3},
			wantAnyOf: []int{3},
		},
		{
			name:      "size three descending",
			nn:        []int{3, 2, 1},
			wantAnyOf: []int{3},
		},
		{
			name:      "size three one max",
			nn:        []int{1, 3, 2},
			wantAnyOf: []int{3},
		},
		{
			name:      "size three two max",
			nn:        []int{2, 1, 3},
			wantAnyOf: []int{2, 3},
		},
		// size 4
		{
			name:      "size four ascending",
			nn:        []int{1, 2, 3, 4},
			wantAnyOf: []int{4},
		},
		{
			name:      "size four descending",
			nn:        []int{4, 3, 2, 1},
			wantAnyOf: []int{4},
		},
		{
			name:      "size four one max",
			nn:        []int{1, 2, 4, 3},
			wantAnyOf: []int{4},
		},
		{
			name:      "size four two max",
			nn:        []int{1, 4, 2, 3},
			wantAnyOf: []int{3, 4},
		},
		{
			name:      "size four two max at begin and end",
			nn:        []int{4, 1, 2, 3},
			wantAnyOf: []int{3, 4},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := maxima.Find(tc.nn)

			if !anyOf(t, tc.wantAnyOf, got) {
				t.Errorf("maxima.Find(%v) = %d; wany any of %v", tc.nn, got, tc.wantAnyOf)
			}
		})
	}
}

func anyOf(t *testing.T, nn []int, x int) bool {
	t.Helper()
	for _, n := range nn {
		if n == x {
			return true
		}
	}
	return false
}
