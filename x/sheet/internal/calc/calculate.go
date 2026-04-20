// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package calc evaluates formulas.
package calc

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/skhal/lab/x/sheet/internal/ast"
)

// ErrCalculate means the AST has an error can can't be calculated.
var ErrCalculate = errors.New("calculate error")

// RefCalculator calculates reference value.
type RefCalculator interface {
	Calculate(string) (float64, error) // calculate reference cell value
}

// Calculate evaluates a formula node and skips other types of nodes. It
// returns an error if evaluation fails.
func Calculate(n ast.Node, ref RefCalculator) (float64, error) {
	c := &calculator{ref}
	return c.Calculate(n)
}

type calculator struct {
	ref RefCalculator
}

// Calculate calculates the value of the node. It uses RefCalculator to get
// the value of a reference.
func (c *calculator) Calculate(node ast.Node) (float64, error) {
	switch n := node.(type) {
	case *ast.NumberNode:
		return c.calcNum(n)
	case *ast.BinOpNode:
		return c.calcBinOp(n)
	case *ast.RefNode:
		return c.ref.Calculate(n.Ref)
	}
	return 0, ErrCalculate
}

func (c *calculator) calcNum(n *ast.NumberNode) (float64, error) {
	return strconv.ParseFloat(n.Number, 64)
}

func (c *calculator) calcBinOp(n *ast.BinOpNode) (_ float64, err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		e, ok := r.(error)
		if !ok {
			return
		}
		err = fmt.Errorf("%w: %s", ErrCalculate, e)
	}()
	const (
		opPlus  = "+"
		opMinus = "-"
	)
	switch n.Op {
	case opPlus:
		op := newBinaryOperator(c, plus, n.Left, n.Right)
		return op.Calculate(), nil
	case opMinus:
		op := newBinaryOperator(c, minus, n.Left, n.Right)
		return op.Calculate(), nil
	}

	return 0, fmt.Errorf("%w: unsupported operator %q", ErrCalculate, n.Op)
}

type binaryOperator struct {
	c        *calculator
	f        func(x, y float64) float64
	lhs, rhs ast.Node
}

func newBinaryOperator(c *calculator, op func(x, y float64) float64, lhs, rhs ast.Node) *binaryOperator {
	return &binaryOperator{c, op, lhs, rhs}
}

// Calculate executes binary operator after calculating left and right operands.
func (op *binaryOperator) Calculate() float64 {
	return op.f(must(op.c.Calculate(op.lhs)), must(op.c.Calculate(op.rhs)))
}

func plus(x, y float64) float64 {
	return x + y
}

func minus(x, y float64) float64 {
	return x - y
}

func must(n float64, err error) float64 {
	if err != nil {
		panic(err)
	}
	return n
}
