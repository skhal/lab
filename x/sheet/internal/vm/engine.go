// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vm

import (
	"fmt"

	"github.com/skhal/lab/x/sheet/internal/ast"
)

// Engine uses bytecode (instructions set) for intermediate
// representation in parsed cells.
type Engine struct{}

// Parse parses cell content and returns bytecode.
func (Engine) Parse(s string) (any, error) {
	ast, err := ast.Parse(s)
	if err != nil {
		return nil, err
	}
	return Compile(ast)
}

// Calculate run bytecode through Virtual Machine. It reports an error if IR is
// not bytecode.
func (Engine) Calculate(data any, refcal func(string) (float64, error)) (float64, error) {
	switch v := data.(type) {
	case InstructionsSet:
		return Run(&v, refcal)
	default:
		return 0, fmt.Errorf("unsupported IR - %T", v)
	}
}
