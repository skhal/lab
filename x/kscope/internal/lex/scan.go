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

// ErrScan means error scanning next token.
var ErrScan = errors.New("scan error")

const runeDot = '.'

// scanFunc is a scan state. It reads data from the reader and returns parsed
// token along with the next scan state to process.
//
// It should return an error along with nil token and next state in case of a
// scan error.
type scanFunc func(rd *bufReader) (*Token, scanFunc, error)

func scan(rd *bufReader) (*Token, scanFunc, error) {
	ignoreWhile(rd, unicode.IsSpace)
	r, ok := rd.Peek()
	if !ok {
		return nil, nil, nil
	}
	// keep-sorted start skip_lines=1,-1
	switch {
	case unicode.IsDigit(r), r == runeDot:
		return scanNumber(rd)
	case unicode.IsLetter(r):
		return scanIdentifier(rd)
	}
	// keep-sorted end
	err := fmt.Errorf("%w: %d: unsupported character '%v'", ErrScan, rd.Pos(), r)
	return nil, nil, err
}

var commands = map[string]TokenKind{
	"def":    TokDef,
	"extern": TokExt,
}

// scanIdentifier scans an identifier.
//
//	ident  = letter [ alnum ]
//	letter = "a" .. "z" | "A" .. "Z"
//	alnum  = letter | digit
func scanIdentifier(rd *bufReader) (*Token, scanFunc, error) {
	readWhile(rd, unicode.IsLetter)
	readWhile(rd, func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r)
	})
	tok := genToken(rd, TokIdent)
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
func scanNumber(rd *bufReader) (*Token, scanFunc, error) {
	readWhile(rd, unicode.IsDigit)
	if r, ok := rd.Peek(); ok && r == runeDot {
		rd.Read() // skip dot
		readWhile(rd, unicode.IsDigit)
	}
	s, start, end := rd.Text()
	tok := &Token{Kind: TokNum, Val: s, Pos: Position{Start: start, End: end}}
	return tok, scan, nil
}

// genToken generates a token of specified kind using test and position from
// the reader.
func genToken(rd *bufReader, tk TokenKind) *Token {
	s, start, end := rd.Text()
	return &Token{Kind: tk, Val: s, Pos: Position{Start: start, End: end}}
}

// ignoreWhile ignores consecutive characters for which predicate f returns
// true.
func ignoreWhile(rc *bufReader, f func(rune) bool) {
	readWhile(rc, f)
	rc.Ignore()
}

// readWhile reads consecutive characters for which predicate f returns true.
// end of stream.
func readWhile(rd *bufReader, f func(rune) bool) {
	for {
		r, ok := rd.Read()
		if !ok {
			break
		}
		if !f(r) {
			rd.Unread()
			break
		}
	}
}
