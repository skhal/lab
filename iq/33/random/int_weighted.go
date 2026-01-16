// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package random

func IntWeighted(weights []int, rand func(n int) int) int {
	ww, wmax := cumulativeDistributionFunction(weights)
	x := rand(wmax)
	return pick(ww, x)
}

func pick(ww []int, n int) int {
	left := 0
	right := len(ww)
	for left < right {
		mid := left + (right-left)/2
		w := ww[mid]
		switch {
		case n > w:
			left = mid + 1
		case n == w:
			return mid + 1
		case mid > 0 && n < ww[mid-1]:
			right = mid
		default:
			return mid
		}
	}
	return left
}

func cumulativeDistributionFunction(ww []int) ([]int, int) {
	sum := 0
	ss := make([]int, 0, len(ww))
	for _, w := range ww {
		sum += w
		ss = append(ss, sum)
	}
	return ss, sum
}
