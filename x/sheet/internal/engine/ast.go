// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engine

import (
	"fmt"

	"github.com/skhal/lab/x/sheet/internal/ast"
)

// AST engine uses AST for intermediate representation.
type AST struct{}

// Parse parses a cell value into an AST node.
func (AST) Parse(s string) (any, error) {
	return ast.Parse(s)
}

// Calculate evaluates cell's AST node.
func (AST) Calculate(data any, refcal func(string) (float64, error)) (float64, error) {
	switch ir := data.(type) {
	case ast.Node:
		return ast.Calculate(ir, refcal)
	default:
		return 0, fmt.Errorf("unsupported IR - %T", ir)
	}
}
