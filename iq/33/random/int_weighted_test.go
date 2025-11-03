// Copyright 2025 Samvel Khalatyan. All rights reserved.

package random_test

import (
	"testing"

	"github.com/skhal/lab/iq/33/random"
)

type want struct {
	n    int
	wmax int
}

type test struct {
	name string
	ww   []int
	rand int
	want want
}

func TestIntWeighted(t *testing.T) {
	tests := []struct {
		name  string
		tests []test
	}{
		// Size 1
		{
			name: "size one",
			tests: []test{
				{
					name: "wmax 1 rand 0",
					ww:   []int{1},
					rand: 0,
					want: want{
						n:    0,
						wmax: 1,
					},
				},
				{
					name: "wmax 2 rand 0",
					ww:   []int{2},
					rand: 0,
					want: want{
						n:    0,
						wmax: 2,
					},
				},
				{
					name: "wmax 2 rand 1",
					ww:   []int{2},
					rand: 1,
					want: want{
						n:    0,
						wmax: 2,
					},
				},
			},
		},
		// Size 2
		{
			name: "size two",
			tests: []test{
				{
					name: "wmax 2 rand 0",
					ww:   []int{1, 1},
					rand: 0,
					want: want{
						n:    0,
						wmax: 2,
					},
				},
				{
					name: "wmax 2 rand 1",
					ww:   []int{1, 1},
					rand: 1,
					want: want{
						n:    1,
						wmax: 2,
					},
				},
				{
					name: "wmax 3 ww 1 2 rand 0",
					ww:   []int{1, 2},
					rand: 0,
					want: want{
						n:    0,
						wmax: 3,
					},
				},
				{
					name: "wmax 3 ww 1 2 rand 1",
					ww:   []int{1, 2},
					rand: 1,
					want: want{
						n:    1,
						wmax: 3,
					},
				},
				{
					name: "wmax 3 ww 1 2 rand 2",
					ww:   []int{1, 2},
					rand: 2,
					want: want{
						n:    1,
						wmax: 3,
					},
				},
				{
					name: "wmax 3 ww 2 1 rand 0",
					ww:   []int{2, 1},
					rand: 0,
					want: want{
						n:    0,
						wmax: 3,
					},
				},
				{
					name: "wmax 3 ww 2 1 rand 1",
					ww:   []int{2, 1},
					rand: 1,
					want: want{
						n:    0,
						wmax: 3,
					},
				},
				{
					name: "wmax 3 ww 2 1 rand 2",
					ww:   []int{2, 1},
					rand: 2,
					want: want{
						n:    1,
						wmax: 3,
					},
				},
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) { testIntWeighted(t, tc.tests) })
	}
}

func testIntWeighted(t *testing.T, tests []test) {
	t.Helper()
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var gotWMax int
			got := random.IntWeighted(tc.ww, func(n int) int { gotWMax = n; return tc.rand })

			if gotWMax != tc.want.wmax {
				t.Errorf("random.IntWeighted(%v, func(wmax int) int {})  unexpected wmax %d; want %d", tc.ww, gotWMax, tc.want.wmax)
			}
			if got != tc.want.n {
				t.Errorf("random.IntWeighted(%v, func(int) int { return %d }) = %d; want %d", tc.ww, tc.rand, got, tc.want.n)
			}
		})
	}
}
