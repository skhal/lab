// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shiftzeros

import "iter"

func Shift(nn []int) {
	offset := find(nn, isZero)
	zi := offset
	for ni := range findAll(nn[offset:], isNonZero) {
		swap(nn, zi, offset+ni)
		zi += 1
	}
}

func find(nn []int, f func(n int) bool) int {
	for i, n := range nn {
		if f(n) {
			return i
		}
	}
	return len(nn)
}

func findAll(nn []int, f func(n int) bool) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i, n := range nn {
			if !f(n) {
				continue
			}
			if !yield(i) {
				return
			}
		}
	}
}

func isZero(n int) bool {
	return n == 0
}

func isNonZero(n int) bool {
	return n != 0
}

func swap(nn []int, i, j int) {
	if i == j {
		return
	}
	nn[i], nn[j] = nn[j], nn[i]
}
