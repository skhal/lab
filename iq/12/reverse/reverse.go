// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reverse

import (
	"bytes"
	"strconv"
)

type Node struct {
	Value int
	Next  *Node
}

func NewNode(val int) *Node {
	return &Node{
		Value: val,
	}
}

type List struct {
	Head *Node
}

func NewList(nn ...int) *List {
	list := new(List)
	var last *Node
	for _, n := range nn {
		node := NewNode(n)
		if last == nil {
			list.Head = node
		} else {
			last.Next = node
		}
		last = node
	}
	return list
}

func (l List) String() string {
	const (
		ByteOpen   = '['
		ByteClose  = ']'
		ByteString = ' '
	)
	buf := new(bytes.Buffer)
	buf.WriteByte(ByteOpen)
	for node := l.Head; node != nil; node = node.Next {
		buf.WriteString(strconv.Itoa(node.Value))
		if node.Next != nil {
			buf.WriteByte(ByteString)
		}
	}
	buf.WriteByte(ByteClose)
	return buf.String()
}

func (l *List) Reverse() {
	node := l.Head
	if node == nil {
		// empty
		return
	}
	if node.Next == nil {
		// one item
		return
	}
	var prev *Node
	for node != nil {
		next := node.Next
		node.Next = prev
		prev, node = node, next
	}
	l.Head = prev
}
