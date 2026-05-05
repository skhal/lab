// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lex

import "unicode/utf8"

// blockReader is a text block reader. A text block is a slice of string,
// [start:end]. Use [Read] to grow the block, the end-index, then call [Text]
// to retrieve the block and reset the block indices to read the next block.
//
// blockReader also supports retrieval of the next rune without advancing the
// current position, [Peek].
//
// The [Unread] call decreases the block size - it drops the last rune in the
// block. It become noop if current block's size becomes zero.
//
// EXAMPLE
//
//	r := newBufReader("abc def")
//	r.Read()   // ('a', true) start:0 end:1
//	r.Read()   // ('b', true) start:0 end:2
//	r.Read()   // ('c', true) start:0 end:3
//	r.Peek()   // (' ', true) start:0 end:3
//	r.Unread() // start: 0 end:2
//	r.Text()   // ("ab", 0)
type blockReader struct {
	pos   *Positioner
	str   string
	block block
}

type block struct {
	start int
	end   int
}

// Empty returns true if the block is empty, e.g. the start and end indices are
// equal.
func (b *block) Empty() bool {
	return b.start == b.end
}

func newBlockReader(s string, pos *Positioner) *blockReader {
	return &blockReader{str: s, pos: pos}
}

// Ignore resets current block. It sets block's start position to the end. This
// effectively skips the block and make the reader ready to read the next
// block.
func (br *blockReader) Ignore() {
	br.block.start = br.block.end
}

// Peek returns the next rune but does not advance the read-position. It
// returns a second parameter to indicate whether the read was successful.
func (br *blockReader) Peek() (rune, bool) {
	r, sz := utf8.DecodeRuneInString(br.str[br.block.end:])
	if sz == 0 {
		return 0, false
	}
	if sz == 1 && r == utf8.RuneError {
		return 0, false
	}
	return r, true
}

// Pos returns current read-position.
func (br *blockReader) Pos() Position {
	return br.pos.PosIndex(br.block.end)
}

// Read returns the next rune and advances the read-position. It returns a
// second parameter to indicate whether the read was successful.
func (br *blockReader) Read() (rune, bool) {
	r, sz := utf8.DecodeRuneInString(br.str[br.block.end:])
	if sz == 0 {
		return 0, false
	}
	if sz == 1 && r == utf8.RuneError {
		return 0, false
	}
	if r == runeEOL {
		br.pos.push(br.block.end)
	}
	br.block.end += sz
	return r, true
}

// Text returns current block and its position. After the call, the reader is
// ready to read the next block.
func (br *blockReader) Text() (s string, pos int) {
	defer br.Ignore()
	return br.str[br.block.start:br.block.end], br.block.start
}

// Unread reduces current block by one rune. It is a noop operation if the
// block is empty.
func (br *blockReader) Unread() {
	if br.block.Empty() {
		return
	}
	r, sz := utf8.DecodeLastRuneInString(br.str[:br.block.end])
	if sz == 0 {
		return
	}
	if sz == 1 && r == utf8.RuneError {
		return
	}
	if r == runeEOL {
		br.pos.pop()
	}
	br.block.end -= sz
}
