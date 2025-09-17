// Copyright 2025 Samvel Khalatyan. All rights reserved.

package threesum_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/iq/twopointer/threesum"
)

type Triplet = threesum.Triplet

var tests = []struct {
	name string
	nn   []int
	want []Triplet
}{
	// nagative tests
	{
		name: "empty",
	},
	{
		name: "one item",
		nn:   []int{1},
	},
	{
		name: "two items",
		nn:   []int{1, 2},
	},
	{
		name: "three positive items",
		nn:   []int{1, 2, 3},
	},
	{
		name: "three idential items",
		nn:   []int{1, 1, 1},
	},
	{
		name: "three negative items",
		nn:   []int{-1, -2, -3},
	},
	{
		name: "three idential negative items",
		nn:   []int{-1, -1, -1},
	},
	// positive tests
	// -- 3 items
	{
		name: "three items one triplet",
		nn:   []int{1, 2, -3},
		want: []Triplet{Triplet{-3, 1, 2}},
	},
	// -- 4 items
	{
		name: "four items",
		nn:   []int{1, 2, -3, 4},
		want: []Triplet{Triplet{-3, 1, 2}},
	},
	{
		name: "four items with duplicate low",
		nn:   []int{1, 2, -3, -3},
		want: []Triplet{Triplet{-3, 1, 2}},
	},
	{
		name: "four items with duplicate high",
		nn:   []int{1, 2, -3, 2},
		want: []Triplet{Triplet{-3, 1, 2}},
	},
	{
		name: "four items with duplicate middle",
		nn:   []int{1, 2, -3, 1},
		want: []Triplet{Triplet{-3, 1, 2}},
	},
	// -- 5 items
	{
		name: "five items one triplet",
		nn:   []int{2, 4, -6, 3, 4},
		want: []Triplet{Triplet{-6, 2, 4}},
	},
	{
		name: "five items two triplets",
		nn:   []int{2, 4, -6, 3, 3},
		want: []Triplet{Triplet{-6, 2, 4}, Triplet{-6, 3, 3}},
	},
}

func TestFind(t *testing.T) {
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := threesum.Find(tc.nn)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("threesum.Find(%v) mismatch (-want, +got):\n%s", tc.nn, diff)
			}
		})
	}
}

func TestFindWithOptimizations(t *testing.T) {
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := threesum.FindWithOptimizations(tc.nn)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("threesum.FindWithOptimizations(%v) mismatch (-want, +got):\n%s", tc.nn, diff)
			}
		})
	}
}
