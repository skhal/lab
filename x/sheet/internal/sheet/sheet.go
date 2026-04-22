// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sheet implements a cells table.
package sheet

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"

	"github.com/skhal/lab/x/sheet/internal/ast"
	"github.com/skhal/lab/x/sheet/internal/calc"
	"github.com/skhal/lab/x/sheet/internal/parse"
)

// ErrCell means there is an error in the cell value.
var ErrCell = errors.New("cell error")

// Sheet is a table of cells, organized into columns and rows.
type Sheet struct {
	data map[string]*cell
}

type cell struct {
	ID   string
	Text string

	ast ast.Node

	calculated bool
	value      float64
}

// New creates a cells table.
func New() *Sheet {
	return &Sheet{make(map[string]*cell)}
}

// Set places a value to the cell.
func (s *Sheet) Set(id, val string) error {
	n, err := parse.Parse(val)
	if err != nil {
		return fmt.Errorf("%w: set %s to %q: %s", ErrCell, id, val, err)
	}
	s.data[id] = &cell{ID: id, Text: val, ast: n}
	return nil
}

// Calculate parses cell content. It returns an error if any of the cell fails
// to parse.
func (s *Sheet) Calculate() error {
	rc := newRefCalculator(s)
	for id, c := range s.data {
		if err := s.calculate(id, c, rc); err != nil {
			return fmt.Errorf("%w %s: calculate: %s", ErrCell, id, err)
		}
	}
	return nil
}

func (s *Sheet) calculate(id string, c *cell, rc *refCalculator) error {
	if c.calculated {
		return nil
	}
	res, err := calc.Calculate(c.ast, rc)
	if err != nil {
		return fmt.Errorf("%s: %s", id, err)
	}
	c.value = res
	c.calculated = true
	return nil
}

type refCalculator struct {
	s    *Sheet
	seen map[string]bool
	path []string // references path for cycle detection
}

func newRefCalculator(s *Sheet) *refCalculator {
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
		return 0, nil
	}
	if err := rc.s.calculate(id, c, rc); err != nil {
		return 0, err
	}
	return c.value, nil
}

// VisitAll calls function f on every cell. It passes cell identifier and
// calculated value.
func (s *Sheet) VisitAll(f func(id, cell string, val float64) bool) {
	kk := slices.Collect(maps.Keys(s.data))
	slices.Sort(kk)
	for _, id := range kk {
		c := s.data[id]
		if !f(id, c.Text, c.value) {
			break
		}
	}
}

// Write writes the sheet to the writer in binary format. It returns an error
// if it fails to write data.
func (s *Sheet) Write(w io.Writer) error {
	kk := slices.Collect(maps.Keys(s.data))
	slices.Sort(kk)
	enc := gob.NewEncoder(w)
	for id := range slices.Values(kk) {
		c := s.data[id]
		if err := enc.Encode(c); err != nil {
			return fmt.Errorf("write: %s", err)
		}
	}
	return nil
}

// Read reads the sheet from the reader. It resets the sheet if there is any
// data in cells.
func (s *Sheet) Read(r io.Reader) error {
	if len(s.data) != 0 {
		s.data = New().data
	}
	var (
		c   cell
		dec = gob.NewDecoder(r)
	)
	for {
		err := dec.Decode(&c)
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return fmt.Errorf("read: %s", err)
		}
		s.Set(c.ID, c.Text)
	}
}
