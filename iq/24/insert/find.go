// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package insert

type Index int

func FindInsertIndex(nn []int, n int) Index {
	left := Index(0)
	right := Index(len(nn))
	for left < right {
		mid := left + (right-left)/2
		switch x := nn[mid]; {
		case n < x:
			right = mid
		case n > x:
			left = mid + 1
		default:
			return mid
		}
	}
	return left
}
