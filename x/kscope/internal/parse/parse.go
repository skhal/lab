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

var errRootToken = errors.New("invalid root token") // NOEXPORT

// Parse parses the code s as a file content with all declarations.
func Parse(s string) (ast.Node, error) {
	return parse(s, (*parser).ParseCode)
}

// ParseExpr parses a single expression, e.g. a declaration, binary expression,
// etc.
func ParseExpr(s string) (ast.Node, error) {
	return parse(s, (*parser).ParseExpr)
}

func parse(s string, pf func(*parser) (ast.Node, error)) (ast.Node, error) {
	var lx lex.Lexer
	toks, _ := lx.Lex(s)
	next, stop := iter.Pull(toks)
	defer stop()
	r := &tokenReader{reader: readerFunc(next)}
	n, err := pf(&parser{r: r})
	if err != nil {
		return nil, err
	}
	if lx.Err() != nil {
		return nil, lx.Err()
	}
	return n, nil
}

// parser is a Recursive Descent Parser (RDP) to parse a sequence of tokens
// into an AST.
type parser struct {
	r *tokenReader
}

// ParseCode parses tokens into an AST.
func (p *parser) ParseCode() (ast.Node, error) {
	f, err := p.parseCode()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrParse, err)
	}
	return f, nil
}

func (p *parser) parseCode() (*ast.Code, error) {
	var decls []*ast.Decl
	for {
		n, err := p.parseDecl()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		decls = append(decls, &ast.Decl{Node: n})
	}
	var f *ast.Code
	if len(decls) != 0 {
		f = &ast.Code{Decls: decls}
	}
	return f, nil
}

func (p *parser) parseDecl() (ast.Node, error) {
	tok, ok := p.r.Peek()
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
	return nil, fmt.Errorf("%s: %w", tok, errRootToken)
}

// ParseExpr parses an expression, e.g. a declaration, binary expression, etc.
func (p *parser) ParseExpr() (ast.Node, error) {
	n, err := p.parseDecl()
	switch {
	case err == io.EOF:
		return nil, nil
	case errors.Is(err, errRootToken):
		n, err = p.parseExpression()
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrParse, err)
	}
	return n, err
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
	var params []string
	for pnum := 1; ; pnum++ {
		if tok, ok := p.r.Peek(); !ok || tok.Kind == lex.TokRpar {
			break
		}
		par, ok := p.r.Read()
		if !ok || par.Kind != lex.TokIdent {
			return nil, fmt.Errorf("func %s: parse param %d: missing name", ident.Val, pnum)
		}
		params = append(params, par.Val)
		if tok, ok := p.r.Peek(); ok && tok.Kind == lex.TokComma {
			p.r.Read()
		}
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
	if len(params) != 0 {
		n.Params = params
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
	// keep-sorted start skip_lines=1,-1
	switch tok.Kind {
	case lex.TokIdent:
		if next, ok := p.r.Peek(); ok && next.Kind == lex.TokLpar {
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
	next, ok := p.r.Peek()
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
	if tok, ok := p.r.Read(); !ok || tok.Kind != lex.TokRpar {
		return nil, fmt.Errorf("missing right paraenthesis")
	}
	return expr, nil
}

func (p *parser) parseVar() (ast.Node, error) {
	p.r.Read() // skip TokVar
	ident, ok := p.r.Read()
	if !ok || ident.Kind != lex.TokIdent {
		return nil, fmt.Errorf("missing variable identifier")
	}
	if tok, ok := p.r.Read(); !ok || tok.Kind != lex.TokAssign {
		return nil, fmt.Errorf("var %s: missing assignment", ident.Val)
	}
	val, err := p.parseExpression()
	if err != nil {
		return nil, fmt.Errorf("var %s: parse value: %s", ident.Val, err)
	}
	n := ast.Var{Name: ident.Val, Val: val}
	return n, nil
}
