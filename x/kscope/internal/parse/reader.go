// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parse

import (
	"github.com/skhal/lab/x/kscope/internal/lex"
)

type reader interface {
	// Read returns the next token and a flag to indicate whether a token exists.
	Read() (lex.Token, bool)
}

// peekerReader is a tokens reader that can read or peek the next token.
type peekerReader struct {
	reader
	// next is the next token, cached by the last [Peek] call, and reset by
	// the next call to [Read].
	next *struct {
		tok lex.Token
		ok  bool
	}
}

// Peek returns the next token after last [Read] and a flag to indicate whether
// a token exists. Repeated calls return the same next token.
func (r *peekerReader) Peek() (lex.Token, bool) {
	if r.next == nil {
		tok, ok := r.reader.Read()
		r.next = &struct {
			tok lex.Token
			ok  bool
		}{tok, ok}
	}
	return r.next.tok, r.next.ok
}

// Read returns the next token and a flag to indicate whether a token exists.
func (r *peekerReader) Read() (lex.Token, bool) {
	if r.next != nil {
		cache := r.next
		r.next = nil
		return cache.tok, cache.ok
	}
	return r.reader.Read()
}
