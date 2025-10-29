// Copyright 2025 Samvel Khalatyan. All rights reserved.

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
