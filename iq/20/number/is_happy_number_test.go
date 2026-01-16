// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package number_test

import (
	"strconv"
	"testing"

	"github.com/skhal/lab/iq/20/number"
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
		{208, true},
		{931, true},
		// Not happy
		{2, false},
		{168, false},
		{936, false},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(strconv.Itoa(tc.n), func(t *testing.T) {
			got := number.IsHappyNumber(tc.n)

			if got != tc.want {
				t.Errorf("number.IsHappyNumber(%d) = %v; want %v", tc.n, got, tc.want)
			}
		})
	}
}
