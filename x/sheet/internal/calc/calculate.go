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

// Calculate evaluates a formula node and skips other types of nodes. It
// returns an error if evaluation fails.
func Calculate(node ast.Node) (float64, error) {
	switch n := node.(type) {
	case *ast.NumberNode:
		return calculateNumber(n)
	case *ast.BinOpNode:
		return calculateBinaryOperator(n)
	}
	return 0, ErrCalculate
}

func calculateNumber(n *ast.NumberNode) (float64, error) {
	return strconv.ParseFloat(n.Number, 64)
}

func calculateBinaryOperator(n *ast.BinOpNode) (_ float64, err error) {
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
		op := newBinaryOperator(plus, n.Left, n.Right)
		return op.Calculate(), nil
	case opMinus:
		op := newBinaryOperator(minus, n.Left, n.Right)
		return op.Calculate(), nil
	}

	return 0, fmt.Errorf("%w: unsupported operator %q", ErrCalculate, n.Op)
}

type binaryOperator struct {
	f        func(x, y float64) float64
	lhs, rhs ast.Node
}

func newBinaryOperator(op func(x, y float64) float64, lhs, rhs ast.Node) *binaryOperator {
	return &binaryOperator{op, lhs, rhs}
}

// Calculate executes binary operator after calculating left and right operands.
func (op *binaryOperator) Calculate() float64 {
	return op.f(mustCalculate(op.lhs), mustCalculate(op.rhs))
}

func plus(x, y float64) float64 {
	return x + y
}

func minus(x, y float64) float64 {
	return x - y
}

func mustCalculate(n ast.Node) float64 {
	res, err := Calculate(n)
	if err != nil {
		panic(err)
	}
	return res
}
