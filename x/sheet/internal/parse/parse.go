// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package parse implements lexical analysis of cell contents.
package parse

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/skhal/lab/x/sheet/internal/ast"
)

// ErrParse means the was en error parsing data.
var ErrParse = errors.New("parse error")

// Parse runs lexical and syntacit analysis of s. It returns an AST node upon
// success or a non-nil error in case of failure.
func Parse(s string) (ast.Node, error) {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrParse, s)
	}
	return &ast.NumberNode{Number: n}, nil
}
