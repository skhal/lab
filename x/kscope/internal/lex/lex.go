// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lex

import (
	"iter"
)

// Lexer parses a string into a sequence of tokens. A zero value lexer is ready
// to parse a string.
type Lexer struct {
	err error
}

// Lex parses a string s into a sequence of tokens.
func (lx *Lexer) Lex(s string) iter.Seq[Token] {
	lx.err = nil
	return func(yield func(Token) bool) {
		sc := newScanner(newBufReader(s))
		for sc.Scan() {
			if !yield(*sc.Token()) {
				break
			}
		}
		lx.err = sc.Err()
	}
}

// Err last parse error if any.
func (lx *Lexer) Err() error {
	return lx.err
}
