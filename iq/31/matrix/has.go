// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package matrix

type M map[int][]int

func Has(m M, n int) bool {
	if len(m) == 0 {
		return false
	}
	left := 0
	right := len(m) * len(m[0])
	for left < right {
		mid := left + (right-left)/2
		x := getValue(m, mid)
		switch {
		case n < x:
			right = mid
		case x < n:
			left = mid + 1
		default:
			return true
		}
	}
	return false
}

func getValue(m M, index int) int {
	r := index / len(m[0])
	c := index % len(m[0])
	return m[r][c]
}
