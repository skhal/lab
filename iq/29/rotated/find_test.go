// Copyright 2025 Samvel Khalatyan. All rights reserved.

package rotated_test

import (
	"testing"

	"github.com/skhal/lab/iq/29/rotated"
)

type test struct {
	name string
	nn   []int
	x    int
	want rotated.Index
}

func TestFind(t *testing.T) {
	tests := []struct {
		group string
		tests []test
	}{
		{
			group: "empty",
			tests: []test{
				{
					name: "miss",
					x:    1,
					want: rotated.IndexError,
				},
			},
		},
		{
			group: "size one",
			tests: []test{
				{
					name: "hit",
					nn:   []int{1},
					x:    1,
					want: 0,
				},
				{
					name: "miss",
					nn:   []int{1},
					x:    2,
					want: rotated.IndexError,
				},
			},
		},
		{
			group: "size two no rotation",
			tests: []test{
				{
					name: "hit first",
					nn:   []int{1, 3},
					x:    1,
					want: 0,
				},
				{
					name: "hit second",
					nn:   []int{1, 3},
					x:    3,
					want: 1,
				},
				{
					name: "miss below",
					nn:   []int{1, 3},
					x:    0,
					want: rotated.IndexError,
				},
				{
					name: "miss between first and second",
					nn:   []int{1, 3},
					x:    2,
					want: rotated.IndexError,
				},
				{
					name: "miss above",
					nn:   []int{1, 3},
					x:    4,
					want: rotated.IndexError,
				},
			},
		},
		{
			group: "size two rotate one",
			tests: []test{
				{
					name: "hit first",
					nn:   []int{3, 1},
					x:    3,
					want: 0,
				},
				{
					name: "hit second",
					nn:   []int{3, 1},
					x:    1,
					want: 1,
				},
				{
					name: "miss below",
					nn:   []int{3, 1},
					x:    0,
					want: rotated.IndexError,
				},
				{
					name: "miss between first and second",
					nn:   []int{3, 1},
					x:    2,
					want: rotated.IndexError,
				},
				{
					name: "miss above",
					nn:   []int{3, 1},
					x:    4,
					want: rotated.IndexError,
				},
			},
		},
		{
			group: "size three no rotation",
			tests: []test{
				{
					name: "hit first",
					nn:   []int{1, 3, 5},
					x:    1,
					want: 0,
				},
				{
					name: "hit second",
					nn:   []int{1, 3, 5},
					x:    3,
					want: 1,
				},
				{
					name: "hit third",
					nn:   []int{1, 3, 5},
					x:    5,
					want: 2,
				},
				{
					name: "miss below",
					nn:   []int{1, 3, 5},
					x:    0,
					want: rotated.IndexError,
				},
				{
					name: "miss between first and second",
					nn:   []int{1, 3, 5},
					x:    2,
					want: rotated.IndexError,
				},
				{
					name: "miss between second and third",
					nn:   []int{1, 3, 5},
					x:    4,
					want: rotated.IndexError,
				},
				{
					name: "miss above",
					nn:   []int{1, 3, 5},
					x:    6,
					want: rotated.IndexError,
				},
			},
		},
		{
			group: "size three rotate one",
			tests: []test{
				{
					name: "hit first",
					nn:   []int{5, 1, 3},
					x:    5,
					want: 0,
				},
				{
					name: "hit second",
					nn:   []int{5, 1, 3},
					x:    1,
					want: 1,
				},
				{
					name: "hit third",
					nn:   []int{5, 1, 3},
					x:    3,
					want: 2,
				},
				{
					name: "miss below",
					nn:   []int{5, 1, 3},
					x:    0,
					want: rotated.IndexError,
				},
				{
					name: "miss between first and second",
					nn:   []int{5, 1, 3},
					x:    4,
					want: rotated.IndexError,
				},
				{
					name: "miss between second and third",
					nn:   []int{5, 1, 3},
					x:    2,
					want: rotated.IndexError,
				},
				{
					name: "miss above",
					nn:   []int{5, 1, 3},
					x:    6,
					want: rotated.IndexError,
				},
			},
		},
		{
			group: "size three rotate two",
			tests: []test{
				{
					name: "hit first",
					nn:   []int{3, 5, 1},
					x:    3,
					want: 0,
				},
				{
					name: "hit second",
					nn:   []int{3, 5, 1},
					x:    5,
					want: 1,
				},
				{
					name: "hit third",
					nn:   []int{3, 5, 1},
					x:    1,
					want: 2,
				},
				{
					name: "miss below",
					nn:   []int{3, 5, 1},
					x:    0,
					want: rotated.IndexError,
				},
				{
					name: "miss between first and second",
					nn:   []int{3, 5, 1},
					x:    4,
					want: rotated.IndexError,
				},
				{
					name: "miss between second and third",
					nn:   []int{3, 5, 1},
					x:    2,
					want: rotated.IndexError,
				},
				{
					name: "miss above",
					nn:   []int{3, 5, 1},
					x:    6,
					want: rotated.IndexError,
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.group, func(t *testing.T) { testFind(t, tc.tests) })
	}
}

func testFind(t *testing.T, tests []test) {
	t.Helper()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := rotated.Find(tc.nn, tc.x)

			if got != tc.want {
				t.Errorf("rotated.Find(%v, %d) = %d; want %d", tc.nn, tc.x, got, tc.want)
			}
		})
	}
}
