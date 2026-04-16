// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package calc evaluates formulas.
package calc

import "github.com/skhal/lab/x/sheet/internal/ast"

// Calculate evaluates a formula node and skips other types of nodes. It
// returns an error if evaluation fails.
func Calculate(node ast.Node) error {
	switch n := node.(type) {
	case *ast.FormulaNode:
		return calcFormula(n)
	default:
		return nil
	}
}

func calcFormula(n *ast.FormulaNode) error {
	n.Result = n.Number.Value()
	return nil
}
