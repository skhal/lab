// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engine

import (
	"fmt"

	"github.com/skhal/lab/x/sheet/internal/parse"
	"github.com/skhal/lab/x/sheet/internal/vm"
)

// VirtualMachine uses bytecode (instructions set) for intermediate
// representation in parsed cells.
type VirtualMachine struct{}

// NewVirtualMachine creates a new VM engine.
func NewVirtualMachine() *VirtualMachine {
	return new(VirtualMachine)
}

// Parse parses cell content and returns bytecode.
func (eng *VirtualMachine) Parse(s string) (any, error) {
	ast, err := parse.Parse(s)
	if err != nil {
		return nil, err
	}
	return vm.Compile(ast)
}

// Calculate run bytecode through Virtual Machine. It reports an error if IR is
// not bytecode.
func (eng *VirtualMachine) Calculate(data any, refcal func(string) (float64, error)) (float64, error) {
	switch ir := data.(type) {
	case *vm.InstructionsSet:
		return vm.Run(ir, refcal)
	default:
		return 0, fmt.Errorf("unsupported IR - %T", ir)
	}
}
