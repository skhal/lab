// Copyright 2025 Samvel Khalatyan. All rights reserved.

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
