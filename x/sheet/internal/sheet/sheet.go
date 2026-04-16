// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sheet implements a cells table.
package sheet

import (
	"errors"
	"fmt"
	"maps"
	"slices"

	"github.com/skhal/lab/x/sheet/internal/ast"
	"github.com/skhal/lab/x/sheet/internal/calc"
	"github.com/skhal/lab/x/sheet/internal/parse"
)

// ErrCell means there is an error in the cell value.
var ErrCell = errors.New("cell error")

type sheet struct {
	data map[string]*cell
}

// New creates a cells table.
func New() *sheet {
	return &sheet{make(map[string]*cell)}
}

// Set places a value to the cell.
func (s *sheet) Set(id, val string) error {
	n, err := parse.Parse(val)
	if err != nil {
		return fmt.Errorf("%w: set %s to %q: %s", ErrCell, id, val, err)
	}
	s.data[id] = &cell{Text: val, Node: n}
	return nil
}

// Calculate parses cell content. It returns an error if any of the cell fails
// to parse.
func (s *sheet) Calculate() error {
	for id, c := range s.data {
		if err := calc.Calculate(c.Node); err != nil {
			return fmt.Errorf("calculate %s: %s", id, err)
		}
	}
	return nil
}

// VisitAll calls function f on every cell. It passes cell identifier and
// calculated value.
func (s *sheet) VisitAll(f func(cell string, val float64) bool) {
	kk := slices.Collect(maps.Keys(s.data))
	slices.Sort(kk)
	for _, id := range kk {
		c := s.data[id]
		if !f(id, c.Node.Value()) {
			break
		}
	}
}

type cell struct {
	Text string
	Node ast.Node
}
