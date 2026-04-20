// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lex

import (
	"fmt"
	"io"
	"unicode"
)

type stateFunc func(*lexer) stateFunc

var (
	digits     = []byte(`0123456789`)
	whitespace = []byte(` \t`)
)

// scanState is the default state of the scanner. It skips whitespace and
// advances to the next supported state.
func scanState(lx *lexer) stateFunc {
	lx.scan(whitespace)
	lx.ignore()
	switch r, err := lx.peek(); {
	case err == io.EOF:
		return eofState
	case err != nil:
		lx.err = err
		return errorState
	case unicode.IsNumber(r):
		return numberState
	case r == '+':
		return plusState
	case r == '-':
		return minusState
	default:
		lx.err = fmt.Errorf("unsupported text at %d -  %q", lx.pos, lx.b[lx.pos:])
		return errorState
	}
}

// numberState parses a floating value number with non-empty integral part. It
// emits parsed number token and advances to the scanState.
func numberState(lx *lexer) stateFunc {
	lx.scan(digits)
	switch r, err := lx.peek(); {
	case err != nil:
		// read failed - the next state will handle the error
	case r == '.':
		lx.read()
		lx.scan(digits)
	}
	lx.emit(TokenNumber)
	return scanState
}

func plusState(lx *lexer) stateFunc {
	lx.read() // ignore err - previous state has peeked into the next rune
	lx.emit(TokenPlus)
	return scanState
}

func minusState(lx *lexer) stateFunc {
	lx.read() // ignore err - previous state has peeked into the next rune
	lx.emit(TokenMinus)
	return scanState
}

// errorState emits an error token and advances to eofState.
func errorState(lx *lexer) stateFunc {
	lx.emit(TokenError)
	return eofState
}

// eofState terminates the sequence of the states.
var eofState stateFunc = nil
