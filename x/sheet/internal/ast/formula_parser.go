// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

import (
	"fmt"
	"iter"
	"strings"

	"github.com/skhal/lab/x/sheet/internal/lex"
)

// formulaParser parses formula without "=" prefix into an AST tree.
//
// Context free grammar (CFG):
//
//	Expr       = Operand | BinaryExpr
//
//	Operand    = Number | Reference | Range | Call | IfCall | Group
//
//	Number     = Digit { Digit } [ "." { Digit } ]
//	Digit      = "0" .. "9"
//
//	Reference  = Letter Digit { Digit }
//	Letter     = "A" .. "Z"
//
//	Range      = Reference ":" Reference
//
//	Call       = Identifier "(" [ ArgList ] ")"
//	Identifier = Letter { Letter }
//	ArgList    = Expr { "," ArgsList }
//
//	IfCall     = "IF" "(" IfArgs ")"
//	IfArgs     = RelExpr "," Expr "," Expr
//	RelExpr    = Expr RelOp Expr
//	RelOp      = "==" | "!=" | "<" | "<=" | ">" | ">="
//
//	Group      =  "(" Expr ")"
//
//	BinaryExpr = Expr Op Expr
//	Op         = AddOp | MulOp
//	AddOp      = "+" | "-"
//	MulOp      = "*" | "/"
type formulaParser struct {
	next func() (lex.Token, bool)
	peek func() (lex.Token, bool)
}

// Parse parses formula b and returns root AST node.
func (p *formulaParser) Parse(b []byte) (Node, error) {
	next, stop := iter.Pull(lex.Lex(b))
	defer stop()
	cache := struct {
		valid bool
		tok   lex.Token
		ok    bool
	}{}
	p.next = func() (lex.Token, bool) {
		if cache.valid {
			cache.valid = false
			return cache.tok, cache.ok
		}
		return next()
	}
	p.peek = func() (lex.Token, bool) {
		if cache.valid {
			return cache.tok, cache.ok
		}
		cache.tok, cache.ok = next()
		cache.valid = true
		return cache.tok, cache.ok
	}
	n, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	tok, ok := p.next()
	if ok {
		return nil, fmt.Errorf("unexpected token %s", tok.Text)
	}
	return n, nil
}

func (p *formulaParser) parseExpr() (Node, error) {
	lhs, err := p.parseOperand()
	if err != nil {
		return nil, err
	}
	op, ok := p.peek()
	if !ok {
		return lhs, nil
	}
	switch op.Type {
	case lex.TokenPlus, lex.TokenMinus, lex.TokenMultiply, lex.TokenDivide:
		return p.parseBinaryExpr(lhs)
	}
	return lhs, nil
}

func (p *formulaParser) parseBinaryExpr(lhs Node) (Node, error) {
	op, _ := p.next() // the operator is guaranteed by the called
	opNext, okNext := p.peek()
	rhs, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	n := &BinOpNode{
		Op:    op.Text,
		Left:  lhs,
		Right: rhs,
	}
	switch op.Type {
	case lex.TokenMultiply, lex.TokenDivide:
		if okNext && opNext.Type == lex.TokenLpar {
			break
		}
		return rotateLeft(n), nil
	}
	return n, nil
}

func rotateLeft(n *BinOpNode) Node {
	rhs, ok := n.Right.(*BinOpNode)
	if !ok {
		return n
	}
	n.Right, rhs.Left = rhs.Left, n
	return rhs
}

func (p *formulaParser) parseOperand() (Node, error) {
	tok, ok := p.next()
	if !ok {
		return nil, fmt.Errorf("expected operand")
	}
	switch tok.Type {
	case lex.TokenError:
		return nil, tok.Err
	case lex.TokenNumber:
		return p.parseNumber(tok)
	case lex.TokenIdent:
		if next, ok := p.peek(); ok && next.Type == lex.TokenLpar {
			return p.parseCall(tok)
		}
		return p.parseReference(tok)
	case lex.TokenRange:
		return p.parseRange(tok)
	case lex.TokenLpar:
		return p.parseGroup()
	}
	return nil, fmt.Errorf("invalid token - %s", tok.Type)
}

