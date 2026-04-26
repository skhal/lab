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

// RefNode is a cell reference.
type RefNode struct {
	Ref string // cell reference
}

// CallNode is a function call.
type CallNode struct {
	Name string // function name
	Args []Node // function arguments, can be empty
}

// RangeNode represents a cells rang, e.g. A1:A5
type RangeNode struct {
	From string // from cell
	To   string // to cell (inclusive)
}

// IfNode is if-clause, e.g. IF(Cond, IfPassExpr, IfFailExpr).
type IfNode struct {
	Cond   *RelOpNode // a comparison binary expression
	IfPass Node       // an expression to be executed when if-condition is true
	IfFail Node       // an expression to be executed when if-condition is false
}

// RelOpNode is a comparison binary operation, e.g. "A1 < 5".
type RelOpNode struct {
	Op    string // comparison operator: ==, !=, <, <=, >, >=
	Left  Node   // left operand
	Right Node   // right operand
}
