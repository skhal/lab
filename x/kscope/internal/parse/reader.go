// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parse

import (
	"github.com/skhal/lab/x/kscope/internal/lex"
)

// reader returns the next token.
type reader interface {
	Read() (lex.Token, bool) // return the next token and ok boolean flag.
}

// readerFunc adopts a function to the [reader] interface.
type readerFunc func() (lex.Token, bool)

// Read implements [reader] interface.
func (f readerFunc) Read() (lex.Token, bool) {
	return f()
}

// tokenReader is a token reader. It can peek and reed the next token.
type tokenReader struct {
	reader
	// next caches the next token retrieved by the last call to [Peek].
	next *struct {
		tok lex.Token
		ok  bool
	}
}

// Peek retrieves the next token. It returns a second parameter to indicate
// whether the next token exists.
//
// Repeated calls to Peek return the same result.
func (tr *tokenReader) Peek() (lex.Token, bool) {
	if tr.next == nil {
		tr.cacheNext()
	}
	return tr.next.tok, tr.next.ok
}

func (tr *tokenReader) cacheNext() {
	tok, ok := tr.reader.Read()
	tr.next = &struct {
		tok lex.Token
		ok  bool
	}{tok, ok}
}

// Read returns the next token. It uses cached value if the last call was to
// [Peek] and resets the cache so that next calls to Read retrieve the next
// token from the reader.
//
// Read returns a second parameter to indicate whether the next token exists.
func (tr *tokenReader) Read() (lex.Token, bool) {
	if tr.next != nil {
		defer func() { tr.next = nil }()
		return tr.next.tok, tr.next.ok
	}
	return tr.reader.Read()
}
