// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

import "fmt"

// Node is any AST node.
type Node any

// Number is a number literal node.
type Number struct {
	Val float64 // parsed value of the number.
}

// String prints the number.
func (n Number) String() string { return fmt.Sprintf("%.1f", n.Val) }

// BinExpr describes a binary expression: left op right.
type BinExpr struct {
	Left  Node  // left operand
	Right Node  // right operand
	Op    BinOp // operator
}

// BinOp enumerates binary operators
//
//go:generate stringer -type=BinOp -linecomment
type BinOp int8

const (
	_ BinOp = iota

	// keep-sorted start
	BinOpDiv   // /
	BinOpMinus // -
	BinOpMul   // *
	BinOpPlus  // +
	// keep-sorted end
)

// String prints binary expression.
func (expr BinExpr) String() string {
	return fmt.Sprintf("%s %v %s", expr.Left, expr.Op, expr.Right)
}
