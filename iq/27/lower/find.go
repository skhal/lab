// Copyright 2025 Samvel Khalatyan. All rights reserved.

package lower

func Find(nn []int, x int) (lb int, ok bool) {
	left := 0
	right := len(nn)
	for left < right {
		mid := left + (right-left)/2
		n := nn[mid]
		switch {
		case n < x:
			left = mid + 1
		case x < n:
			right = mid
		default:
			return n, true
		}
	}
	if right == len(nn) {
		return
	}
	return nn[right], true
}
