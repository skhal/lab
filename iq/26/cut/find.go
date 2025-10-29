// Copyright 2025 Samvel Khalatyan. All rights reserved.

package cut

import "slices"

func Find(nn []int, k int) int {
	left := 0
	right := slices.Max(nn)
	for left < right {
		mid := left + (right-left)/2 + 1
		if x := sumAbove(nn, mid); x < k {
			right = mid - 1
			continue
		}
		left = mid
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
