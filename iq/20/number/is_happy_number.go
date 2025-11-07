// Copyright 2025 Samvel Khalatyan. All rights reserved.

package number

import "iter"

const happyNumEnd = 1

// IsHappyNumber reports if the number n is "happy", i.e., iterative sum of
// squares of its digits leads to 1.
func IsHappyNumber(n int) bool {
	for slow, fast := range runNumbers(n) {
		if fast == happyNumEnd {
			break
		}
		if slow == fast {
			return false
		}
	}
	return true
}

func runNumbers(n int) iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		if n <= happyNumEnd {
			return
		}
		slow, fast := n, getNextNumber(n)
		for fast != happyNumEnd {
			if !yield(slow, fast) {
				break
			}
			slow = getNextNumber(slow)
			fast = getNextNumber(fast)
			fast = getNextNumber(fast)
		}
	}
}

func getNextNumber(n int) int {
	next := 0
	for ; n > 0; n /= 10 {
		i := n % 10
		if i == 0 {
			continue
		}
		next += i * i
	}
	return next
}
