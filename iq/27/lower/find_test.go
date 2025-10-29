// Copyright 2025 Samvel Khalatyan. All rights reserved.

package lower_test

import (
	"testing"

	"github.com/skhal/lab/iq/27/lower"
)

func TestFind(t *testing.T) {
	tests := []struct {
		name   string
		nn     []int
		x      int
		want   int
		wantOk bool
	}{
		// size 1
		{
			name:   "size one below",
			nn:     []int{1},
			x:      0,
			want:   1,
			wantOk: true,
		},
		{
			name:   "size one match",
			nn:     []int{1},
			x:      1,
			want:   1,
			wantOk: true,
		},
		{
			name: "size one above",
			nn:   []int{1},
			x:    2,
		},
		// size 2
		{
			name:   "size two below",
			nn:     []int{1, 3},
			x:      0,
			want:   1,
			wantOk: true,
		},
		{
			name:   "size two match first",
			nn:     []int{1, 3},
			x:      1,
			want:   1,
			wantOk: true,
		},
		{
			name:   "size two below second",
			nn:     []int{1, 3},
			x:      2,
			want:   3,
			wantOk: true,
		},
		{
			name:   "size two match second",
			nn:     []int{1, 3},
			x:      3,
			want:   3,
			wantOk: true,
		},
		{
			name: "size two above",
			nn:   []int{1, 3},
			x:    4,
		},
		// size 3
		{
			name:   "size three below",
			nn:     []int{1, 3, 5},
			x:      0,
			want:   1,
			wantOk: true,
		},
		{
			name:   "size three match first",
			nn:     []int{1, 3, 5},
			x:      1,
			want:   1,
			wantOk: true,
		},
		{
			name:   "size three below second",
			nn:     []int{1, 3, 5},
			x:      2,
			want:   3,
			wantOk: true,
		},
		{
			name:   "size three match second",
			nn:     []int{1, 3, 5},
			x:      3,
			want:   3,
			wantOk: true,
		},
		{
			name:   "size three below third",
			nn:     []int{1, 3, 5},
			x:      4,
			want:   5,
			wantOk: true,
		},
		{
			name:   "size three match third",
			nn:     []int{1, 3, 5},
			x:      5,
			want:   5,
			wantOk: true,
		},
		{
			name: "size three above",
			nn:   []int{1, 3, 5},
			x:    6,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, gotOk := lower.Find(tc.nn, tc.x)

			if gotOk != tc.wantOk {
				t.Errorf("lower.Find(%v, %d) = _, %v; want %v", tc.nn, tc.x, gotOk, tc.wantOk)
			}
			if got != tc.want {
				t.Errorf("lower.Find(%v, %d) = %d, _; want %d", tc.nn, tc.x, got, tc.want)
			}
		})
	}
}
