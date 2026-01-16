// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cycle

import "iter"

type Node struct {
	Val  int
	Next *Node
}

// Is reports whether the singly linked list has cycles. It uses fast and slow
// pointers algorithm instead of a set of visited nodes.
func Is(node *Node) bool {
	for it := range run(node) {
		if it.slow == it.fast {
			return true
		}
	}
	return false
}

type pointers struct {
	fast, slow *Node
}

func run(node *Node) iter.Seq[pointers] {
	return func(yield func(pointers) bool) {
		slow, fast := node, node
		// There is no need to check the slow for nil since the fast already checked
		// the condition.
		for fast != nil && fast.Next != nil {
			slow = slow.Next
			fast = fast.Next.Next
			if !yield(pointers{slow, fast}) {
				break
			}
		}
	}
}
