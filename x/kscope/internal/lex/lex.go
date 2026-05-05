// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lex

import (
	"iter"
)

// Lexer converts a string into a sequence of tokens. A zero value Lexer is
// ready to use.
type Lexer struct {
	err error
	pos *Positioner
}

// Lex runs lexical tokenization on a string s. It returns a sequence of tokens
// and a [Positioner] to convert token positions to line and column.
func (lx *Lexer) Lex(s string) (iter.Seq[Token], *Positioner) {
	lx.reset()
	return lx.lex(s), lx.pos
}

// reset resets the Lexer state.
func (lx *Lexer) reset() {
	lx.err = nil
	lx.pos = new(Positioner)
}

// lex converts the string into a sequence of tokens.
func (lx *Lexer) lex(s string) iter.Seq[Token] {
	return func(yield func(Token) bool) {
		sc := newScanner(newBlockReader(s, lx.pos))
		for sc.Scan() {
			if !yield(*sc.Token()) {
				break
			}
		}
		lx.err = sc.Err()
	}
}

// Err returns the last encountered error if any.
func (lx *Lexer) Err() error {
	return lx.err
}
