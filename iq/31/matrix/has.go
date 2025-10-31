// Copyright 2025 Samvel Khalatyan. All rights reserved.

package matrix

type M map[int][]int

func Has(m M, n int) bool {
	if len(m) == 0{
		return false
	}
	left := 0
	right := len(m) * len(m[0])
	for left < right {
		mid := left + (right - left) / 2
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
