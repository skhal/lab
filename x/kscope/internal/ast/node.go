// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

import (
	"fmt"
	"strings"
)

// Node is any AST node.
type Node any

// Func is a function definition. It consists of a name and a body, represented
// by a set of statements.
type Func struct {
	Name   string   // function name
	Params []string // parameter names
	Body   []Node   // function body
}

// String prints the function.
func (f Func) String() string {
	var s strings.Builder
	fmt.Fprintf(&s, "def %s(", f.Name)
	for i, p := range f.Params {
		if i > 0 {
			fmt.Fprintf(&s, ", ")
		}
		fmt.Fprintf(&s, "%s", p)
	}
	fmt.Fprintln(&s, ")")
	for _, n := range f.Body {
		fmt.Fprintf(&s, "  %s\n", n)
	}
	return s.String()
}

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

// Call is a function call.
type Call struct {
	Name string // function name
	Args []Node // arguments
}

// String prints a function call.
func (c Call) String() string {
	var s strings.Builder
	fmt.Fprintf(&s, "%s(", c.Name)
	for i, arg := range c.Args {
		if i > 0 {
			fmt.Fprint(&s, ", ")
		}
		fmt.Fprint(&s, arg)
	}
	fmt.Fprint(&s, ")")
	return s.String()
}

// Var is a variable definition, e.g. "var name = 1 + 3".
type Var struct {
	Val  Node   // value expression
	Name string // variable name
}

// String prints the variable.
func (v Var) String() string {
	return fmt.Sprintf("var %s = %s", v.Name, v.Val)
}
