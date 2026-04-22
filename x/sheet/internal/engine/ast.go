// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package engine holds different kinds of engines to drive the sheets.
//
// The engines differ in how they store parsed cell data, i.e. intermediate
// representation (IR), and how to calculate the value from the IR.
//
// For example, an AST engine may use an AST for IR. A VM engine may use an
// instruction set in the form of bytecode and calculate the result using
// virtual machines.
//
// One of the noticeable difference between the IRs is that some can be saved
// while others need to be constructed from scratch.
package engine

import (
	"fmt"

	"github.com/skhal/lab/x/sheet/internal/ast"
	"github.com/skhal/lab/x/sheet/internal/calc"
	"github.com/skhal/lab/x/sheet/internal/parse"
)

// AST engine uses AST for intermediate representation.
type AST struct{}

// NewAST creates an AST engine.
func NewAST() *AST {
	return new(AST)
}

// Parse parses a cell value into an AST node.
func (eng *AST) Parse(s string) (any, error) {
	return parse.Parse(s)
}

// Calculate evaluates cell's AST node.
func (eng *AST) Calculate(data any, refcal func(string) (float64, error)) (float64, error) {
	switch n := data.(type) {
	case ast.Node:
		return calc.Calculate(n, refcal)
	default:
		return 0, fmt.Errorf("unsupported data type - %T", n)
	}
}
