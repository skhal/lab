// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

import (
	"errors"
	"fmt"
	"strconv"
)

// ErrCalculate means the AST has an error can can't be calculated.
var ErrCalculate = errors.New("calculate error")

// Calculate evaluates a formula node and skips other types of nodes. It
// returns an error if evaluation fails.
func Calculate(n Node, refcal func(string) (float64, error)) (float64, error) {
	c := &calculator{refcal}
	return c.Calculate(n)
}

type calculator struct {
	refcal func(string) (float64, error)
}

// Calculate calculates the value of the node. It uses reference calculator to
// get the value of a reference.
func (c *calculator) Calculate(node Node) (float64, error) {
	switch n := node.(type) {
	// keep-sorted start
	case *BinOpNode:
		return c.calcBinOp(n)
	case *CallNode:
		return c.calcCall(n)
	case *IfNode:
		return c.calcIf(n)
	case *NumberNode:
		return c.calcNum(n)
	case *RefNode:
		return c.refcal(n.Ref)
		// keep-sorted end
	}
	return 0, ErrCalculate
}

func (c *calculator) calcNum(n *NumberNode) (float64, error) {
	return strconv.ParseFloat(n.Number, 64)
}

func (c *calculator) calcBinOp(n *BinOpNode) (_ float64, err error) {
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
		// keep-sorted start
		opDivide   = "/"
		opMinus    = "-"
		opMultiply = "*"
		opPlus     = "+"
		// keep-sorted end
	)
	switch n.Op {
	// keep-sorted start
	case opDivide:
		op := newBinaryOperator(c, divide, n.Left, n.Right)
		return op.Calculate(), nil
	case opMinus:
		op := newBinaryOperator(c, minus, n.Left, n.Right)
		return op.Calculate(), nil
	case opMultiply:
		op := newBinaryOperator(c, multiply, n.Left, n.Right)
		return op.Calculate(), nil
	case opPlus:
		op := newBinaryOperator(c, plus, n.Left, n.Right)
		return op.Calculate(), nil
		// keep-sorted end
	}
	return 0, fmt.Errorf("%w: unsupported operator %q", ErrCalculate, n.Op)
}

type binaryOperator struct {
	c        *calculator
	f        func(x, y float64) float64
	lhs, rhs Node
}

func newBinaryOperator(c *calculator, op func(x, y float64) float64, lhs, rhs Node) *binaryOperator {
	return &binaryOperator{c, op, lhs, rhs}
}

// Calculate executes binary operator after calculating left and right operands.
func (op *binaryOperator) Calculate() float64 {
	return op.f(must(op.c.Calculate(op.lhs)), must(op.c.Calculate(op.rhs)))
}

func divide(x, y float64) float64 {
	return x / y
}

func minus(x, y float64) float64 {
	return x - y
}

func multiply(x, y float64) float64 {
	return x * y
}

func plus(x, y float64) float64 {
	return x + y
}

func must(n float64, err error) float64 {
	if err != nil {
		panic(err)
	}
	return n
}

var calls = map[string]func([]float64) float64{
	// keep-sorted start
	"MIN": nil,
	"SUM": callSum,
	// keep-sorted end
}

func (c *calculator) calcCall(n *CallNode) (float64, error) {
	fn, ok := calls[n.Name]
	if !ok {
		return 0, fmt.Errorf("unsupported formula - %s", n.Name)
	}
	if fn == nil {
		return 0, fmt.Errorf("disabled formula - %s", n.Name)
	}
	args := make([]float64, 0, len(n.Args))
	for _, na := range n.Args {
		switch n := na.(type) {
		case *RangeNode:
			rng, err := c.calcRange(n)
			if err != nil {
				return 0, err
			}
			args = append(args, rng...)
			continue
		}
		arg, err := c.Calculate(na)
		if err != nil {
			return 0, err
		}
		args = append(args, arg)
	}
	return fn(args), nil
}

func (c *calculator) calcRange(n *RangeNode) ([]float64, error) {
	cr, err := NewCellScanner(n.From, n.To)
	if err != nil {
		return nil, err
	}
	nn := make([]float64, 0, cr.Len())
	for id := range cr.Scan() {
		res, err := c.Calculate(&RefNode{Ref: id})
		if err != nil {
			return nil, err
		}
		nn = append(nn, res)
	}
	return nn, nil
}

func callSum(nn []float64) float64 {
	var sum float64
	for _, n := range nn {
		sum += n
	}
	return sum
}

func (c *calculator) calcIf(n *IfNode) (float64, error) {
	cond, err := c.calcRelOp(n.Cond)
	if err != nil {
		return 0, err
	}
	if !cond {
		return c.Calculate(n.IfFail)
	}
	return c.Calculate(n.IfPass)
}

func (c *calculator) calcRelOp(n *RelOpNode) (bool, error) {
	op, ok := relOp[n.Op]
	if !ok {
		return false, fmt.Errorf("unsupported comparison operator %s", n.Op)
	}
	x := must(c.Calculate(n.Left))
	y := must(c.Calculate(n.Right))
	return op(x, y), nil
}

var relOp = map[string]func(x, y float64) bool{
	"==": Equal,
	"!=": func(x, y float64) bool { return !Equal(x, y) },
	"<":  func(x, y float64) bool { return x < y },
	"<=": func(x, y float64) bool { return x <= y },
	">":  func(x, y float64) bool { return x > y },
	">=": func(x, y float64) bool { return x >= y },
}
