// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vm

import (
	"errors"
	"fmt"
)

// ErrRun means error running the virtual machine on the instructions set.
var ErrRun = errors.New("run error")

// Run executes the instructions set using a virtual machine (VM).
// Is uses refcal to calculate references.
func Run(bc *InstructionsSet, refcal func(string) (float64, error)) (float64, error) {
	r := runner{refcal: refcal}
	n, err := r.Run(bc)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrRun, err)
	}
	return n, nil
}

// runner implements a virtual machine to run the instructions set.
type runner struct {
	refcal func(string) (float64, error)
	stack  []float64
}

// Run executes the instructions set.
func (r *runner) Run(bc *InstructionsSet) (_ float64, err error) {
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
	r.stack = make([]float64, 0, len(bc.Instructions))
	if err := r.run(bc); err != nil {
		return 0, err
	}
	if len(r.stack) != 1 {
		return 0, fmt.Errorf("invalid instructions set - stack size %d", len(r.stack))
	}
	return r.pop(), nil
}

func (r *runner) run(bc *InstructionsSet) error {
	for _, instruction := range bc.Instructions {
		switch v := instruction.(type) {
		case Number:
			r.push(float64(v))
		case BinOp:
			n, err := r.runBinOp(v)
			if err != nil {
				return err
			}
			r.push(n)
		case Ref:
			n, err := r.refcal(string(v))
			if err != nil {
				return err
			}
			r.push(n)
		case Call:
			n, err := r.runCall(&v)
			if err != nil {
				return err
			}
			r.push(n)
		}
	}
	return nil
}

func (r *runner) runBinOp(op BinOp) (float64, error) {
	y := r.pop()
	x := r.pop()
	switch op {
	case BinOpPlus:
		return x + y, nil
	case BinOpMinus:
		return x - y, nil
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
