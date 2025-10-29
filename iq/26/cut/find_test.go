// Copyright 2025 Samvel Khalatyan. All rights reserved.

package cut_test

import (
	"testing"

	"github.com/skhal/lab/iq/26/cut"
)

func TestFind(t *testing.T) {
	tests := []struct {
		name string
		nn   []int
		k    int
		want int
	}{
		{
			name: "one item",
			nn:   []int{2},
			k:    1,
			want: 1,
		},
		{
			name: "two items in order cut one",
			nn:   []int{2, 3},
			k:    1,
			want: 2,
		},
		{
			name: "two items out of order cut one",
			nn:   []int{3, 2},
			k:    2,
			want: 1,
		},
		{
			name: "two items in order cut two",
			nn:   []int{2, 3},
			k:    2,
			want: 1,
		},
		{
			name: "two items out of order cut two",
			nn:   []int{3, 2},
			k:    2,
			want: 1,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := cut.Find(tc.nn, tc.k)

			if got != tc.want {
				t.Errorf("cut.Find(%v, %d) = %d; want %d", tc.nn, tc.k, got, tc.want)
			}
		})
	}
}
