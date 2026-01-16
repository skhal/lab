// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package palindrome

import (
	"fmt"
	"iter"
)

type Node struct {
	val  int
	next *Node
}

func (node *Node) String() string {
	if node == nil {
		return "<nil>"
	}
	var nn []int
	for node != nil {
		nn = append(nn, node.val)
		node = node.next
	}
	return fmt.Sprint(nn)
}

func NewList(nn ...int) *Node {
	var (
		head *Node
		last *Node
	)
	for _, n := range nn {
		node := &Node{
			val: n,
		}
		switch {
		case head == nil:
			head = node
		default:
			last.next = node
		}
		last = node
	}
	return head
}

func Is(n *Node) bool {
	for leftNode, rightNode := range split(n) {
		if leftNode.val != rightNode.val {
			return false
		}
	}
	return true
}

func split(n *Node) iter.Seq2[*Node, *Node] {
	return func(yield func(*Node, *Node) bool) {
		if n == nil {
			return
		}
		if n.next == nil {
			return
		}
		backward, forward := splitInHalf(n)
		for backward != nil && forward != nil {
			if !yield(backward.node, forward) {
				break
			}
			backward = backward.next
			forward = forward.next
		}
	}
}

type backwardNode struct {
	node *Node
	next *backwardNode
}

func splitInHalf(n *Node) (*backwardNode, *Node) {
	var (
		backward *backwardNode
		forward  = n
	)
	for step := 1; n != nil; n, step = n.next, step+1 {
		switch step % 2 {
		case 0:
			forward = forward.next
		case 1:
			backward = &backwardNode{
				node: forward,
				next: backward,
			}
		}
	}
	return backward, forward
}
