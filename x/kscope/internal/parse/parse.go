// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parse

import (
	"errors"
	"fmt"
	"iter"
	"strconv"

	"github.com/skhal/lab/x/kscope/internal/ast"
	"github.com/skhal/lab/x/kscope/internal/lex"
)

// ErrParse means there is an error in parsing.
var ErrParse = errors.New("parse error")

// Parse returns an Abstract Syntax Tree (AST) representing parsed string s.
func Parse(s string) (ast.Node, error) {
	var lx lex.Lexer
	next, stop := iter.Pull(lx.Lex(s))
	defer stop()
	r := &peekerReader{reader: readerFunc(next)}
	par := &parser{r: r}
	n, err := par.Parse()
	if err != nil {
		return nil, err
	}
	if lx.Err() != nil {
		return nil, lx.Err()
	}
	return n, nil
}

// readerFunc adopts a function that returns next token to the [reader]
// interface.
type readerFunc func() (lex.Token, bool)

// Read returns the next token and a flag to indicate whether a token exists.
func (f readerFunc) Read() (lex.Token, bool) {
	return f()
}

// parser is a Recursive Descent Parser (RDP) to parse a sequence of tokens
// into an AST.
type parser struct {
	r *peekerReader
}

// Parse parses tokens into an AST.
func (p *parser) Parse() (ast.Node, error) {
	n, err := p.parse()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrParse, err)
	}
	return n, nil
}

func (p *parser) parse() (ast.Node, error) {
	if _, ok := p.r.Peek(); !ok {
		// end of stream
		return nil, nil
	}
	return p.parseExpression()
}

// parseExpression parses an expression.
func (p *parser) parseExpression() (ast.Node, error) {
	lhs, err := p.parseOperand()
	if err != nil {
		return nil, err
	}
	tok, ok := p.r.Peek()
	if !ok {
		return lhs, nil
	}
	switch tok.Kind {
	case lex.TokPlus, lex.TokMinus:
		return p.parseBinExpr(lhs)
	}
	return nil, fmt.Errorf("unsupported token %s", tok)
}

func (p *parser) parseOperand() (ast.Node, error) {
	tok, ok := p.r.Read()
	if !ok {
		return nil, fmt.Errorf("missing expression")
	}
	switch tok.Kind {
	case lex.TokNum:
		return parseNumber(tok)
	}
	return nil, fmt.Errorf("unsupported token %s", tok)
}

var binOps = map[lex.TokenKind]ast.BinOp{
	lex.TokPlus:  ast.BinOpPlus,
	lex.TokMinus: ast.BinOpMinus,
}

func (p *parser) parseBinExpr(lhs ast.Node) (ast.Node, error) {
	tok, ok := p.r.Read()
	if !ok {
		return nil, fmt.Errorf("missing binary operator")
	}
	op, ok := binOps[tok.Kind]
	if !ok {
		return nil, fmt.Errorf("unsupported binary operator %s", tok)
	}
	rhs, err := p.parseOperand()
	if err != nil {
		return nil, fmt.Errorf("operator %s: right operand: %s", tok, err)
	}
	return ast.BinExpr{Op: op, Left: lhs, Right: rhs}, nil
}

// parseNumber parses token as a numbee
func parseNumber(tok lex.Token) (ast.Node, error) {
	v, err := strconv.ParseFloat(tok.Val, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse number - %s", tok)
	}
	return ast.Number{Val: v}, nil
}
