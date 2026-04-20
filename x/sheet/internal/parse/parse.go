// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package parse implements lexical analysis of cell contents.
package parse

import (
	"errors"
	"fmt"
	"iter"
	"regexp"
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
	s, ok := strings.CutPrefix(s, formulaPrefix)
	s = strings.TrimSpace(s)
	if !ok {
		return parseCell(s)
	}
	return parseFormula(s)
}

var cellRx = regexp.MustCompile(`^\d+(?:\.\d*)?$`)

// parseCell parses a cell without formula.
func parseCell(s string) (ast.Node, error) {
	if len(s) == 0 {
		return nil, fmt.Errorf("%w: empty cell", ErrParse)
	}
	if !cellRx.MatchString(s) {
		return nil, fmt.Errorf("%w: invalid cell", ErrParse)
	}
	return &ast.NumberNode{Number: s}, nil
}

func parseFormula(s string) (ast.Node, error) {
	if len(s) == 0 {
		return nil, fmt.Errorf("%w: empty formula", ErrParse)
	}
	p := &formulaParser{}
	return p.Parse([]byte(s))
}

// formulaParser parses formula without "=" prefix into an AST tree.
//
// Context free grammar (CFG):
//
//	Expr       = Operand | BinaryExpr
//	Operand    = Number | "(" Expr ")"
//	BinaryExpr = Expr Op Expr
//	Op         = "+" | "-"
type formulaParser struct {
	next  func() (lex.Token, bool)
	depth int // parentheses level
}

// Parse parses formula b and returns root AST node.
func (p *formulaParser) Parse(b []byte) (ast.Node, error) {
	var stop func()
	p.next, stop = iter.Pull(lex.Lex(b))
	defer stop()
	return p.parseExpr()
}

func (p *formulaParser) parseExpr() (ast.Node, error) {
	lhs, err := p.parseOperand()
	if err != nil {
		return nil, err
	}
	op, ok := p.next()
	if !ok {
		// no more tokens
		return lhs, nil
	}
	switch op.Type {
	case lex.TokenPlus, lex.TokenMinus:
		return p.parseBinaryExpr(lhs, op)
	case lex.TokenRpar:
		if p.depth == 0 {
			return nil, fmt.Errorf("%w: unbalanced right parenthesis", ErrParse)
		}
		p.depth--
		return lhs, nil
	}
	return nil, ErrParse
}

func (p *formulaParser) parseOperand() (ast.Node, error) {
	tok, ok := p.next()
	if !ok {
		// no more tokens
		return nil, fmt.Errorf("%w: expected operand", ErrParse)
	}
	parseLpar := func() (ast.Node, error) {
		depth := p.depth
		p.depth++
		n, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if p.depth != depth {
			return nil, fmt.Errorf("%w: unbalances parentheses", ErrParse)
		}
		return n, nil
	}
	switch tok.Type {
	case lex.TokenError:
		return nil, fmt.Errorf("%w: %s", ErrParse, tok.Err)
	case lex.TokenNumber:
		return &ast.NumberNode{Number: tok.Text}, nil
	case lex.TokenLpar:
		return parseLpar()
	}
	return nil, fmt.Errorf("%w: invalid token - %s", ErrParse, tok.Type)
}

func (p *formulaParser) parseBinaryExpr(lhs ast.Node, op lex.Token) (ast.Node, error) {
	// binary expression
	rhs, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return &ast.BinOpNode{
		Op:    op.Text,
		Left:  lhs,
		Right: rhs,
	}, nil
}
