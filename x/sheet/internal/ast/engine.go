// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

import (
	"fmt"
)

// Engine uses AST for intermediate representation.
type Engine struct{}

// Parse parses a cell value into an AST node.
func (Engine) Parse(s string) (any, error) {
	return Parse(s)
}

// Calculate evaluates cell's AST node.
func (Engine) Calculate(data any, refcal func(string) (float64, error)) (float64, error) {
	switch ir := data.(type) {
	case Node:
		return Calculate(ir, refcal)
	default:
		return 0, fmt.Errorf("unsupported IR - %T", ir)
	}
}

// WriteIR disables intermediate representation in the Sheet write/read state.
func (Engine) WriteIR() bool { return false }
