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

var whitespace = []byte(` \t`)

// scanState is the default state of the scanner. It skips whitespace and
// advances to the next supported state.
func scanState(lx *lexer) stateFunc {
	const (
		plus  = '+'
		minus = '-'
		lpar  = '('
		rpar  = ')'
	)
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
	case unicode.IsLetter(r):
		return identifierState
	case r == plus:
		return genState(TokenPlus)
	case r == minus:
		return genState(TokenMinus)
	case r == lpar:
		return genState(TokenLpar)
	case r == rpar:
		return genState(TokenRpar)
	default:
		lx.err = fmt.Errorf("unsupported text at %d -  %q", lx.pos, lx.b[lx.pos:])
		return errorState
	}
}

// numberState parses a floating value number with non-empty integral part. It
// emits parsed number token and advances to the scanState.
func numberState(lx *lexer) stateFunc {
	lx.scanFunc(unicode.IsDigit)
	switch r, err := lx.peek(); {
	case err != nil:
		// read failed - the next state will handle the error
	case r == '.':
		lx.read()
		lx.scanFunc(unicode.IsDigit)
	}
	lx.emit(TokenNumber)
	return scanState
}

// identifierState parses an identifier that consists of any number of letters
// optionally followed by numbers, e.g. "abC123" or "Ab". The function is case
// insensitive.
func identifierState(lx *lexer) stateFunc {
	lx.scanFunc(unicode.IsLetter)
	switch r, err := lx.peek(); {
	case err != nil:
	case unicode.IsNumber(r):
		lx.scanFunc(unicode.IsNumber)
	}
	lx.emit(TokenIdent)
	return scanState
}

func genState(tok tokenType) stateFunc {
	return func(lx *lexer) stateFunc {
		lx.read()
		lx.emit(tok)
		return scanState
	}
}

// errorState emits an error token and advances to eofState.
func errorState(lx *lexer) stateFunc {
	lx.emit(TokenError)
	return eofState
}

// eofState terminates the sequence of the states.
var eofState stateFunc = nil
