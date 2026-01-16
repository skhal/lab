// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package midpoint

type Node struct {
	Val  int
	Next *Node
}

func Find(node *Node) *Node {
	it := node
	for node != nil && node.Next != nil {
		node = node.Next.Next
		it = it.Next
	}
	return it
}
