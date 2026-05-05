// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parse

import (
	"errors"
	"fmt"
	"io"
	"iter"
	"strconv"

	"github.com/skhal/lab/x/kscope/internal/ast"
	"github.com/skhal/lab/x/kscope/internal/lex"
)

// ErrParse means there is an error in parsing.
var ErrParse = errors.New("parse error")

// errNotDeclaration means the root element is not a declaration.
var errNotDeclaration = errors.New("invalid root token") // NOEXPORT

// Parse parses the code segment s.
func Parse(s string) (ast.Node, error) {
	return parse(s, (*parser).ParseCode)
}

// ParseExpr parses a single expression. A declaration is a valid expression.
func ParseExpr(s string) (ast.Node, error) {
	return parse(s, (*parser).ParseExpr)
}

type parseFunc func(*parser) (ast.Node, error)

// parse uses [parser] method expression pf to parse code in s.
func parse(s string, pf parseFunc) (node ast.Node, err error) {
	var lx lex.Lexer
	defer func() {
		if err != nil {
			return
		}
		err = lx.Err()
	}()
	toks, _ := lx.Lex(s)
	next, stop := iter.Pull(toks)
	defer stop()
	r := newTokenReader(readerFunc(next))
	return pf(&parser{tr: r})
}

// parser is a Recursive Descent Parser (RDP). It converts a sequence of tokens
// into an Abstract Syntax Tree (AST).
type parser struct {
	tr *tokenReader
}

// ParseCode parses a code segment. A code segment only holds declarations.
// It is an exported API to wrap all internal errors into ErrParse.
func (p *parser) ParseCode() (ast.Node, error) {
	f, err := p.parseCode()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrParse, err)
	}
	return f, nil
}

// parseCode parse a code segment.
func (p *parser) parseCode() (*ast.Code, error) {
	var decls []*ast.Decl
	for {
		decl, err := p.parseDecl()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		decls = append(decls, &ast.Decl{Node: decl})
	}
	if len(decls) == 0 {
		return nil, nil
	}
	return &ast.Code{Decls: decls}, nil
}

// ParseExpr parses an expression. A declaration is a valid expression. It is
// a public API to wrap parsing errors into ErrParse.
func (p *parser) ParseExpr() (ast.Node, error) {
	n, err := p.parseDecl()
	switch {
	case err == io.EOF:
		return nil, nil
	case errors.Is(err, errNotDeclaration):
		n, err = p.parseExpression()
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrParse, err)
	}
	return n, err
}

// parseDecl parses a declaration of a variable or a function.
func (p *parser) parseDecl() (ast.Node, error) {
	tok, ok := p.tr.Peek()
	if !ok {
		// end of stream
		return nil, io.EOF
	}
	// keep-sorted start skip_lines=1,-1
	switch tok.Kind {
	case lex.TokDef:
		return p.parseFunc()
	case lex.TokVar:
		return p.parseVar()
	}
	// keep-sorted end
	return nil, fmt.Errorf("%s: %w", tok, errNotDeclaration)
}

// parseFunc parses a function definition.
func (p *parser) parseFunc() (ast.Node, error) {
	p.tr.Read() // skip TokDef
	ident, ok := p.tr.Read()
	if !ok || ident.Kind != lex.TokIdent {
		return nil, fmt.Errorf("missing function identifier")
	}
	if tok, ok := p.tr.Read(); !ok || tok.Kind != lex.TokLpar {
		return nil, fmt.Errorf("func %s: missing args left parenthesis", ident.Val)
	}
	params, err := p.parseFuncParams()
	if err != nil {
		return nil, fmt.Errorf("func %s: %s", ident.Val, err)
	}
	if tok, ok := p.tr.Read(); !ok || tok.Kind != lex.TokRpar {
		return nil, fmt.Errorf("func %s: missing args right parenthesis", ident.Val)
	}
	body, err := p.parseExpression()
	if err != nil {
		return nil, fmt.Errorf("func %s: parse body: %s", ident.Val, err)
	}
	node := ast.Func{
		Name: ident.Val,
		Body: []ast.Node{
			body,
		},
	}
	if len(params) != 0 {
		node.Params = params
	}
	return node, nil
}

