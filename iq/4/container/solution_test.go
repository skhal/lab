// Copyright 2025 Samvel Khalatyan. All rights reserved.

package container_test

import (
	"testing"

	"github.com/skhal/lab/iq/4/container"
)

func TestFind(t *testing.T) {
	tests := []struct {
		name string
		nn   []int
		want int
	}{
		{
			name: "empty",
		},
		{
			name: "one item",
			nn:   []int{1},
		},
		{
			name: "two items one zero",
			nn:   []int{0, 1},
		},
		{
			name: "two items same non-zero",
			nn:   []int{1, 1},
			want: 1,
		},
		{
			name: "two items non-zero different",
			nn:   []int{1, 2},
			want: 1,
		},
		{
			name: "three items flat",
			nn:   []int{1, 1, 1},
			want: 2,
		},
		{
			name: "three items ascending",
			nn:   []int{1, 2, 3},
			want: 2,
		},
		{
			name: "three items descending",
			nn:   []int{3, 2, 1},
			want: 2,
		},
		{
			name: "four items same edges",
			nn:   []int{1, 4, 3, 1},
			want: 3,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := container.Find(tc.nn)
			if want := container.Volume(tc.want); got != want {
				t.Errorf("container.Find(%v) = %d, want %d", tc.nn, got, want)
			}
		})
	}
}
