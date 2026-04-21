// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parse

import (
	"fmt"
	"iter"
	"strings"

	"github.com/skhal/lab/x/sheet/internal/ast"
	"github.com/skhal/lab/x/sheet/internal/lex"
)

// formulaParser parses formula without "=" prefix into an AST tree.
//
// Context free grammar (CFG):
//
//	Expr       = Operand | BinaryExpr
//
//	Operand    = Number | Identifier | Range | Call | "(" Expr ")"
//	Identifier = Letter Digit
//	Range      = Identifier ":" Identifier
//	Call       = Identifier "(" [ ArgList ] ")"
//	ArgList    = Expr { "," ArgsList }
//
//	BinaryExpr = Expr Op Expr
//	Op         = "+" | "-"
type formulaParser struct {
	next func() (lex.Token, bool)
	peek func() (lex.Token, bool)
}

// Parse parses formula b and returns root AST node.
func (p *formulaParser) Parse(b []byte) (ast.Node, error) {
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

func (p *formulaParser) parseExpr() (ast.Node, error) {
	lhs, err := p.parseOperand()
	if err != nil {
		return nil, err
	}
	op, ok := p.peek()
	if !ok {
		return lhs, nil
	}
	switch op.Type {
	case lex.TokenPlus, lex.TokenMinus:
		p.next() // discard peek cache
		return p.parseBinaryExpr(lhs, op)
	}
	return lhs, nil
}

func (p *formulaParser) parseOperand() (ast.Node, error) {
	tok, ok := p.next()
	if !ok {
		return nil, fmt.Errorf("expected operand")
	}
	switch tok.Type {
	// keep-sorted start
	case lex.TokenError:
		return nil, tok.Err
	case lex.TokenIdent: // identifier or a function call
		next, ok := p.peek()
		if ok && next.Type == lex.TokenLpar {
			return p.parseCall(tok)
		}
		return p.parseIdentifier(tok)
	case lex.TokenLpar: // an expression in parentheses
		n, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if tok, ok = p.next(); !ok || tok.Type != lex.TokenRpar {
			return nil, fmt.Errorf("unbalanced parentheses")
		}
		return n, nil
	case lex.TokenNumber:
		return &ast.NumberNode{Number: tok.Text}, nil
	case lex.TokenRange:
		return p.parseRange(tok)
		// keep-sorted end
	}
	return nil, fmt.Errorf("invalid token - %s", tok.Type)
}

func (p *formulaParser) parseIdentifier(ident lex.Token) (ast.Node, error) {
	if !refRx.MatchString(ident.Text) {
		return nil, fmt.Errorf("invalid identifier %s", ident.Text)
	}
	return &ast.RefNode{Ref: ident.Text}, nil
}

func (p *formulaParser) parseRange(tok lex.Token) (ast.Node, error) {
	const (
		sep    = ":"
		fields = 2
	)
	items := strings.SplitN(tok.Text, sep, fields)
	if len(items) != fields {
		// should not happen as long as the lexer and the parser are in sync.
		return nil, fmt.Errorf("invalid range %s", tok.Text)
	}
	return &ast.RangeNode{From: items[0], To: items[1]}, nil
}

func (p *formulaParser) parseCall(ident lex.Token) (ast.Node, error) {
	if !callRx.MatchString(ident.Text) {
		return nil, fmt.Errorf("invalid function name %s", ident.Text)
	}
	p.next() // skip left-parenthesis
	var args []ast.Node
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
	return &ast.CallNode{Name: ident.Text, Args: args}, nil
}

func (p *formulaParser) parseArgs() ([]ast.Node, error) {
	var args []ast.Node
	for {
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		tok, ok := p.peek()
		if !ok {
			break
		}
		if tok.Type != lex.TokenComma {
			break
		}
		p.next()
	}
	return args, nil
}

func (p *formulaParser) parseBinaryExpr(lhs ast.Node, op lex.Token) (ast.Node, error) {
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
