// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package upper

func Find(nn []int, x int) (ub int, ok bool) {
	left := 0
	right := len(nn)
	for left < right {
		mid := left + (right-left)/2
		n := nn[mid]
		switch {
		case x < n:
			right = mid
		case x > n:
			left = mid + 1
		default:
			mid += 1
			if mid == len(nn) {
				return
			}
			return nn[mid], true
		}
	}
	if left == len(nn) {
		return
	}
	return nn[left], true
}
