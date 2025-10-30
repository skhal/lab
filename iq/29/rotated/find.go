// Copyright 2025 Samvel Khalatyan. All rights reserved.

package rotated

type Index int

const IndexError = -1

func Find(nn []int, x int) Index {
	left := 0
	right := len(nn)
	isRotated := func() bool {
		return nn[left] > nn[right-1]
	}
	canMoveRight := func(x int) bool {
		if !isRotated() {
			return false
		}
		return x <= nn[right-1]
	}
	moveRight := func(mid int) {
		left = mid + 1
	}
	canMoveLeft := func(x int) bool {
		if !isRotated() {
			return false
		}
		return nn[left] <= x
	}
	moveLeft := func(mid int) {
		right = mid
	}
	for left < right {
		mid := left + (right-left)/2
		n := nn[mid]
		switch {
		case x < n:
			if canMoveRight(x) {
				moveRight(mid)
				break
			}
			moveLeft(mid)
		case x > n:
			if canMoveLeft(x) {
				moveLeft(mid)
				right = mid
				break
			}
			moveRight(mid)
		default:
			return Index(mid)
		}
	}
	return IndexError
}