func (p *formulaParser) parseNumber(tok lex.Token) (Node, error) {
	if !numRx.MatchString(tok.Text) {
		return nil, fmt.Errorf("not a number %q", tok.Text)
	}
	return &NumberNode{Number: tok.Text}, nil
}

func (p *formulaParser) parseCall(ident lex.Token) (Node, error) {
	if !callRx.MatchString(ident.Text) {
		return nil, fmt.Errorf("invalid function name %s", ident.Text)
	}
	const ifCall = "IF"
	if ident.Text == ifCall {
		return p.parseIfCall()
	}
	return p.parseFuncCall(ident)
}

func (p *formulaParser) parseFuncCall(ident lex.Token) (Node, error) {
	p.next() // ignore left-parenthesis
	var args []Node
	if tok, ok := p.peek(); ok && tok.Type != lex.TokenRpar {
		var err error
		args, err = p.parseArgs()
		if err != nil {
			return nil, err
		}
	}
	if tok, ok := p.next(); !ok || tok.Type != lex.TokenRpar {
		return nil, fmt.Errorf("unbalanced parentheses")
	}
	return &CallNode{Name: ident.Text, Args: args}, nil
}

func (p *formulaParser) parseArgs() ([]Node, error) {
	var args []Node
	for {
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		if tok, ok := p.peek(); !ok || tok.Type != lex.TokenComma {
			break
		}
		p.next() // ignore comma
	}
	return args, nil
}

func (p *formulaParser) parseReference(ident lex.Token) (Node, error) {
	if !refRx.MatchString(ident.Text) {
		return nil, fmt.Errorf("invalid identifier %s", ident.Text)
	}
	return &RefNode{Ref: ident.Text}, nil
}

func (p *formulaParser) parseRange(tok lex.Token) (Node, error) {
	const (
		sep    = ":"
		fields = 2
	)
	items := strings.SplitN(tok.Text, sep, fields)
	if len(items) != fields {
		// should not happen as long as the lexer and the parser are in sync.
		return nil, fmt.Errorf("invalid range %s", tok.Text)
	}
	return &RangeNode{From: items[0], To: items[1]}, nil
}

func (p *formulaParser) parseGroup() (Node, error) {
	n, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if tok, ok := p.next(); !ok || tok.Type != lex.TokenRpar {
		return nil, fmt.Errorf("unbalanced parentheses")
	}
	return n, nil
}

func (p *formulaParser) parseIfCall() (Node, error) {
	p.next() // ignore left-parenthesis
	cond, ifPass, ifFail, err := p.parseIfArgs()
	if err != nil {
		return nil, err
	}
	if tok, ok := p.next(); !ok || tok.Type != lex.TokenRpar {
		return nil, fmt.Errorf("if: unbalanced parentheses")
	}
	return &IfNode{Cond: cond, IfPass: ifPass, IfFail: ifFail}, nil
}

func (p *formulaParser) parseIfArgs() (*RelOpNode, Node, Node, error) {
	cond, err := p.parseRelExpr()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("if condition: %s", err)
	}
	if tok, ok := p.next(); !ok || tok.Type != lex.TokenComma {
		return nil, nil, nil, fmt.Errorf("missing comma after if condition")
	}
	ifPass, err := p.parseExpr()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("if pass branch: %s", err)
	}
	if tok, ok := p.next(); !ok || tok.Type != lex.TokenComma {
		return nil, nil, nil, fmt.Errorf("missing comma after if pass branch")
	}
	ifFail, err := p.parseExpr()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("if fail branch: %s", err)
	}
	return cond, ifPass, ifFail, nil
}

func (p *formulaParser) parseRelExpr() (*RelOpNode, error) {
	lhs, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	op, ok := p.next()
	if !ok {
		return nil, fmt.Errorf("missing operator in binary expression")
	}
	switch op.Type {
	case lex.TokenEqual, lex.TokenNotEqual:
	case lex.TokenLess, lex.TokenLessOrEqual:
	case lex.TokenGreater, lex.TokenGreaterOrEqual:
	default:
		err := fmt.Errorf("unexpected operator %s - want compare operator", op.Type)
		return nil, err
	}
	rhs, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	n := &RelOpNode{
		Op:    op.Text,
		Left:  lhs,
		Right: rhs,
	}
	return n, nil
}
