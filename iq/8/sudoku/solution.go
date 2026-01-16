// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// Board:
//
//      A B C   D E F   G H I
//  1 | . . . | . . . | . . . |
//  2 | . . . | . . . | . . . |
//  3 | . . . | . . . | . . . |
//      -----   -----   -----
//  4 | . . . | . . . | . . . |
//  5 | . . . | . . . | . . . |
//  6 | . . . | . . . | . . . |
//      -----   -----   -----
//  7 | . . . | . . . | . . . |
//  8 | . . . | . . . | . . . |
//  9 | . . . | . . . | . . . |
//      -----   -----   -----

package sudoku

import (
	"bytes"
	"fmt"
	"strconv"
)

type ColID int

const (
	ColA ColID = iota
	ColB
	ColC
	ColD
	ColE
	ColF
	ColG
	ColH
	ColI
	NumCols
)

func newColID(boxid BoxID) ColID {
	return ColID(boxid / 9)
}

func (c ColID) String() string {
	const offset = int('A')
	return string(byte(int(c) + offset))
}

type RowID int

const (
	Row1 RowID = iota
	Row2
	Row3
	Row4
	Row5
	Row6
	Row7
	Row8
	Row9
	NumRows
)

func newRowID(boxid BoxID) RowID {
	return RowID(boxid % 9)
}

func (r RowID) String() string {
	const offset = 1
	return strconv.Itoa(int(r) + offset)
}

const boxEmpty = 0

type Row [NumCols]int

func (r *Row) String() string {
	const (
		// keep-sorted start
		blockSeparator = " | "
		boxSeparator   = ' '
		emptyChar      = '.'
		emptyRow       = ". . . | . . . | . . ."
		// keep-sorted end
	)
	if r == nil {
		return emptyRow
	}
	buf := new(bytes.Buffer)
	for c := ColA; c < NumCols; c++ {
		n := r[c]
		switch n {
		case boxEmpty:
			buf.WriteByte(emptyChar)
		default:
			buf.WriteString(strconv.Itoa(n))
		}
		switch c {
		case ColC, ColF:
			buf.WriteString(blockSeparator)
		case ColI:
		default:
			buf.WriteByte(boxSeparator)
		}
	}
	return buf.String()
}

type BoxID int

const (
	A1 BoxID = iota
	A2
	A3
	A4
	A5
	A6
	A7
	A8
	A9
	B1
	B2
	B3
	B4
	B5
	B6
	B7
	B8
	B9
	C1
	C2
	C3
	C4
	C5
	C6
	C7
	C8
	C9
	D1
	D2
	D3
	D4
	D5
	D6
	D7
	D8
	D9
	E1
	E2
	E3
	E4
	E5
	E6
	E7
	E8
	E9
	F1
	F2
	F3
	F4
	F5
	F6
	F7
	F8
	F9
	G1
	G2
	G3
	G4
	G5
	G6
	G7
	G8
	G9
	H1
	H2
	H3
	H4
	H5
	H6
	H7
	H8
	H9
	I1
	I2
	I3
	I4
	I5
	I6
	I7
	I8
	I9
)

type Box struct {
	ID  BoxID
	Num int
}

type Board struct {
	rows [NumRows]*Row
}

func NewBoard(bb ...Box) *Board {
	board := &Board{}
	for _, b := range bb {
		rowid := newRowID(b.ID)
		row := board.rows[rowid]
		if row == nil {
			row = &Row{}
			board.rows[rowid] = row
		}
		colid := newColID(b.ID)
		row[colid] = b.Num
	}
	return board
}

func (b *Board) String() string {
	const blockSpaceOffset = "    "
	writeHeader := func(b *bytes.Buffer) {
		const blockSpaceSeparator = "   "
		const (
			eol   = '\n'
			space = ' '
		)
		b.WriteString(blockSpaceOffset)
		for c := ColA; c < NumCols; c++ {
			b.WriteString(c.String())
			switch c {
			case ColC, ColF:
				b.WriteString(blockSpaceSeparator)
			case ColI: // nothing
			default:
				b.WriteByte(space)
			}
		}
		b.WriteByte(eol)
	}
	writeSeparator := func(b *bytes.Buffer) {
		const separator = "-----   -----   -----"
		fmt.Fprintf(b, "%s%s\n", blockSpaceOffset, separator)
	}
	buf := new(bytes.Buffer)
	writeHeader(buf)
	for r := Row1; r < NumRows; r++ {
		row := b.rows[r]
		fmt.Fprintf(buf, "%v | %v |\n", r, row)
		switch r {
		case Row3, Row6, Row9:
			writeSeparator(buf)
		}
	}
	return buf.String()
}

func (b *Board) IsValid() bool {
	v := newValidator()
	return v.validate(b)
}

type validator struct {
	boxes map[RowID]map[ColID]int
}

func newValidator() *validator {
	return &validator{
		boxes: make(map[RowID]map[ColID]int),
	}
}

func (v *validator) validate(b *Board) bool {
	if ok := v.validateRows(b); !ok {
		return false
	}
	if ok := v.validateColumns(); !ok {
		return false
	}
	return v.validateBlocks()
}

func (v *validator) validateRows(b *Board) bool {
	for r := Row1; r < NumRows; r++ {
		row := b.rows[r]
		if row == nil {
			continue
		}
		seen := make(map[int]struct{})
		for c, n := range row {
			switch n {
			case boxEmpty: // do nothing
			default:
				if _, ok := seen[n]; ok {
					return false
				}
				seen[n] = struct{}{}
				v.fillBox(r, ColID(c), n)
			}
		}
	}
	return true
}

func (v *validator) fillBox(r RowID, c ColID, n int) {
	row, ok := v.boxes[r]
	if !ok {
		row = make(map[ColID]int)
		v.boxes[r] = row
	}
	row[c] = n
}

func (v *validator) validateColumns() bool {
	for c := ColA; c < NumCols; c++ {
		if !v.isValidColumn(c) {
			return false
		}
	}
	return true
}

func (v *validator) isValidColumn(c ColID) bool {
	seen := make(map[int]struct{})
	for r := Row1; r < NumRows; r++ {
		row, ok := v.boxes[r]
		if !ok {
			continue
		}
		n, ok := row[c]
		if !ok {
			continue
		}
		if _, ok := seen[n]; ok {
			return false
		}
		seen[n] = struct{}{}
	}
	return true
}

func (v *validator) validateBlocks() bool {
	blocks := []struct {
		start BoxID
		end   BoxID
	}{
		{start: A1, end: C3},
		{start: D1, end: F3},
		{start: G1, end: I3},
		{start: A4, end: C6},
		{start: D4, end: F6},
		{start: G4, end: I6},
		{start: A7, end: C9},
		{start: D7, end: F9},
		{start: G7, end: I9},
	}
	for _, b := range blocks {
		if !v.isValidBlock(b.start, b.end) {
			return false
		}
	}
	return true
}

func (v *validator) isValidBlock(start, end BoxID) bool {
	seen := make(map[int]struct{})
	for r, rend := newRowID(start), newRowID(end); r <= rend; r++ {
		for c, cend := newColID(start), newColID(end); c <= cend; c++ {
			row, ok := v.boxes[r]
			if !ok {
				continue
			}
			n, ok := row[c]
			if !ok {
				continue
			}
			if _, ok := seen[n]; ok {
				return false
			}
			seen[n] = struct{}{}
		}
	}
	return true
}
