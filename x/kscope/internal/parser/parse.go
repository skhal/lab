// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"errors"
	"fmt"
	"io"
	"iter"
	"strconv"

	"github.com/skhal/lab/x/kscope/internal/ast"
	"github.com/skhal/lab/x/kscope/internal/lex"
)

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
	toks, posr := lx.Lex(s)
	next, stop := iter.Pull(toks)
	defer stop()
	r := newTokenReader(readerFunc(next))
	p := newParser(r, posr)
	return pf(p)
}

// parser is a Recursive Descent Parser (RDP). It converts a sequence of tokens
// into an Abstract Syntax Tree (AST).
type parser struct {
	tr   *tokenReader
	posr *lex.Positioner
}

// newParser creates a parser with supplied token reader.
func newParser(r *tokenReader, p *lex.Positioner) *parser {
	return &parser{tr: r, posr: p}
}

// ParseCode parses a code segment. A code segment only holds declarations.
// It is an exported API to wrap all internal errors into ErrParse.
func (p *parser) ParseCode() (ast.Node, error) {
	f, err := p.parseCode()
	switch {
	case errors.Is(err, ErrParse):
		return nil, err
	case err != nil:
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
	case errors.Is(err, ErrParse):
		return nil, err
	case err != nil:
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
	return nil, p.error(errNotDeclaration)
}

// parseFunc parses a function definition.
func (p *parser) parseFunc() (ast.Node, error) {
	p.tr.Read() // skip TokDef
	ident, ok := p.tr.Read()
	switch {
	case !ok:
		return nil, p.errorf("missing function identifier")
	case ident.Kind != lex.TokIdent:
		return nil, p.errorf("not a function identifier")
	}
	switch tok, ok := p.tr.Read(); {
	case !ok:
		return nil, p.errorf("missing left parenthesis")
	case tok.Kind != lex.TokLpar:
		return nil, p.errorf("not left parenthesis")
	}
	params, err := p.parseFuncParams()
	if err != nil {
		return nil, err
	}
	switch tok, ok := p.tr.Read(); {
	case !ok:
		return nil, p.errorf("missing right parenthesis")
	case tok.Kind != lex.TokRpar:
		return nil, p.errorf("not right parenthesis")
	}
	if _, ok := p.tr.Peek(); !ok {
		return nil, p.errorf("missing function body")
	}
	body, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	node := &ast.Func{
		Name: ident.Val,
		Body: []ast.Node{
			body,
		},
		Params: params,
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
		par, _ := p.tr.Read()
		if par.Kind != lex.TokIdent {
			return nil, p.errorf("not a parameter name")
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

// parseOperand parses an operand of an expression. It might be a number,
// identifier, a function call, or a grouped expression enclosed in
// parentheses.
func (p *parser) parseOperand() (ast.Node, error) {
	tok, ok := p.tr.Read()
	if !ok {
		return nil, p.errorf("missing operand")
	}
	// keep-sorted start skip_lines=1,-1
	switch tok.Kind {
	case lex.TokIdent:
		if next, ok := p.tr.Peek(); ok && next.Kind == lex.TokLpar {
			return p.parseCall(tok)
		}
		return &ast.Ident{Name: tok.Val}, nil
	case lex.TokLpar:
		return p.parseGroup(tok)
	case lex.TokNum:
		return p.parseNumber(tok)
	}
	// keep-sorted end
	return nil, p.errorf("not operand %s", tok.Val)
}

// parseCall parses a function call.
func (p *parser) parseCall(ident lex.Token) (ast.Node, error) {
	// left parenthesis
	if _, ok := p.tr.Read(); !ok {
		return nil, p.errorf("missing left parenthesis")
	}
	args, err := p.parseCallArgs()
	if err != nil {
		return nil, err
	}
	if tok, ok := p.tr.Read(); !ok || tok.Kind != lex.TokRpar {
		return nil, p.errorf("missing right parenthesis")
	}
	node := &ast.Call{
		Name: ident.Val,
		Args: args,
	}
	return node, nil
}

// parseCallArgs parses a function call arguments
func (p *parser) parseCallArgs() ([]ast.Node, error) {
	var args []ast.Node
	for {
		if tok, ok := p.tr.Peek(); !ok || tok.Kind == lex.TokRpar {
			break
		}
		arg, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		if tok, ok := p.tr.Peek(); ok && tok.Kind == lex.TokComma {
			p.tr.Read()
		}
	}
	return args, nil
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
		return nil, p.errorf("missing binary operator")
	}
	op, ok := binOps[tok.Kind]
	if !ok {
		return nil, p.errorf("invalid operator %s", tok)
	}
	next, ok := p.tr.Peek()
	if !ok {
		return nil, p.errorf("missing right operand")
	}
	rhs, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	n := &ast.BinExpr{Op: op, Left: lhs, Right: rhs}
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
func rotateLeft(n *ast.BinExpr) *ast.BinExpr {
	right, ok := n.Right.(*ast.BinExpr)
	if !ok {
		return n
	}
	n.Right = right.Left
	right.Left = n
	return right
}

// parseNumber parses token as a numbee
func (p *parser) parseNumber(tok lex.Token) (ast.Node, error) {
	v, err := strconv.ParseFloat(tok.Val, 64)
	if err != nil {
		return nil, p.errorf("not a number - %s", tok.Val)
	}
	return &ast.Number{Val: v}, nil
}

func (p *parser) parseGroup(_ lex.Token) (ast.Node, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	if tok, ok := p.tr.Read(); !ok || tok.Kind != lex.TokRpar {
		return nil, p.errorf("missing right paraenthesis")
	}
	return expr, nil
}

func (p *parser) parseVar() (ast.Node, error) {
	p.tr.Read() // skip TokVar
	ident, ok := p.tr.Read()
	switch {
	case !ok:
		return nil, p.errorf("missing variable identifier")
	case ident.Kind != lex.TokIdent:
		return nil, p.errorf("want variable identifier")
	}
	switch tok, ok := p.tr.Read(); {
	case !ok:
		return nil, p.errorf("missing assignment")
	case tok.Kind != lex.TokAssign:
		return nil, p.errorf("want assignment")
	}
	val, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	n := &ast.Var{Name: ident.Val, Val: val}
	return n, nil
}

func (p *parser) error(err error) error {
	return newParseError(p.posr.Pos(), err)
}

func (p *parser) errorf(format string, arg ...any) error {
	err := fmt.Errorf(format, arg...)
	return p.error(err)
}
