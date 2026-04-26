// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vm

import (
	"errors"
	"fmt"

	"github.com/skhal/lab/x/sheet/internal/ast"
)

// ErrRun means error running the virtual machine on the instructions set.
var ErrRun = errors.New("run error")

// Run executes the instructions set using a virtual machine (VM).
// Is uses refcal to calculate references.
func Run(iset *InstructionsSet, refcal func(string) (float64, error)) (float64, error) {
	r := runner{refcal: refcal}
	n, err := r.Run(iset)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrRun, err)
	}
	return n, nil
}

// runner implements a virtual machine to run the instructions set.
type runner struct {
	refcal func(string) (float64, error)
	stack  []float64
	ip     int // instruction position
	iplen  int
}

// Run executes the instructions set.
func (r *runner) Run(iset *InstructionsSet) (_ float64, err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		e, ok := r.(error)
		if !ok {
			return
		}
		err = e
	}()
	r.stack = make([]float64, 0, len(iset.Instructions))
	if err := r.run(iset); err != nil {
		return 0, err
	}
	if len(r.stack) == 0 {
		return 0, fmt.Errorf("empty stack")
	}
	return r.pop(), nil
}

func (r *runner) run(iset *InstructionsSet) error {
	for r.ip, r.iplen = 0, len(iset.Instructions); r.ip < r.iplen; r.ip++ {
		inst := iset.Instructions[r.ip]
		switch inst.Type {
		case InstTypeNumber:
			r.push(inst.Number)
		case InstTypeBinOp:
			n, err := r.runBinOp(inst.BinOp)
			if err != nil {
				return err
			}
			r.push(n)
		case InstTypeRef:
			n, err := r.refcal(inst.Ref)
			if err != nil {
				return err
			}
			r.push(n)
		case InstTypeCall:
			n, err := r.runCall(inst.Call)
			if err != nil {
				return err
			}
			r.push(n)
		case InstTypeIfCall:
			if err := r.runIfCall(inst.IfCall); err != nil {
				return err
			}
		case InstTypeJump:
			if err := r.runJump(inst.JumpOffset); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *runner) runBinOp(op BinOp) (float64, error) {
	y := r.pop()
	x := r.pop()
	switch op {
	// keep-sorted start
	case BinOpDivide:
		return x / y, nil
	case BinOpMinus:
		return x - y, nil
	case BinOpMultiply:
		return x * y, nil
	case BinOpPlus:
		return x + y, nil
		// keep-sorted end
	}
	return 0, fmt.Errorf("unsupported operation %s", op)
}

func (r *runner) pop() float64 {
	if len(r.stack) == 0 {
		panic(fmt.Errorf("pop: empty stack"))
	}
	n := r.stack[len(r.stack)-1]
	r.stack = r.stack[:len(r.stack)-1]
	return n
}

func (r *runner) push(n float64) {
	r.stack = append(r.stack, n)
}

func (r *runner) runCall(c *Call) (float64, error) {
	args := make([]float64, c.Args)
	for i := c.Args; i > 0; i-- {
		args[i-1] = r.pop()
	}
	switch c.Func {
	case FuncSum:
		return sum(args...), nil
	}
	return 0, fmt.Errorf("unsupported function %s", c.Func)
}

func sum(nn ...float64) float64 {
	var res float64
	for _, n := range nn {
		res += n
	}
	return res
}

var relOps = map[RelOp]func(x, y float64) bool{
	RelOpEqual:          ast.Equal,
	RelOpNotEqual:       func(x, y float64) bool { return !ast.Equal(x, y) },
	RelOpLess:           func(x, y float64) bool { return x < y },
	RelOpLessOrEqual:    func(x, y float64) bool { return x <= y },
	RelOpGreater:        func(x, y float64) bool { return x > y },
	RelOpGreaterOrEqual: func(x, y float64) bool { return x >= y },
}

func (r *runner) runIfCall(ifc *IfCall) error {
	y := r.pop()
	x := r.pop()
	op, ok := relOps[ifc.RelOp]
	if !ok {
		return fmt.Errorf("unsupported comparison operator %s", ifc.RelOp)
	}
	if op(x, y) {
		return nil
	}
	return r.runJump(ifc.IfFail)
}

func (r *runner) runJump(offset JumpOffset) error {
	if offset == 0 {
		return fmt.Errorf("empty jump offset")
	}
	r.ip += int(offset)
	return nil
}
