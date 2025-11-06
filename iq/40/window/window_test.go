// Copyright 2025 Samvel Khalatyan. All rights reserved.

package window_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/skhal/lab/iq/40/window"
)

func ExampleW() {
	nn := []int{1, 4, 2, 3, 1, 2}
	mm := make([]int, 0, len(nn))
	for w, _ := window.New(nn, 3); w.Slide(); {
		mm = append(mm, w.Max())
	}
	fmt.Printf("%v\n", mm)
	// Output:
	// [1 4 4 4 3 3]
}

type test struct {
	name string
	nn   []int
	want []int
}

func TestW(t *testing.T) {
	tests := []struct {
		size  int
		tests []test
	}{
		{
			size: 1,
			tests: []test{
				{
					name: "one item",
					nn:   []int{1},
					want: []int{1},
				},
				{
					name: "two items ascend",
					nn:   []int{1, 2},
					want: []int{1, 2},
				},
				{
					name: "two items descend",
					nn:   []int{2, 1},
					want: []int{2, 1},
				},
				{
					name: "three items ascend",
					nn:   []int{1, 2, 3},
					want: []int{1, 2, 3},
				},
				{
					name: "three items descend",
					nn:   []int{3, 2, 1},
					want: []int{3, 2, 1},
				},
				{
					name: "three items up and down",
					nn:   []int{1, 3, 2},
					want: []int{1, 3, 2},
				},
				{
					name: "three items down and up",
					nn:   []int{3, 1, 2},
					want: []int{3, 1, 2},
				},
			},
		},
		{
			size: 2,
			tests: []test{
				{
					name: "one item",
					nn:   []int{1},
					want: []int{1},
				},
				{
					name: "two items ascend",
					nn:   []int{1, 2},
					want: []int{1, 2},
				},
				{
					name: "two items descend",
					nn:   []int{2, 1},
					want: []int{2, 2},
				},
				{
					name: "three items ascend",
					nn:   []int{1, 2, 3},
					want: []int{1, 2, 3},
				},
				{
					name: "three items descend",
					nn:   []int{3, 2, 1},
					want: []int{3, 3, 2},
				},
				{
					name: "three items up and down",
					nn:   []int{1, 3, 2},
					want: []int{1, 3, 3},
				},
				{
					name: "three items down and up",
					nn:   []int{3, 1, 2},
					want: []int{3, 3, 2},
				},
			},
		},
		{
			size: 3,
			tests: []test{
				{
					name: "one item",
					nn:   []int{1},
					want: []int{1},
				},
				{
					name: "two items ascend",
					nn:   []int{1, 2},
					want: []int{1, 2},
				},
				{
					name: "two items descend",
					nn:   []int{2, 1},
					want: []int{2, 2},
				},
				{
					name: "three items ascend",
					nn:   []int{1, 2, 3},
					want: []int{1, 2, 3},
				},
				{
					name: "three items descend",
					nn:   []int{3, 2, 1},
					want: []int{3, 3, 3},
				},
				{
					name: "three items up and down",
					nn:   []int{1, 3, 2},
					want: []int{1, 3, 3},
				},
				{
					name: "three items down and up",
					nn:   []int{3, 1, 2},
					want: []int{3, 3, 3},
				},
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("size=%d", tc.size), func(t *testing.T) {
			testW(t, tc.size, tc.tests)
		})
	}
}

func testW(t *testing.T, size int, tests []test) {
	t.Helper()
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			w := mustCreateWindow(t, tc.nn, size)
			got := make([]int, 0, len(tc.nn))

			for w.Slide() {
				got = append(got, w.Max())
			}

			if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("window.New(%v, %d) mismatch (-want, +got):\n%s", tc.nn, size, diff)
			}
		})
	}
}

func mustCreateWindow(t *testing.T, nn []int, size int) *window.W {
	t.Helper()
	w, err := window.New(nn, size)
	if err != nil {
		t.Fatalf("window.New(%d) failed with error: %s", size, err)
	}
	return w
}
