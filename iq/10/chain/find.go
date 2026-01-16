// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chain

import "iter"

func Find(nn []int) []int {
	set := make(map[int]struct{})
	for _, n := range nn {
		set[n] = struct{}{}
	}
	var longestChain []int
	for chain := range findChains(set) {
		if len(chain) <= len(longestChain) {
			continue
		}
		longestChain = chain[:]
	}
	return longestChain
}

type C []int

func findChains(nn map[int]struct{}) iter.Seq[C] {
	return func(yield func(C) bool) {
		for n := range nn {
			if _, ok := nn[n-1]; ok {
				continue
			}
			chain := newChain(nn, n)
			if !yield(chain) {
				break
			}
		}
	}
}

func newChain(nn map[int]struct{}, start int) C {
	cc := []int{start}
	exists := func(n int) bool {
		_, ok := nn[n]
		return ok
	}
	for n := start + 1; exists(n); n++ {
		cc = append(cc, n)
	}
	return cc
}
