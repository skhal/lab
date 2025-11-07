// Copyright 2025 Samvel Khalatyan. All rights reserved.

package flatten

import "fmt"

func NewTree(nn ...int) *TreeNode {
	var (
		head *TreeNode
		curr *TreeNode
	)
	for _, n := range nn {
		node := &TreeNode{
			val: n,
		}
		if head == nil {
			head = node
			curr = node
			continue
		}
		curr.next = node
		curr = node
	}
	return head
}

type TreeNode struct {
	val   int
	next  *TreeNode
	child *TreeNode
}

func (tn *TreeNode) Get(val int) *TreeNode {
	for n := tn; n != nil; n = n.next {
		if n.val == val {
			return n
		}
	}
	return nil
}

func (tn *TreeNode) SetChild(child *TreeNode) {
	tn.child = child
}

type Node struct {
	val  int
	next *Node
}

func NewList(nn ...int) *Node {
	var (
		head *Node
		curr *Node
	)
	for _, n := range nn {
		node := &Node{
			val: n,
		}
		if head == nil {
			head = node
			curr = node
			continue
		}
		curr.next = node
		curr = node
	}
	return head
}

func (n *Node) Slice() []int {
	var vv []int
	for n != nil {
		vv = append(vv, n.val)
		n = n.next
	}
	return vv
}

func (n *Node) String() string {
	var vv []int
	for ; n != nil; n = n.next {
		vv = append(vv, n.val)
	}
	return fmt.Sprint(vv)
}

func Do(tn *TreeNode) *Node {
	var (
		head *Node
		tail *Node
	)
	for queue := []*TreeNode{tn}; len(queue) > 0; queue = queue[1:] {
		h, t, q := flattenOneLayer(queue[0])
		switch {
		case head == nil:
			head = h
		default:
			tail.next = h
		}
		tail = t
		queue = append(queue, q...)
	}
	return head
}

func flattenOneLayer(tn *TreeNode) (head *Node, tail *Node, queue []*TreeNode) {
	for ; tn != nil; tn = tn.next {
		n := &Node{val: tn.val}
		switch {
		case head == nil:
			head = n
		default:
			tail.next = n
		}
		tail = n
		if tn.child != nil {
			queue = append(queue, tn.child)
		}
	}
	return head, tail, queue
}
