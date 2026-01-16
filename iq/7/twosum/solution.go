// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
