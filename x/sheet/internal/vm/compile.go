// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vm

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/skhal/lab/x/sheet/internal/ast"
)

// ErrCompile means error to compile AST into the instructions set.
var ErrCompile = errors.New("compile error")

// Compile converts AST to the instructions set.
func Compile(n ast.Node) (InstructionsSet, error) {
	if n == nil {
		return InstructionsSet{}, nil
	}
	c := new(compiler)
	iset, err := c.Compile(n)
	if err != nil {
		return InstructionsSet{}, fmt.Errorf("%w: %s", ErrCompile, err)
	}
	return iset, nil
}

// compiler generates the instructions set for IR.
type compiler struct {
	iset []Inst // instructions set
}

// Compile converts an Abstract Syntax Tree (AST) to the instructions set.
func (c *compiler) Compile(node ast.Node) (InstructionsSet, error) {
	if err := c.compile(node); err != nil {
		return InstructionsSet{}, err
	}
	return InstructionsSet{c.iset}, nil
}

func (c *compiler) compile(node ast.Node) error {
	switch n := node.(type) {
	// keep-sorted start
	case *ast.BinOpNode:
		return c.compileBinOp(n)
	case *ast.CallNode:
		return c.compileCall(n)
	case *ast.NumberNode:
		return c.compileNumber(n)
	case *ast.RefNode:
		return c.compileRef(n)
		// keep-sorted end
	default:
		return fmt.Errorf("unsupported node %T", n)
	}
}

func (c *compiler) compileNumber(num *ast.NumberNode) error {
	n, err := strconv.ParseFloat(num.Number, 64)
	if err != nil {
		return err
	}
	c.push(Inst{Type: InstTypeNumber, Number: n})
	return nil
}

func (c *compiler) compileBinOp(op *ast.BinOpNode) error {
	if err := c.compile(op.Left); err != nil {
		return err
	}
	if err := c.compile(op.Right); err != nil {
		return err
	}
	const (
		plus  = "+"
		minus = "-"
	)
	switch op.Op {
	case plus:
		c.push(Inst{Type: InstTypeBinOp, BinOp: BinOpPlus})
		return nil
	case minus:
		c.push(Inst{Type: InstTypeBinOp, BinOp: BinOpMinus})
		return nil
	}
	return fmt.Errorf("unsupported binary operator %s", op.Op)
}

func (c *compiler) compileCall(call *ast.CallNode) error {
	var args int
	for _, arg := range call.Args {
		switch n := arg.(type) {
		case *ast.RangeNode:
			num, err := c.compileRange(n)
			if err != nil {
				return err
			}
			args += num
			continue
		}
		if err := c.compile(arg); err != nil {
			return err
		}
		args++
	}
	fn, ok := calls[call.Name]
	if !ok {
		return fmt.Errorf("unsupported call %s", call.Name)
	}
	c.push(Inst{Type: InstTypeCall, Call: &Call{Func: fn, Args: args}})
	return nil
}

func (c *compiler) compileRange(n *ast.RangeNode) (int, error) {
	cr, err := ast.NewCellScanner(n.From, n.To)
	if err != nil {
		return 0, err
	}
	for id := range cr.Scan() {
		c.push(Inst{Type: InstTypeRef, Ref: id})
	}
	return cr.Len(), nil
}

func (c *compiler) compileRef(ref *ast.RefNode) error {
	c.push(Inst{Type: InstTypeRef, Ref: ref.Ref})
	return nil
}

func (c *compiler) push(v Inst) {
	c.iset = append(c.iset, v)
}
