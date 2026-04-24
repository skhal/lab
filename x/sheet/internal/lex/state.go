// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lex

import (
	"bytes"
	"fmt"
	"io"
	"unicode"
)

type stateFunc func(*lexer) stateFunc

var whitespace = []byte(` \t`)

func isWhitespace(r rune) bool { return bytes.ContainsRune(whitespace, r) }

// scanState is the default state of the scanner. It skips whitespace and
// advances to the next supported state.
func scanState(lx *lexer) stateFunc {
	const (
		// keep-sorted start
		comma    = ','
		divide   = '/'
		lpar     = '('
		minus    = '-'
		multiply = '*'
		plus     = '+'
		rpar     = ')'
		// keep-sorted end
	)
	lx.ScanFunc(isWhitespace)
	lx.Ignore()
	switch r, err := lx.Peek(); {
	case err == io.EOF:
		return eofState
	case err != nil:
		return errorState(err)
	case unicode.IsNumber(r):
		return numberState
	case unicode.IsLetter(r):
		return identifierState
	// keep-sorted start
	case r == comma:
		return genState(TokenComma)
	case r == divide:
		return genState(TokenDivide)
	case r == lpar:
		return genState(TokenLpar)
	case r == minus:
		return genState(TokenMinus)
	case r == multiply:
		return genState(TokenMultiply)
	case r == plus:
		return genState(TokenPlus)
	case r == rpar:
		return genState(TokenRpar)
	// keep-sorted end
	default:
		err := fmt.Errorf("unsupported text at %d -  %q", lx.pos, lx.b[lx.pos:])
		return errorState(err)
	}
}

// numberState parses a floating value number with non-empty integral part. It
// emits parsed number token and advances to the scanState.
func numberState(lx *lexer) stateFunc {
	lx.ScanFunc(unicode.IsDigit)
	switch r, err := lx.Peek(); {
	case err != nil:
		// read failed - the next state will handle the error
	case r == '.':
		lx.Read()
		lx.ScanFunc(unicode.IsDigit)
	}
	lx.Emit(TokenNumber)
	return scanState
}

// identifierState parses an identifier that consists of any number of letters
// optionally followed by numbers, e.g. "abC123" or "Ab". The function is case
// insensitive.
func identifierState(lx *lexer) stateFunc {
	lx.ScanFunc(unicode.IsLetter)

	if r, err := lx.Peek(); err != nil || !unicode.IsNumber(r) {
		lx.Emit(TokenIdent)
		return scanState
	}
	lx.ScanFunc(unicode.IsNumber)

	const colon = ':'
	if r, err := lx.Peek(); err != nil || r != colon {
		lx.Emit(TokenIdent)
		return scanState
	}
	lx.Read() // skip colon

	var (
		pos = lx.pos
		b   = lx.b[lx.pos:]
	)
	if r, err := lx.Peek(); err != nil || !unicode.IsLetter(r) {
		err := fmt.Errorf("missing identifier after colon at %d - %q", pos, b)
		return errorState(err)
	}
	lx.ScanFunc(unicode.IsLetter)

	if r, err := lx.Peek(); err != nil || !unicode.IsNumber(r) {
		err := fmt.Errorf("invalid identifier after colon at %d - %q", pos, b)
		return errorState(err)
	}
	lx.ScanFunc(unicode.IsNumber)

	lx.Emit(TokenRange)
	return scanState
}

func genState(tok tokenType) stateFunc {
	return func(lx *lexer) stateFunc {
		lx.Read()
		lx.Emit(tok)
		return scanState
	}
}

// errorState emits an error token and advances to eofState.
func errorState(err error) stateFunc {
	return func(lx *lexer) stateFunc {
		lx.Err = err
		lx.Emit(TokenError)
		return eofState
	}
}

// eofState terminates the sequence of the states.
var eofState stateFunc = nil
