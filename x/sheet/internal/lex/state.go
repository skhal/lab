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

const (
	// keep-sorted start
	runeComma       = ','
	runeDivide      = '/'
	runeEqual       = '='
	runeExclamation = '!'
	runeGreater     = '>'
	runeLess        = '<'
	runeLpar        = '('
	runeMinus       = '-'
	runeMultiply    = '*'
	runePlus        = '+'
	runeRpar        = ')'
	// keep-sorted end
)

func isWhitespace(r rune) bool { return bytes.ContainsRune(whitespace, r) }

// scanState is the default state of the scanner. It skips whitespace and
// advances to the next supported state.
func scanState(lx *lexer) stateFunc {
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
	case r == runeComma:
		return genState(TokenComma)
	case r == runeDivide:
		return genState(TokenDivide)
	case r == runeEqual:
		return equalState
	case r == runeExclamation:
		return notEqualState
	case r == runeGreater:
		return greaterState
	case r == runeLess:
		return lessState
	case r == runeLpar:
		return genState(TokenLpar)
	case r == runeMinus:
		return genState(TokenMinus)
	case r == runeMultiply:
		return genState(TokenMultiply)
	case r == runePlus:
		return genState(TokenPlus)
	case r == runeRpar:
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

func genState(tok TokenType) stateFunc {
	return func(lx *lexer) stateFunc {
		lx.Read()
		lx.Emit(tok)
		return scanState
	}
}

func equalState(lx *lexer) stateFunc {
	lx.Read() // ignore the equal sign
	switch r, err := lx.Read(); {
	case err != nil:
		return errorState(err)
	case r != runeEqual:
		err := fmt.Errorf("unexpected '=%v', want ==", r)
		return errorState(err)
	}
	lx.Emit(TokenEqual)
	return scanState
}

func notEqualState(lx *lexer) stateFunc {
	lx.Read() // ignore the exclamation mark
	switch r, err := lx.Read(); {
	case err != nil:
		return errorState(err)
	case r != runeEqual:
		err := fmt.Errorf("unexpected '!%v', want !=", r)
		return errorState(err)
	}
	lx.Emit(TokenNotEqual)
	return scanState
}

func lessState(lx *lexer) stateFunc {
	lx.Read() // ignore the less sign
	if r, err := lx.Peek(); err != nil || r != runeEqual {
		lx.Emit(TokenLess)
		return scanState
	}
	lx.Read() // ignore the equal sign
	lx.Emit(TokenLessOrEqual)
	return scanState
}

func greaterState(lx *lexer) stateFunc {
	lx.Read() // ignore the greater sign
	if r, err := lx.Peek(); err != nil || r != runeEqual {
		lx.Emit(TokenGreater)
		return scanState
	}
	lx.Read() // ignore the equal sign
	lx.Emit(TokenGreaterOrEqual)
	return scanState
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
