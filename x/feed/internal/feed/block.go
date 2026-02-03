// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package feed

import "iter"

// Block is a range of indices [Low, High).
type Block struct {
	Low  int // start index (inclusive)
	High int // end index (exclusive)
}

// EqualSizeBlocks generates a sequence of blocks of a given size out of n
// items.
func EqualSizeBlocks(size, n int) iter.Seq[Block] {
	return func(yield func(Block) bool) {
		for i, h := 0, 0; i < n; i = h {
			h = min(n, i+size)
			if !yield(Block{i, h}) {
				break
			}
		}
	}
}
