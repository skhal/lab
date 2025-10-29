// Copyright 2025 Samvel Khalatyan. All rights reserved.

package bounds_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/iq/25/bounds"
)

type testCase struct {
	name string
	nn   []int
	n    int
	want bounds.Bounds
}

func TestFind(t *testing.T) {
	tests := []struct {
		group string
		tests []testCase
	}{
		{
			group: "size 0",
			tests: []testCase{
				{
					name: "empty",
					want: bounds.BoundsError,
				},
			},
		},
		// size 1
		{
			group: "size 1",
			tests: []testCase{
				{
					name: "hit",
					nn:   []int{1},
					n:    1,
					want: bounds.Bounds{Left: 0, Right: 0},
				},
				{
					name: "miss left",
					nn:   []int{1},
					n:    0,
					want: bounds.BoundsError,
				},
				{
					name: "miss right",
					nn:   []int{1},
					n:    2,
					want: bounds.BoundsError,
				},
			},
		},
		// size 2
		{
			group: "size 2 no dups",
			tests: []testCase{
				{
					name: "hit first",
					nn:   []int{1, 3},
					n:    1,
					want: bounds.Bounds{Left: 0, Right: 0},
				},
				{
					name: "hit second",
					nn:   []int{1, 3},
					n:    3,
					want: bounds.Bounds{Left: 1, Right: 1},
				},
				{
					name: "miss left",
					nn:   []int{1, 3},
					n:    0,
					want: bounds.BoundsError,
				},
				{
					name: "miss between",
					nn:   []int{1, 3},
					n:    2,
					want: bounds.BoundsError,
				},
				{
					name: "miss right",
					nn:   []int{1, 3},
					n:    4,
					want: bounds.BoundsError,
				},
			},
		},
		{
			group: "size 2 dups",
			tests: []testCase{
				{
					name: "hit",
					nn:   []int{1, 1},
					n:    1,
					want: bounds.Bounds{Left: 0, Right: 1},
				},
				{
					name: "miss left",
					nn:   []int{1, 1},
					n:    0,
					want: bounds.BoundsError,
				},
				{
					name: "miss right",
					nn:   []int{1, 1},
					n:    2,
					want: bounds.BoundsError,
				},
			},
		},
		// size 3
		{
			group: "size 3 no dups",
			tests: []testCase{
				{
					name: "hit first",
					nn:   []int{1, 3, 5},
					n:    1,
					want: bounds.Bounds{Left: 0, Right: 0},
				},
				{
					name: "hit second",
					nn:   []int{1, 3, 5},
					n:    3,
					want: bounds.Bounds{Left: 1, Right: 1},
				},
				{
					name: "hit third",
					nn:   []int{1, 3, 5},
					n:    5,
					want: bounds.Bounds{Left: 2, Right: 2},
				},
				{
					name: "miss left",
					nn:   []int{1, 3, 5},
					n:    0,
					want: bounds.BoundsError,
				},
				{
					name: "miss between first and second",
					nn:   []int{1, 3, 5},
					n:    2,
					want: bounds.BoundsError,
				},
				{
					name: "miss between second and third",
					nn:   []int{1, 3, 5},
					n:    2,
					want: bounds.BoundsError,
				},
				{
					name: "miss right",
					nn:   []int{1, 3, 5},
					n:    6,
					want: bounds.BoundsError,
				},
			},
		},
		{
			group: "size 3 first two dups",
			tests: []testCase{
				{
					name: "hit first",
					nn:   []int{1, 1, 3},
					n:    1,
					want: bounds.Bounds{Left: 0, Right: 1},
				},
				{
					name: "hit third",
					nn:   []int{1, 1, 3},
					n:    3,
					want: bounds.Bounds{Left: 2, Right: 2},
				},
				{
					name: "miss left",
					nn:   []int{1, 1, 3},
					n:    0,
					want: bounds.BoundsError,
				},
				{
					name: "miss between first and second",
					nn:   []int{1, 1, 3},
					n:    2,
					want: bounds.BoundsError,
				},
				{
					name: "miss right",
					nn:   []int{1, 1, 3},
					n:    4,
					want: bounds.BoundsError,
				},
			},
		},
		{
			group: "size 3 last two dups",
			tests: []testCase{
				{
					name: "hit first",
					nn:   []int{1, 3, 3},
					n:    1,
					want: bounds.Bounds{Left: 0, Right: 0},
				},
				{
					name: "hit third",
					nn:   []int{1, 3, 3},
					n:    3,
					want: bounds.Bounds{Left: 1, Right: 2},
				},
				{
					name: "miss left",
					nn:   []int{1, 3, 3},
					n:    0,
					want: bounds.BoundsError,
				},
				{
					name: "miss between first and second",
					nn:   []int{1, 3, 3},
					n:    2,
					want: bounds.BoundsError,
				},
				{
					name: "miss right",
					nn:   []int{1, 3, 3},
					n:    4,
					want: bounds.BoundsError,
				},
			},
		},
		{
			group: "size 3 all dups",
			tests: []testCase{
				{
					name: "hit",
					nn:   []int{1, 1, 1},
					n:    1,
					want: bounds.Bounds{Left: 0, Right: 2},
				},
				{
					name: "miss left",
					nn:   []int{1, 1, 1},
					n:    0,
					want: bounds.BoundsError,
				},
				{
					name: "miss right",
					nn:   []int{1, 1, 1},
					n:    2,
					want: bounds.BoundsError,
				},
			},
		},
		// size 4
		{
			group: "size 4 no dups",
			tests: []testCase{
				{
					name: "hit first",
					nn:   []int{1, 3, 5, 7},
					n:    1,
					want: bounds.Bounds{Left: 0, Right: 0},
				},
				{
					name: "hit second",
					nn:   []int{1, 3, 5, 7},
					n:    3,
					want: bounds.Bounds{Left: 1, Right: 1},
				},
				{
					name: "hit third",
					nn:   []int{1, 3, 5, 7},
					n:    5,
					want: bounds.Bounds{Left: 2, Right: 2},
				},
				{
					name: "hit fourth",
					nn:   []int{1, 3, 5, 7},
					n:    7,
					want: bounds.Bounds{Left: 3, Right: 3},
				},
				{
					name: "miss left",
					nn:   []int{1, 3, 5, 7},
					n:    0,
					want: bounds.BoundsError,
				},
				{
					name: "miss between first and second",
					nn:   []int{1, 3, 5, 7},
					n:    2,
					want: bounds.BoundsError,
				},
				{
					name: "miss between second and third",
					nn:   []int{1, 3, 5, 7},
					n:    4,
					want: bounds.BoundsError,
				},
				{
					name: "miss between third and fourth",
					nn:   []int{1, 3, 5, 7},
					n:    6,
					want: bounds.BoundsError,
				},
				{
					name: "miss right",
					nn:   []int{1, 3, 5, 7},
					n:    8,
					want: bounds.BoundsError,
				},
			},
		},
		{
			group: "size 4 first two dups",
			tests: []testCase{
				{
					name: "hit first",
					nn:   []int{1, 1, 3, 5},
					n:    1,
					want: bounds.Bounds{Left: 0, Right: 1},
				},
				{
					name: "hit third",
					nn:   []int{1, 1, 3, 5},
					n:    3,
					want: bounds.Bounds{Left: 2, Right: 2},
				},
				{
					name: "hit fourth",
					nn:   []int{1, 1, 3, 5},
					n:    5,
					want: bounds.Bounds{Left: 3, Right: 3},
				},
				{
					name: "miss left",
					nn:   []int{1, 1, 3, 5},
					n:    0,
					want: bounds.BoundsError,
				},
				{
					name: "miss between second and third",
					nn:   []int{1, 1, 3, 5},
					n:    2,
					want: bounds.BoundsError,
				},
				{
					name: "miss between third and fourth",
					nn:   []int{1, 1, 3, 5},
					n:    4,
					want: bounds.BoundsError,
				},
				{
					name: "miss right",
					nn:   []int{1, 1, 3, 5},
					n:    6,
					want: bounds.BoundsError,
				},
			},
		},
		{
			group: "size 4 first three dups",
			tests: []testCase{
				{
					name: "hit first",
					nn:   []int{1, 1, 1, 3},
					n:    1,
					want: bounds.Bounds{Left: 0, Right: 2},
				},
				{
					name: "hit fourth",
					nn:   []int{1, 1, 1, 3},
					n:    3,
					want: bounds.Bounds{Left: 3, Right: 3},
				},
				{
					name: "miss left",
					nn:   []int{1, 1, 1, 3},
					n:    0,
					want: bounds.BoundsError,
				},
				{
					name: "miss between third and fourth",
					nn:   []int{1, 1, 1, 3},
					n:    2,
					want: bounds.BoundsError,
				},
				{
					name: "miss right",
					nn:   []int{1, 1, 1, 3},
					n:    4,
					want: bounds.BoundsError,
				},
			},
		},
		{
			group: "size 4 second two dups",
			tests: []testCase{
				{
					name: "hit first",
					nn:   []int{1, 3, 3, 5},
					n:    1,
					want: bounds.Bounds{Left: 0, Right: 0},
				},
				{
					name: "hit second",
					nn:   []int{1, 3, 3, 5},
					n:    3,
					want: bounds.Bounds{Left: 1, Right: 2},
				},
				{
					name: "hit fourth",
					nn:   []int{1, 3, 3, 5},
					n:    5,
					want: bounds.Bounds{Left: 3, Right: 3},
				},
				{
					name: "miss left",
					nn:   []int{1, 3, 3, 5},
					n:    0,
					want: bounds.BoundsError,
				},
				{
					name: "miss between first and second",
					nn:   []int{1, 3, 3, 5},
					n:    2,
					want: bounds.BoundsError,
				},
				{
					name: "miss between third and fourth",
					nn:   []int{1, 3, 3, 5},
					n:    4,
					want: bounds.BoundsError,
				},
				{
					name: "miss right",
					nn:   []int{1, 3, 3, 5},
					n:    6,
					want: bounds.BoundsError,
				},
			},
		},
		{
			group: "size 4 second three dups",
			tests: []testCase{
				{
					name: "hit first",
					nn:   []int{1, 3, 3, 3},
					n:    1,
					want: bounds.Bounds{Left: 0, Right: 0},
				},
				{
					name: "hit second",
					nn:   []int{1, 3, 3, 3},
					n:    3,
					want: bounds.Bounds{Left: 1, Right: 3},
				},
				{
					name: "miss left",
					nn:   []int{1, 3, 3, 3},
					n:    0,
					want: bounds.BoundsError,
				},
				{
					name: "miss between first and second",
					nn:   []int{1, 3, 3, 3},
					n:    2,
					want: bounds.BoundsError,
				},
				{
					name: "miss right",
					nn:   []int{1, 3, 3, 3},
					n:    4,
					want: bounds.BoundsError,
				},
			},
		},
		{
			group: "size 4 last two dups",
			tests: []testCase{
				{
					name: "hit first",
					nn:   []int{1, 3, 5, 5},
					n:    1,
					want: bounds.Bounds{Left: 0, Right: 0},
				},
				{
					name: "hit second",
					nn:   []int{1, 3, 5, 5},
					n:    3,
					want: bounds.Bounds{Left: 1, Right: 1},
				},
				{
					name: "hit third",
					nn:   []int{1, 3, 5, 5},
					n:    5,
					want: bounds.Bounds{Left: 2, Right: 3},
				},
				{
					name: "miss left",
					nn:   []int{1, 3, 5, 5},
					n:    0,
					want: bounds.BoundsError,
				},
				{
					name: "miss between first and second",
					nn:   []int{1, 3, 5, 5},
					n:    2,
					want: bounds.BoundsError,
				},
				{
					name: "miss right",
					nn:   []int{1, 3, 5, 5},
					n:    6,
					want: bounds.BoundsError,
				},
			},
		},
		{
			group: "size 4 last three dups",
			tests: []testCase{
				{
					name: "hit first",
					nn:   []int{1, 3, 3, 3},
					n:    1,
					want: bounds.Bounds{Left: 0, Right: 0},
				},
				{
					name: "hit second",
					nn:   []int{1, 3, 3, 3},
					n:    3,
					want: bounds.Bounds{Left: 1, Right: 3},
				},
				{
					name: "miss left",
					nn:   []int{1, 3, 3, 3},
					n:    0,
					want: bounds.BoundsError,
				},
				{
					name: "miss between first and second",
					nn:   []int{1, 3, 3, 3},
					n:    2,
					want: bounds.BoundsError,
				},
				{
					name: "miss right",
					nn:   []int{1, 3, 3, 3},
					n:    4,
					want: bounds.BoundsError,
				},
			},
		},
		{
			group: "size 4 all dups",
			tests: []testCase{
				{
					name: "hit",
					nn:   []int{1, 1, 1},
					n:    1,
					want: bounds.Bounds{Left: 0, Right: 2},
				},
				{
					name: "miss left",
					nn:   []int{1, 1, 1},
					n:    0,
					want: bounds.BoundsError,
				},
				{
					name: "miss right",
					nn:   []int{1, 1, 1},
					n:    2,
					want: bounds.BoundsError,
				},
			},
		},
	}
	for _, tg := range tests {
		tg := tg
		t.Run(tg.group, func(t *testing.T) { testFind(t, tg.tests) })
	}
}

func testFind(t *testing.T, tests []testCase) {
	t.Helper()
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := bounds.Find(tc.nn, tc.n)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("bounds.Find(%v, %d) mismatch (-want, +got):\n%s", tc.nn, tc.n, diff)
			}
		})
	}
}
