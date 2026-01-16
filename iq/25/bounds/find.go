// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bounds

const indexError = -1

var BoundsError = Bounds{indexError, indexError}

type Bounds struct {
	Left  int
	Right int
}

func Find(nn []int, n int) Bounds {
	left := find(nn, n, findLeft)
	if left == -1 {
		return BoundsError
	}
	right := find(nn, n, findRight)
	if right == -1 {
		return BoundsError
	}
	return Bounds{Left: left, Right: right}
}

func find(nn []int, n int, f func(nn []int, n int, left, mid, right int) int) int {
	for left, right := 0, len(nn); left < right; {
		mid := left + (right-left)/2
		switch x := nn[mid]; {
		case n < x:
			right = mid
		case n > x:
			left = mid + 1
		default:
			return f(nn, n, left, mid, right)
		}
	}
	return indexError
}

func findLeft(nn []int, n int, left, mid, right int) int {
	i := find(nn[left:mid], n, findLeft)
	if i == indexError {
		return mid
	}
	return left + i
}

func findRight(nn []int, n int, left, mid, right int) int {
	i := find(nn[mid+1:right], n, findRight)
	if i == indexError {
		return mid
	}
	return mid + 1 + i
}
