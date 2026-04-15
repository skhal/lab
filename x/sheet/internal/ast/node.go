// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ast defines Abstract Syntax Tree (AST) nodes for a sheet cell.
package ast

// Node represents any node.
type Node interface {
	Value() float64 // retrieve node value
}

// NumberNode is a number value.
type NumberNode struct {
	Number float64 // node value
}

// Value returns the value of the NumberNode.
func (n *NumberNode) Value() float64 {
	return n.Number
}
