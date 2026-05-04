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
	switch tok, ok := p.r.Peek(); {
	case !ok:
		// end of stream
		return nil, nil
	case tok.Kind == lex.TokDef:
		return p.parseFunc()
	}
	return p.parseExpression()
}

// parseFunc parses a function definition.
func (p *parser) parseFunc() (ast.Node, error) {
	p.r.Read() // skip TokDef
	ident, ok := p.r.Read()
	if !ok || ident.Kind != lex.TokIdent {
		return nil, fmt.Errorf("missing function identifier")
	}
	if tok, ok := p.r.Read(); !ok || tok.Kind != lex.TokLpar {
		return nil, fmt.Errorf("func %s: missing args left parenthesis", ident.Val)
	}
	if tok, ok := p.r.Read(); !ok || tok.Kind != lex.TokRpar {
		return nil, fmt.Errorf("func %s: missing args right parenthesis", ident.Val)
	}
	body, err := p.parseExpression()
	if err != nil {
		return nil, fmt.Errorf("func %s: parse body: %s", ident.Val, err)
	}
	n := ast.Func{
		Name: ident.Val,
		Body: []ast.Node{
			body,
		},
	}
	return n, nil
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
	case lex.TokPlus, lex.TokMinus, lex.TokMul, lex.TokDiv:
		return p.parseBinExpr(lhs)
	}
	return lhs, nil
}

func (p *parser) parseOperand() (ast.Node, error) {
	tok, ok := p.r.Read()
	if !ok {
		return nil, fmt.Errorf("missing expression")
	}
	switch tok.Kind {
	case lex.TokIdent:
		if next, ok := p.r.Peek(); !ok || next.Kind != lex.TokLpar {
			break
		}
		return p.parseCall(tok)
	case lex.TokNum:
		return parseNumber(tok)
	}
	return nil, fmt.Errorf("unsupported token %s", tok)
}

func (p *parser) parseCall(ident lex.Token) (ast.Node, error) {
	// left parenthesis
	if _, ok := p.r.Read(); !ok {
		return nil, fmt.Errorf("call %s: missing left parenthesis", ident.Val)
	}
	var args []ast.Node
	for {
		if tok, ok := p.r.Peek(); !ok || tok.Kind == lex.TokRpar {
			break
		}
		arg, err := p.parseExpression()
		if err != nil {
			return nil, fmt.Errorf("call %s: %s", ident.Val, err)
		}
		args = append(args, arg)
		if tok, ok := p.r.Peek(); ok && tok.Kind == lex.TokComma {
			// ignore comma
			p.r.Read()
		}
	}
	if tok, ok := p.r.Read(); !ok || tok.Kind != lex.TokRpar {
		return nil, fmt.Errorf("call %s: missing right parenthesis", ident.Val)
	}
	n := ast.Call{
		Name: ident.Val,
		Args: args,
	}
	return n, nil
}

var binOps = map[lex.TokenKind]ast.BinOp{
	// keep-sorted start
	lex.TokDiv:   ast.BinOpDiv,
	lex.TokMinus: ast.BinOpMinus,
	lex.TokMul:   ast.BinOpMul,
	lex.TokPlus:  ast.BinOpPlus,
	// keep-sorted end
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
	rhs, err := p.parseExpression()
	if err != nil {
		return nil, fmt.Errorf("operator %s: right operand: %s", tok, err)
	}
	n := ast.BinExpr{Op: op, Left: lhs, Right: rhs}
	if op == ast.BinOpDiv || op == ast.BinOpMul {
		// give preference to the first operator op1 in "a op1 b op2 c", e.g.:
		// "1 * 2 + 3" becomes {Op:Plus Left:{Op:Mul Left:1 Right:2} Right:3}
		// "1 * 2 / 3" becomes {Op:Div Left:{Op:Mul Left:1 Right:2} Right:3}
		n = rotateLeft(n)
	}
	return n, nil
}

// rotateLeft rotates nodes in a binary expression counter-clockwise
// (aka right-hand rule) to make BinExpr.Right root node if is it a binary
// expression.
func rotateLeft(n ast.BinExpr) ast.BinExpr {
	right, ok := n.Right.(ast.BinExpr)
	if !ok {
		return n
	}
	n.Right = right.Left
	right.Left = n
	return right
}

// parseNumber parses token as a numbee
func parseNumber(tok lex.Token) (ast.Node, error) {
	v, err := strconv.ParseFloat(tok.Val, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse number - %s", tok)
	}
	return ast.Number{Val: v}, nil
}
