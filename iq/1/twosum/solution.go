// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package twosum

// Find returns indices of the first pair of items from nn that add up to x.
func Find(nn []int, x int) []int {
	if len(nn) < 2 {
		return nil
	}
	left, right := 0, len(nn)-1
	for left < right {
		sum := nn[left] + nn[right]
		if sum < x {
			left += 1
		} else if sum > x {
			right -= 1
		} else {
			return []int{left, right}
		}
	}
	return nil
}
