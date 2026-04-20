// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ast defines Abstract Syntax Tree (AST) nodes for a sheet cell.
package ast

// Node represents any node.
type Node any

// NumberNode is a number value.
type NumberNode struct {
	Number string // number text value
}

// BinOpNode is a binary operation of the form "left op right".
type BinOpNode struct {
	Op    string // binary operator
	Left  Node   // left operand
	Right Node   // right operand
}
