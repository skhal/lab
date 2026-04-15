// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sheet implements a cells table.
package sheet

import (
	"fmt"

	"github.com/skhal/lab/x/sheet/internal/ast"
	parser "github.com/skhal/lab/x/sheet/internal/parse"
)

type sheet struct {
	data map[string]*cell
}

// New creates a cells table.
func New() *sheet {
	return &sheet{make(map[string]*cell)}
}

// Set places a value to the cell.
func (s *sheet) Set(id, val string) {
	s.data[id] = &cell{Text: val}
}

// Calculate parses cell content. It returns an error if any of the cell fails
// to parse.
func (s *sheet) Calculate() error {
	for id, c := range s.data {
		if err := s.parse(c); err != nil {
			return fmt.Errorf("calculate %s: %s", id, err)
		}
	}
	return nil
}

func (s *sheet) parse(c *cell) error {
	n, err := parser.Parse(c.Text)
	if err != nil {
		return err
	}
	c.Node = n
	return nil
}

// VisitAll calls function f on every cell. It passes cell identifier and
// calculated value.
func (s *sheet) VisitAll(f func(cell string, val float64) bool) {
	for id, c := range s.data {
		if !f(id, c.Node.Value()) {
			break
		}
	}
}

type cell struct {
	Text string
	Node ast.Node
}
