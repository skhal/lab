// Copyright 2025 Samvel Khalatyan. All rights reserved.

package cut

import "slices"

func Find(nn []int, k int) int {
	left := 0
	right := slices.Max(nn)
	for left < right {
		mid := left + (right-left)/2
		x := sumAbove(nn, mid+1)
		switch {
		case x < k:
			right = mid
		default:
			left = mid + 1
		}
	}
	return left
}

func sumAbove(nn []int, c int) int {
	sum := 0
	for _, n := range nn {
		if n > c {
			sum += n - c
		}
	}
	return sum
}
