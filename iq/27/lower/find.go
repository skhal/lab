// Copyright 2025 Samvel Khalatyan. All rights reserved.

package lower

func Find(nn []int, x int) (lb int, ok bool) {
	right := len(nn)
	for left := 0; left < right; {
		mid := left + (right-left)/2
		switch n := nn[mid]; {
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
