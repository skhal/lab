// Copyright 2025 Samvel Khalatyan. All rights reserved.

package twosum

func Find(nn []int, x int) []int {
	seen := make(map[int]int)
	for i, n := range nn {
		target := x - n
		if j, ok := seen[target]; ok {
			return []int{i, j}
		}
		seen[n] = i
	}
	return nil
}
