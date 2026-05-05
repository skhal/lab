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

// Code is the top-level container to represent a code segment.
type Code struct {
	Decls []*Decl // top-level declarations
}

// String prints the code segment.
func (f Code) String() string {
	var s strings.Builder
	for i, d := range f.Decls {
		if i > 0 {
			fmt.Fprintln(&s)
		}
		fmt.Fprint(&s, d)
	}
	return s.String()
}

// Decl is a top-level declaration in the file.
type Decl struct {
	Node // variable or a function
}

// String prints the declaration.
func (d Decl) String() string {
	return fmt.Sprintf("%s", d.Node)
}

// Func is a function definition.
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
	for i, n := range f.Body {
		if i > 0 {
			fmt.Fprintln(&s)
		}
		fmt.Fprintf(&s, "  %s", n)
	}
	return s.String()
}

// Number is a number literal node.
type Number struct {
	Val float64 // parsed value of the number.
}

// String prints the number.
func (n Number) String() string { return fmt.Sprintf("%.1f", n.Val) }

// Ident is an identifier used in an expression, e.g. a variable in "x + 3"
type Ident struct {
	Name string // identifier name
}

// String prints the identifier.
func (i Ident) String() string { return i.Name }

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
