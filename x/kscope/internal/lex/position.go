// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lex

import (
	"fmt"
	"slices"
)

// Position is the line and column position of a token in the code.
type Position struct {
	Line int // line number (count from 1)
	Col  int // column number (count from 1)
}

// String prints the position as a colon separated pair.
func (l Position) String() string {
	return fmt.Sprintf("%d:%d", l.Line, l.Col)
}

// Positioner returns the position, the line and column number, within the
// source code.
type Positioner struct {
	eols []int // sorted indises of new lines in the code.
	pos  int   // position of the last read character
}

// Pos returns current read position.
func (p *Positioner) Pos() Position {
	return p.PosIndex(p.pos)
}

// PosToken converts token index in the code to the line and column number. It
// assumes that the code is ASCII in calculating the column number, i.e. every
// character is a single byte rune.
func (p *Positioner) PosToken(tok Token) Position {
	return p.PosIndex(tok.pos)
}

// PosIndex converts code index to the line and column number. See [Pos] for
// encoding assumptions.
func (p *Positioner) PosIndex(n int) Position {
	if n < 0 {
		return Position{}
	}
	return p.posIndex(n)
}

func (p *Positioner) posIndex(n int) Position {
	ln := p.findLine(n)
	col := p.findCol(n, ln)
	return Position{Line: ln + 1, Col: col}
}
func (p *Positioner) findLine(pos int) int {
	// oftentimes the code wants the position of the last token, e.g. errors.
	if len(p.eols) == 0 || pos > p.eols[len(p.eols)-1] {
		return len(p.eols)
	}
	ln, _ := slices.BinarySearch(p.eols, pos)
	return ln
}

func (p *Positioner) findCol(pos int, ln int) int {
	if ln == 0 {
		return pos + 1
	}
	return pos - p.eols[ln-1] // -1 for previous line
}

// push adds an end-of-line index.
func (p *Positioner) push(eol int) {
	p.eols = append(p.eols, eol)
}

// pop removes last end-of-line index.
func (p *Positioner) pop() {
	p.eols = p.eols[:len(p.eols)-1]
}

// position stores current position.
func (p *Positioner) position(n int) {
	p.pos = n
}
