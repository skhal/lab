// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lex

type scanner struct {
	r     *bufReader
	state scanFunc
	tok   *Token
	err   error
}

func newScanner(r *bufReader) *scanner {
	return &scanner{r: r, state: scan}
}

// Scan scans for the next token and returns true upon finding one, otherwise
// false.
func (sc *scanner) Scan() bool {
	if sc.err != nil || sc.state == nil {
		return false
	}
	sc.tok, sc.state, sc.err = sc.state(sc.r)
	return sc.tok != nil
}

// Token returns last parsed token.
func (sc *scanner) Token() *Token {
	return sc.tok
}

// Err returns last parse error if any.
func (sc *scanner) Err() error {
	return sc.err
}
