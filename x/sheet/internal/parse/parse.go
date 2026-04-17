// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package parse implements lexical analysis of cell contents.
package parse

import (
	"errors"
	"fmt"
	"strings"

	"github.com/skhal/lab/x/sheet/internal/ast"
	"github.com/skhal/lab/x/sheet/internal/lex"
)

// ErrParse means the was en error parsing data.
var ErrParse = errors.New("parse error")

// Parse runs lexical and syntacit analysis of s. It returns an AST node upon
// success or a non-nil error in case of failure.
func Parse(s string) (ast.Node, error) {
	const (
		formulaPrefix = "="
	)
	switch s, ok := strings.CutPrefix(s, formulaPrefix); {
	case ok:
		return parseFormula(s)
	case len(s) != 0:
		return parseNode(s)
	}
	return nil, nil
}

// parseFormula parses a formula string without "=" prefix. It returns a
// formula AST node upon success or error.
func parseFormula(s string) (node ast.Node, _ error) {
	for tok := range lex.Lex([]byte(s)) {
		switch tok.Type {
		case lex.TokenError:
			err := fmt.Errorf("%w: formula %q: %s", ErrParse, s, tok.Err)
			return nil, err
		case lex.TokenNumber:
			if node != nil {
				err := fmt.Errorf("%w: formula %q: multiple numbers", ErrParse, s)
				return nil, err
			}
			node = &ast.NumberNode{Number: tok.Text}
		default:
			err := fmt.Errorf("%w: formula %q: unsupported token %s - %q", ErrParse, s, tok.Type, tok.Text)
			return nil, err
		}
	}
	if node == nil {
		return nil, fmt.Errorf("%w: empty formula", ErrParse)
	}
	return
}

func parseNode(s string) (node ast.Node, _ error) {
	for tok := range lex.Lex([]byte(s)) {
		switch tok.Type {
		case lex.TokenNumber:
			if node != nil {
				return nil, fmt.Errorf("%w: multiple values - %q", ErrParse, s)
			}
			node = &ast.NumberNode{Number: tok.Text}
		default:
			return nil, fmt.Errorf("%w: unsupported node value - %q", ErrParse, s)
		}
	}
	if node == nil {
		return nil, fmt.Errorf("%w: empty cell", ErrParse)
	}
	return
}