// parseFuncParams parses function parameters
func (p *parser) parseFuncParams() ([]string, error) {
	var params []string
	for pnum := 1; ; pnum++ {
		if tok, ok := p.tr.Peek(); !ok || tok.Kind == lex.TokRpar {
			break
		}
		par, ok := p.tr.Read()
		if !ok || par.Kind != lex.TokIdent {
			return nil, fmt.Errorf("param %d: missing name", pnum)
		}
		params = append(params, par.Val)
		if tok, ok := p.tr.Peek(); ok && tok.Kind == lex.TokComma {
			p.tr.Read()
		}
	}
	return params, nil
}

// parseExpression parses an expression.
func (p *parser) parseExpression() (ast.Node, error) {
	lhs, err := p.parseOperand()
	if err != nil {
		return nil, err
	}
	if tok, ok := p.tr.Peek(); ok && isBinExprOperator(tok) {
		return p.parseBinExpr(lhs)
	}
	return lhs, nil
}

func isBinExprOperator(tok lex.Token) bool {
	switch tok.Kind {
	case lex.TokPlus, lex.TokMinus, lex.TokMul, lex.TokDiv:
		return true
	}
	return false
}

func (p *parser) parseOperand() (ast.Node, error) {
	tok, ok := p.tr.Read()
	if !ok {
		return nil, fmt.Errorf("missing expression")
	}
	// keep-sorted start skip_lines=1,-1
	switch tok.Kind {
	case lex.TokIdent:
		if next, ok := p.tr.Peek(); ok && next.Kind == lex.TokLpar {
			return p.parseCall(tok)
		}
		return ast.Ident{Name: tok.Val}, nil
	case lex.TokLpar:
		return p.parseGroup(tok)
	case lex.TokNum:
		return parseNumber(tok)
	}
	// keep-sorted end
	return nil, fmt.Errorf("unsupported token %s", tok)
}

func (p *parser) parseCall(ident lex.Token) (ast.Node, error) {
	// left parenthesis
	if _, ok := p.tr.Read(); !ok {
		return nil, fmt.Errorf("call %s: missing left parenthesis", ident.Val)
	}
	var args []ast.Node
	for {
		if tok, ok := p.tr.Peek(); !ok || tok.Kind == lex.TokRpar {
			break
		}
		arg, err := p.parseExpression()
		if err != nil {
			return nil, fmt.Errorf("call %s: %s", ident.Val, err)
		}
		args = append(args, arg)
		if tok, ok := p.tr.Peek(); ok && tok.Kind == lex.TokComma {
			p.tr.Read()
		}
	}
	if tok, ok := p.tr.Read(); !ok || tok.Kind != lex.TokRpar {
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
	tok, ok := p.tr.Read()
	if !ok {
		return nil, fmt.Errorf("missing binary operator")
	}
	op, ok := binOps[tok.Kind]
	if !ok {
		return nil, fmt.Errorf("unsupported binary operator %s", tok)
	}
	next, ok := p.tr.Peek()
	if !ok {
		return nil, fmt.Errorf("missing right operand")
	}
	rhs, err := p.parseExpression()
	if err != nil {
		return nil, fmt.Errorf("operator %s: right operand: %s", tok, err)
	}
	n := ast.BinExpr{Op: op, Left: lhs, Right: rhs}
	if shouldRotateLeft(tok, next) {
		n = rotateLeft(n)
	}
	return n, nil
}

func shouldRotateLeft(tok, next lex.Token) bool {
	if next.Kind == lex.TokLpar {
		return false
	}
	return tok.Kind == lex.TokDiv || tok.Kind == lex.TokMul
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

func (p *parser) parseGroup(_ lex.Token) (ast.Node, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	if tok, ok := p.tr.Read(); !ok || tok.Kind != lex.TokRpar {
		return nil, fmt.Errorf("missing right paraenthesis")
	}
	return expr, nil
}

func (p *parser) parseVar() (ast.Node, error) {
	p.tr.Read() // skip TokVar
	ident, ok := p.tr.Read()
	if !ok || ident.Kind != lex.TokIdent {
		return nil, fmt.Errorf("missing variable identifier")
	}
	if tok, ok := p.tr.Read(); !ok || tok.Kind != lex.TokAssign {
		return nil, fmt.Errorf("var %s: missing assignment", ident.Val)
	}
	val, err := p.parseExpression()
	if err != nil {
		return nil, fmt.Errorf("var %s: parse value: %s", ident.Val, err)
	}
	n := ast.Var{Name: ident.Val, Val: val}
	return n, nil
}
