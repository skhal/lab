// Copyright 2025 Samvel Khalatyan. All rights reserved.

package insert_test

import (
	"testing"

	"github.com/skhal/lab/iq/24/insert"
)

func TestFind(t *testing.T) {
	tests := []struct {
		name string
		nn   []int
		n    int
		want insert.Index
	}{
		{
			name: "empty",
			want: 0,
		},
		// size 1
		{
			name: "one item hit",
			nn:   []int{1},
			n:    1,
			want: 0,
		},
		{
			name: "one item insert left",
			nn:   []int{1},
			n:    0,
			want: 0,
		},
		{
			name: "one item insert right",
			nn:   []int{1},
			n:    2,
			want: 1,
		},
		// size 2
		{
			name: "two items hit first",
			nn:   []int{1, 3},
			n:    1,
			want: 0,
		},
		{
			name: "two items hit second",
			nn:   []int{1, 3},
			n:    3,
			want: 1,
		},
		{
			name: "two items insert left",
			nn:   []int{1, 3},
			n:    0,
			want: 0,
		},
		{
			name: "two items insert middle",
			nn:   []int{1, 3},
			n:    2,
			want: 1,
		},
		{
			name: "two items insert right",
			nn:   []int{1, 3},
			n:    4,
			want: 2,
		},
		// size 3
		{
			name: "three items hit first",
			nn:   []int{1, 3, 5},
			n:    1,
			want: 0,
		},
		{
			name: "three items hit second",
			nn:   []int{1, 3, 5},
			n:    3,
			want: 1,
		},
		{
			name: "three items hit third",
			nn:   []int{1, 3, 5},
			n:    5,
			want: 2,
		},
		{
			name: "three items insert left",
			nn:   []int{1, 3, 5},
			n:    0,
			want: 0,
		},
		{
			name: "three items insert between first and second",
			nn:   []int{1, 3, 5},
			n:    2,
			want: 1,
		},
		{
			name: "three items insert between second and third",
			nn:   []int{1, 3, 5},
			n:    4,
			want: 2,
		},
		{
			name: "three items insert right",
			nn:   []int{1, 3, 5},
			n:    6,
			want: 3,
		},
		// size 4
		{
			name: "four items hit first",
			nn:   []int{1, 3, 5, 7},
			n:    1,
			want: 0,
		},
		{
			name: "four items hit second",
			nn:   []int{1, 3, 5, 7},
			n:    3,
			want: 1,
		},
		{
			name: "four items hit third",
			nn:   []int{1, 3, 5, 7},
			n:    5,
			want: 2,
		},
		{
			name: "four items hit fourth",
			nn:   []int{1, 3, 5, 7},
			n:    7,
			want: 3,
		},
		{
			name: "four items insert left",
			nn:   []int{1, 3, 5, 7},
			n:    0,
			want: 0,
		},
		{
			name: "four items insert between first and second",
			nn:   []int{1, 3, 5, 7},
			n:    2,
			want: 1,
		},
		{
			name: "four items insert between second and third",
			nn:   []int{1, 3, 5, 7},
			n:    4,
			want: 2,
		},
		{
			name: "four items insert between third and fourth",
			nn:   []int{1, 3, 5, 7},
			n:    6,
			want: 3,
		},
		{
			name: "four items insert right",
			nn:   []int{1, 3, 5, 7},
			n:    8,
			want: 4,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := insert.FindInsertIndex(tc.nn, tc.n)

			if got != tc.want {
				t.Errorf("insert.FindInsertIndex(%v, %d) = %d; want %d", tc.nn, tc.n, got, tc.want)
			}
		})
	}
}
