// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

import "math"

// floatEqPrecision defines max per-cent difference allowed to treat two values
// equal.
const floatEqPrecision = 0.001 // 0.1%

// Equal is approximate equality of two floating numbers. The two number are
// considered equal, if relative absolute difference between the two is within
// set precision [floatEqPrecision].
func Equal(x, y float64) bool {
	x, y = math.Min(x, y), math.Max(x, y)
	absDiff := math.Abs(y - x)
	absMean := math.Abs(x + absDiff/2)
	var relChange float64
	if absMean != 0 {
		relChange = absDiff / absMean
	}
	return relChange <= floatEqPrecision
}
