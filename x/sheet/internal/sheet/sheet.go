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
	s.data[id] = &cell{Text: val, AST: n}
	return nil
}

// Calculate parses cell content. It returns an error if any of the cell fails
// to parse.
func (s *sheet) Calculate() error {
	rc := newRefCalculator(s)
	for id, c := range s.data {
		if err := s.calculate(id, c, rc); err != nil {
			return fmt.Errorf("%w %s: calculate: %s", ErrCell, id, err)
		}
	}
	return nil
}

func (s *sheet) calculate(id string, c *cell, rc *refCalculator) error {
	if c.Calculated {
		return nil
	}
	res, err := calc.Calculate(c.AST, rc)
	if err != nil {
		return fmt.Errorf("%s: %s", id, err)
	}
	c.Value = res
	c.Calculated = true
	return nil
}

type refCalculator struct {
	s    *sheet
	seen map[string]bool
	path []string // references path for cycle detection
}

func newRefCalculator(s *sheet) *refCalculator {
	return &refCalculator{s: s, seen: make(map[string]bool)}
}

// Calculate calculates a reference value. It returns an error if calculation
// fails or reference calculator detects a circular dependency.
func (rc *refCalculator) Calculate(id string) (float64, error) {
	if rc.seen[id] {
		// circular dependency
		return 0, fmt.Errorf("circular dependency - %v", rc.path)
	}
	rc.seen[id] = true
	rc.path = append(rc.path, id)
	defer func() {
		rc.path = rc.path[:len(rc.path)-1]
		rc.seen[id] = false
	}()
	c, ok := rc.s.data[id]
	if !ok {
		return 0, fmt.Errorf("invalid cell %s", id)
	}
	if err := rc.s.calculate(id, c, rc); err != nil {
		return 0, err
	}
	return c.Value, nil
}

// VisitAll calls function f on every cell. It passes cell identifier and
// calculated value.
func (s *sheet) VisitAll(f func(id, cell string, val float64) bool) {
	kk := slices.Collect(maps.Keys(s.data))
	slices.Sort(kk)
	for _, id := range kk {
		c := s.data[id]
		if !f(id, c.Text, c.Value) {
			break
		}
	}
}

type cell struct {
	Text string
	AST  ast.Node

	Calculated bool
	Value      float64
}
