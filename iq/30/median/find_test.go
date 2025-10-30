// Copyright 2025 Samvel Khalatyan. All rights reserved.

package median_test

import (
	"testing"

	"github.com/skhal/lab/iq/30/median"
)

type test struct {
	name string
	nn   []int
	mm   []int
	want int
}

func TestFind(t *testing.T) {
	tests := []struct {
		name  string
		tests []test
	}{
		{
			name: "one array",
			tests: []test{
				{
					name: "first only size 1",
					nn:   []int{1},
					want: 1,
				},
				{
					name: "first only size 2",
					nn:   []int{1, 3},
					want: 2,
				},
				{
					name: "first only size 3",
					nn:   []int{1, 3, 5},
					want: 3,
				},
				{
					name: "second only size 1",
					mm:   []int{1},
					want: 1,
				},
				{
					name: "second only size 2",
					mm:   []int{1, 3},
					want: 2,
				},
				{
					name: "second only size 3",
					mm:   []int{1, 3, 5},
					want: 3,
				},
			},
		},
		{
			name: "first size one",
			tests: []test{
				{
					name: "second size one above",
					nn:   []int{1},
					mm:   []int{3},
					want: 2,
				},
				{
					name: "second size two above",
					nn:   []int{1},
					mm:   []int{3, 5},
					want: 3,
				},
				{
					name: "second size one below",
					nn:   []int{3},
					mm:   []int{1},
					want: 2,
				},
				{
					name: "second size two below",
					nn:   []int{5},
					mm:   []int{1, 3},
					want: 3,
				},
				{
					name: "second size two around",
					nn:   []int{3},
					mm:   []int{1, 5},
					want: 3,
				},
			},
		},
		{
			name: "first size two",
			tests: []test{
				{
					name: "second size one above",
					nn:   []int{1, 3},
					mm:   []int{5},
					want: 3,
				},
				{
					name: "second size two above",
					nn:   []int{1, 3},
					mm:   []int{5, 7},
					want: 4,
				},
				{
					name: "second size one below",
					nn:   []int{3, 5},
					mm:   []int{1},
					want: 3,
				},
				{
					name: "second size two below",
					nn:   []int{5, 7},
					mm:   []int{1, 3},
					want: 4,
				},
				{
					name: "second size two around",
					nn:   []int{3, 5},
					mm:   []int{1, 7},
					want: 4,
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) { testFind(t, tc.tests) })
	}
}

func testFind(t *testing.T, tests []test) {
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := median.Find(tc.nn, tc.mm)

			if tc.want != got {
				t.Errorf("median.Find(%v, %v) = %d; want %d", tc.nn, tc.mm, got, tc.want)
			}
		})
	}
}
