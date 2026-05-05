// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lex

import "unicode/utf8"

// bufReader is a buffered reader. It supports scan a string by runes - read,
// unread, peek next - and get a block of at once.
//
// EXAMPLE
//
//	r := newBufReader("abc def")
//	r.Read() // ('a', true)
//	r.Read() // ('b', true)
//	r.Read() // ('c', true)
//	r.Text() // ("abc", 0, 3)
type bufReader struct {
	pos *Positioner
	str string

	// state
	start int // token index
	end   int // current position
}

func newBufReader(s string, pos *Positioner) *bufReader {
	return &bufReader{str: s, pos: pos}
}

// Ignore resets current block.
func (br *bufReader) Ignore() {
	br.start = br.end
}

// Peek returns the next next without advancing the position of last rune in the
// block. It returns a boolean flag to indicate whether the read was successful.
func (br *bufReader) Peek() (rune, bool) {
	r, sz := utf8.DecodeRuneInString(br.str[br.end:])
	if sz == 0 {
		return 0, false
	}
	return r, true
}

// Pos returns current position in the read buffer.
func (br *bufReader) Pos() Position {
	return br.pos.PosIndex(br.end)
}

// Read returns the next next and advances the position of last rune in the
// block. It returns a boolean flag to indicate whether the read was successful.
func (br *bufReader) Read() (rune, bool) {
	r, sz := utf8.DecodeRuneInString(br.str[br.end:])
	if sz == 0 {
		return 0, false
	}
	if r == runeEOL {
		br.pos.push(br.end)
	}
	br.end += sz
	return r, true
}

// Text returns the text of the current block and resets the block, ready to
// read the next one.
func (br *bufReader) Text() (s string, pos int) {
	defer br.Ignore()
	return br.str[br.start:br.end], br.start
}

// Unread unreads the last rune in the block. It is a noop operation if the
// block is empty.
func (br *bufReader) Unread() {
	if br.start == br.end {
		// can't unread beyond the beginning of the token
		return
	}
	r, sz := utf8.DecodeLastRuneInString(br.str[:br.end])
	if sz == 0 {
		return
	}
	if sz == 1 && r == utf8.RuneError {
		return
	}
	if r == runeEOL {
		br.pos.pop()
	}
	br.end -= sz
}
