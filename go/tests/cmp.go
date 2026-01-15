// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tests

import (
	"math"

	"github.com/google/go-cmp/cmp"
)

// EquateFloat64 compares two float64 numbers and reports them equal if their
// relative change is less than or equal to relchg. See:
// https://en.wikipedia.org/wiki/Relative_change
func EquateFloat64(relchg float64) cmp.Option {
	relchg = math.Abs(relchg)
	return cmp.Comparer(func(x, y float64) bool {
		x, y = math.Min(x, y), math.Max(x, y)
		absDiff := math.Abs(y - x)
		absMean := math.Abs(x + absDiff/2)
		var relChange float64
		if absMean != 0 {
			relChange = absDiff / absMean
		}
		return relChange <= relchg
	})
}
