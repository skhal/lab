// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package calc evaluates formulas.
package calc

import (
	"strconv"

	"github.com/skhal/lab/x/sheet/internal/ast"
)

// Calculate evaluates a formula node and skips other types of nodes. It
// returns an error if evaluation fails.
func Calculate(node ast.Node) (float64, error) {
	switch n := node.(type) {
	case *ast.NumberNode:
		return calculateNumber(n)
	}
	return 0, nil
}

func calculateNumber(n *ast.NumberNode) (float64, error) {
	return strconv.ParseFloat(n.Number, 64)
}
