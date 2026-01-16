// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package remove

import "fmt"

type Node struct {
	Value int
	next  *Node
}

type List struct {
	head *Node
	size int
}

func NewList(nn ...int) *List {
	var (
		head *Node
		prev *Node
	)
	for _, n := range nn {
		node := &Node{
			Value: n,
		}
		if head == nil {
			head = node
			prev = node
			continue
		}
		prev.next = node
		prev = node
	}
	return &List{
		head: head,
		size: len(nn),
	}
}

func (l *List) String() string {
	nn := l.Slice()
	return fmt.Sprint(nn)
}

func (l *List) Slice() []int {
	if l == nil {
		return nil
	}
	nn := make([]int, 0, l.size)
	for node := l.head; node != nil; node = node.next {
		nn = append(nn, node.Value)
	}
	return nn
}

func LastNth(l *List, n int) {
	if l == nil {
		return
	}
	prev, ok := findLastNth(l.head, n)
	if !ok {
		return
	}
	switch {
	case prev == nil:
		removeHead(l)
	default:
		removeNode(prev)
	}
}

// findLastNth searches for the n-th node from the end starting from node. It
// returns a reference to the previous to n-th last node and flag whether the
// n-th last node was found. It returns (nil, true) if the length of the chain
// is equal to n starting from node.
func findLastNth(start *Node, n int) (prev *Node, ok bool) {
	size := 0
	for node := start; node != nil; node = node.next {
		size++
		switch {
		case size < n:
		case size == n:
			ok = true
		case prev == nil:
			prev = start
		default:
			prev = prev.next
		}
	}
	return
}

func removeHead(l *List) {
	next := l.head.next
	l.head.next = nil
	l.head = next
}

func removeNode(prev *Node) {
	node := prev.next
	prev.next = node.next
	node.next = nil
}
