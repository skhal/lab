// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lex

// scanner scans for tokens using blockReader.
type scanner struct {
	r     *blockReader
	state scanFunc // current state, can be nil.
	tok   *Token   // last token
	err   error
}

// newScanner creates a token scanner with provided blockReader.
func newScanner(r *blockReader) *scanner {
	return &scanner{r: r, state: scan}
}

// Scan extracts the next token and returns true if the extraction was
// successful. It is noop if the last call to Scan resulted in an error or
// the internal scan state is EOF, represented by nil state.
func (sc *scanner) Scan() bool {
	if sc.err != nil || sc.state == nil {
		return false
	}
	sc.tok, sc.state, sc.err = sc.state(sc.r)
	return sc.tok != nil
}

// Token returns last extracted token.
func (sc *scanner) Token() *Token {
	return sc.tok
}

// Err returns last extraction error if any.
func (sc *scanner) Err() error {
	return sc.err
}
