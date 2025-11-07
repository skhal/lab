// Copyright 2025 Samvel Khalatyan. All rights reserved.

package intersection

import (
	"fmt"
	"iter"
)

type Node struct {
	val  int
	next *Node
}

func NewList(nn ...int) *Node {
	var (
		head *Node
		tail *Node
	)
	for _, n := range nn {
		node := &Node{val: n}
		switch {
		case head == nil:
			head = node
		case tail != nil:
			tail.next = node
		}
		tail = node
	}
	return head
}

func (node *Node) String() string {
	nn := node.ToList()
	return fmt.Sprint(nn)
}

func (node *Node) ToList() []int {
	var nn []int
	for node != nil {
		nn = append(nn, node.val)
		node = node.next
	}
	return nn
}

func (node *Node) Append(tail *Node) *Node {
	head := node
	for node.next != nil {
		node = node.next
	}
	node.next = tail
	return head
}

func Find(a, b *Node) *Node {
	if a == nil || b == nil {
		return nil
	}
	for na, nb := range zip(join(a, b), join(b, a)) {
		if na == nb {
			return na
		}
	}
	return nil
}

func join(a, b *Node) iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		walk := func(n *Node) bool {
			for n != nil {
				if !yield(n) {
					return false
				}
				n = n.next
			}
			return true
		}
		if !walk(a) {
			return
		}
		if !walk(b) {
			return
		}
	}
}

func zip(a, b iter.Seq[*Node]) iter.Seq2[*Node, *Node] {
	return func(yield func(x, y *Node) bool) {
		nexta, stopa := iter.Pull(a)
		defer stopa()
		nextb, stopb := iter.Pull(b)
		defer stopb()
		for {
			x, xok := nexta()
			y, yok := nextb()
			if xok == false && yok == false {
				break
			}
			if !yield(x, y) {
				break
			}
		}
	}
}
