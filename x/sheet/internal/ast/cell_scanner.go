// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

import (
	"cmp"
	"errors"
	"fmt"
	"iter"
	"strconv"
	"strings"
	"unicode"
)

// ErrCellRange means the cell range is invalid.
var ErrCellRange = errors.New("invalid cell range")

// CellScanner scans through cells in the inclusive range [from, to].
//
// It supports inverted range, e.g. "A3:A1". In either case, the scanner
// generates a sequence of cell IDs, sorted by row and column.
type CellScanner struct {
	from cell
	to   cell
}

type cell struct {
	col byte
	row int
}

// NewCellScanner creates a cell scanner. It returns an error if the range
// is invalid, e.g. the cell reference uses double-letter cells (unsupported).
func NewCellScanner(from, to string) (*CellScanner, error) {
	split := func(s string) (c cell, err error) {
		col := strings.TrimRightFunc(s, unicode.IsNumber)
		if len(col) != 1 {
			err = fmt.Errorf("invalid column in %s", s)
			return
		}
		rs := strings.TrimLeftFunc(s, unicode.IsLetter)
		if rs == "" {
			err = fmt.Errorf("invalid row in %s", s)
			return
		}
		row, _ := strconv.Atoi(rs)
		c = cell{col[0], row}
		return
	}
	cellFrom, err := split(from)
	if err != nil {
		return nil, fmt.Errorf("%w: from: %s", ErrCellRange, err)
	}
	cellTo, err := split(to)
	if err != nil {
		return nil, fmt.Errorf("%w: to: %s", ErrCellRange, err)
	}
	cmin, cmax := minmax(cellFrom.col, cellTo.col)
	rmin, rmax := minmax(cellFrom.row, cellTo.row)
	cr := &CellScanner{
		from: cell{cmin, rmin},
		to:   cell{cmax, rmax},
	}
	return cr, nil
}

func minmax[T cmp.Ordered](x, y T) (T, T) {
	return min(x, y), max(x, y)
}

// Len returns the number of cells in the cell range.
func (sc *CellScanner) Len() int {
	// Assuming a single letter columns
	cols := int(sc.to.col-sc.from.col) + 1
	rows := (sc.to.row - sc.from.row + 1)
	return rows * cols
}

// Scan generates a sequence of cell identiifiers in the range.
func (sc *CellScanner) Scan() iter.Seq[string] {
	return func(yield func(string) bool) {
		var id strings.Builder
		for c := sc.from.col; c <= sc.to.col; c++ {
			for r := sc.from.row; r <= sc.to.row; r++ {
				id.Reset()
				id.WriteByte(c)
				id.WriteString(strconv.Itoa(r))
				if !yield(id.String()) {
					return
				}
			}
		}
	}
}
