// Copyright 2025 Samvel Khalatyan. All rights reserved.

package shiftzeros_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/iq/twopointer/shiftzeros"
)

func TestShift(t *testing.T) {
	tests := []struct {
		name string
		nn   []int
		want []int
	}{
		{
			name: "empty",
		},
		// 1 item
		{
			name: "one item not zero",
			nn:   []int{1},
			want: []int{1},
		},
		{
			name: "one item zero",
			nn:   []int{0},
			want: []int{0},
		},
		// 2 items
		{
			name: "two items not zero",
			nn:   []int{1, 2},
			want: []int{1, 2},
		},
		{
			name: "two items first zero",
			nn:   []int{0, 1},
			want: []int{1, 0},
		},
		{
			name: "two items second zero",
			nn:   []int{1, 0},
			want: []int{1, 0},
		},
		{
			name: "two items all zero",
			nn:   []int{0, 0},
			want: []int{0, 0},
		},
		// 3 items
		{
			name: "three items not zero",
			nn:   []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			name: "three items first zero",
			nn:   []int{0, 2, 3},
			want: []int{2, 3, 0},
		},
		{
			name: "three items second zero",
			nn:   []int{1, 0, 3},
			want: []int{1, 3, 0},
		},
		{
			name: "three items third zero",
			nn:   []int{1, 2, 0},
			want: []int{1, 2, 0},
		},
		{
			name: "three items first and second zero",
			nn:   []int{0, 0, 3},
			want: []int{3, 0, 0},
		},
		{
			name: "three items first and third zero",
			nn:   []int{0, 2, 0},
			want: []int{2, 0, 0},
		},
		{
			name: "three items second and third zero",
			nn:   []int{1, 0, 0},
			want: []int{1, 0, 0},
		},
		{
			name: "three items all zero",
			nn:   []int{1, 0, 0},
			want: []int{1, 0, 0},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			nn := tc.nn[:] // for reporting
			shiftzeros.Shift(nn)
			if diff := cmp.Diff(tc.want, nn); diff != "" {
				t.Errorf("shitzeros.Shift(%v) mismatch (-want, +got):\n%s", nn, diff)
			}
		})
	}
}
