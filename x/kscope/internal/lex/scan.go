// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lex

import (
	"errors"
	"fmt"
	"unicode"
)

// ErrScan means error scanning the next token.
var ErrScan = errors.New("scan error")

const (
	// keep-sorted start
	runeComma = ','
	runeDiv   = '/'
	runeDot   = '.'
	runeEOL   = '\n'
	runeEq    = '='
	runeHash  = '#'
	runeLpar  = '('
	runeMinus = '-'
	runeMul   = '*'
	runePlus  = '+'
	runeRpar  = ')'
	// keep-sorted end
)

// scanFunc is a scan state. Every scan state uses low-level read from the
// blockReader to scan the text to get the token and return the next expected
// scan state.
//
// A scan state returns an error with nil token and nil next scan state upon
// error.
type scanFunc func(br *blockReader) (*Token, scanFunc, error)

// singleRuneScanners maps single-run tokens to scan states.
var singleRuneScanners map[rune]scanFunc

func init() {
	singleRuneScanners = map[rune]scanFunc{
		// keep-sorted start
		runeComma: singleRuneScanner(TokComma),
		runeDiv:   singleRuneScanner(TokDiv),
		runeEq:    singleRuneScanner(TokAssign),
		runeHash:  scanComment,
		runeLpar:  singleRuneScanner(TokLpar),
		runeMinus: singleRuneScanner(TokMinus),
		runeMul:   singleRuneScanner(TokMul),
		runePlus:  singleRuneScanner(TokPlus),
		runeRpar:  singleRuneScanner(TokRpar),
		// keep-sorted end
	}
}

// scan is the initial and general state that skips whitespace and scans the
// next token.
func scan(br *blockReader) (*Token, scanFunc, error) {
	ignoreWhile(br, unicode.IsSpace)
	r, ok := br.Peek()
	if !ok {
		return nil, nil, nil
	}
	if scanner, ok := singleRuneScanners[r]; ok {
		return scanner(br)
	}
	// keep-sorted start skip_lines=1,-1
	switch {
	case unicode.IsDigit(r), r == runeDot:
		return scanNumber(br)
	case unicode.IsLetter(r):
		return scanIdentifier(br)
	}
	// keep-sorted end
	err := fmt.Errorf("%s: %w: unsupported character '%v'", br.Pos(), ErrScan, r)
	return nil, nil, err
}

// scanComment reads a comment token starting from current position to the end
// of the line.
func scanComment(br *blockReader) (*Token, scanFunc, error) {
	readWhile(br, func(r rune) bool {
		return r != runeEOL
	})
	return genToken(br, TokComment), scan, nil
}

func singleRuneScanner(tok TokenKind) scanFunc {
	return func(br *blockReader) (*Token, scanFunc, error) {
		br.Read() // consume the next rune
		return genToken(br, tok), scan, nil
	}
}

// reserved words.
var commands = map[string]TokenKind{
	"def":    TokDef, // function declaration
	"extern": TokExtern,
	"var":    TokVar, // variable declaration
}

// scanIdentifier scans an identifier.
//
//	ident  = letter [ alnum ]
//	letter = "a" .. "z" | "A" .. "Z"
//	alnum  = letter | digit
//	digit  = "0" .. "9"
func scanIdentifier(br *blockReader) (*Token, scanFunc, error) {
	readWhile(br, unicode.IsLetter)
	readWhile(br, func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r)
	})
	tok := genToken(br, TokIdent)
	if kind, ok := commands[tok.Val]; ok {
		tok.Kind = kind
	}
	return tok, scan, nil
}

// scanNumber scans a number token.
//
//	number = int | float
//	int    = digit { digit }
//	float  = int "." [ int ] | "." int
//	digit  = "0" .. "9"
func scanNumber(br *blockReader) (*Token, scanFunc, error) {
	readWhile(br, unicode.IsDigit)
	if r, ok := br.Peek(); ok && r == runeDot {
		br.Read() // skip the dot
		readWhile(br, unicode.IsDigit)
	}
	return genToken(br, TokNum), scan, nil
}

// genToken generates a token for current text block of kind tk.
func genToken(br *blockReader, tk TokenKind) *Token {
	s, pos := br.Text()
	return &Token{Kind: tk, Val: s, pos: pos}
}

// ignoreWhile ignores a block of consecutive characters that pass predicate f.
func ignoreWhile(br *blockReader, f func(rune) bool) {
	defer br.Ignore()
	readWhile(br, f)
}

// readWhile reads consecutive characters that pass predicate f.
func readWhile(br *blockReader, f func(rune) bool) {
	for {
		r, ok := br.Read()
		if !ok {
			break
		}
		if !f(r) {
			br.Unread()
			break
		}
	}
}
