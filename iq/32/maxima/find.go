// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package maxima

import "math"

func Find(nn []int) int {
	left := 0
	right := len(nn)
	for left < right {
		mid := left + (right-left)/2
		switch {
		case shouldMoveLeft(nn, mid):
			right = mid
		case shouldMoveRight(nn, mid):
			left = mid + 1
		default:
			return nn[mid]
		}
	}
	if left < len(nn) {
		return nn[left]
	}
	return 0
}

func shouldMoveLeft(nn []int, idx int) bool {
	return prev(nn, idx) > nn[idx]
}

func shouldMoveRight(nn []int, idx int) bool {
	return nn[idx] < next(nn, idx)
}

func prev(nn []int, idx int) int {
	if idx == 0 {
		return math.MinInt
	}
	return nn[idx-1]
}

func next(nn []int, idx int) int {
	idx += 1
	if idx == len(nn) {
		return math.MinInt
	}
	return nn[idx]
}
