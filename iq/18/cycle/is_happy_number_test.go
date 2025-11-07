// Copyright 2025 Samvel Khalatyan. All rights reserved.

package cycle_test

import (
	"strconv"
	"testing"

	"github.com/skhal/lab/iq/18/cycle"
)

func TestIsHappyNumber(t *testing.T) {
	tests := []struct {
		n    int
		want bool
	}{
		// Happy numbers: https://en.wikipedia.org/wiki/Happy_number
		{1, true},
		{7, true},
		{10, true},
		{188, true},
		{931, true},
		// Not happy
		{2, false},
		{168, false},
		{936, false},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(strconv.Itoa(tc.n), func(t *testing.T) {
			got := cycle.IsHappyNumber(tc.n)

			if got != tc.want {
				t.Errorf("cycle.IsHappyNubmer(%d) = %v; want %v", tc.n, got, tc.want)
			}
		})
	}
}
