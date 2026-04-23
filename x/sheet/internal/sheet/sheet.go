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

	"github.com/skhal/lab/x/sheet/internal/engine"
)

// ErrCell means there is an error in the cell value.
var ErrCell = errors.New("cell error")

// Sheet is a table of cells, organized into columns and rows.
type Sheet struct {
	data map[string]*cell
	eng  Engine
}

// Engine is the Sheets table backend to parse and calculate cells.
type Engine interface {
	// Parse cell content
	Parse(string) (any, error)

	// Calculate calculate cell value. The passed function can be used to
	// trigger calculation of the reference cells.
	Calculate(any, func(string) (float64, error)) (float64, error)
}

type cell struct {
	ID   string
	Text string

	ir any

	calculated bool
	result     float64
}

// Option represents a configuration option for the Sheets table.
type Option func(*Sheet)

// WithASTEngine makes Sheets table use AST engine.
func WithASTEngine() Option {
	return func(s *Sheet) {
		if s.eng != nil {
			panic("engine is already set")
		}
		s.eng = engine.NewAST()
	}
}

// New creates a cells table.
func New(opts ...Option) *Sheet {
	s := &Sheet{
		data: make(map[string]*cell),
	}
	for _, opt := range opts {
		opt(s)
	}
	if s.eng == nil {
		opt := WithASTEngine()
		opt(s)
	}
	return s
}

// Set places a value to the cell.
func (s *Sheet) Set(id, val string) error {
	data, err := s.eng.Parse(val)
	if err != nil {
		return fmt.Errorf("%w: set %s to %q: %s", ErrCell, id, val, err)
	}
	s.data[id] = &cell{ID: id, Text: val, ir: data}
	return nil
}

// Calculate parses cell content. It returns an error if any of the cell fails
// to parse.
func (s *Sheet) Calculate() error {
	rc := newCalculator(s)
	return rc.Calculate()
}

// VisitAll calls function f on every cell. It passes cell identifier and
// calculated value.
func (s *Sheet) VisitAll(f func(id, cell string, val float64) bool) {
	kk := slices.Collect(maps.Keys(s.data))
	slices.Sort(kk)
	for _, id := range kk {
		c := s.data[id]
		if !f(id, c.Text, c.result) {
			break
		}
	}
}

// Write writes the sheet to the writer in binary format. It returns an error
// if it fails to write data.
func (s *Sheet) Write(w io.Writer) error {
	enc := gob.NewEncoder(w)
	if err := enc.Encode(s.eng); err != nil {
		return fmt.Errorf("write: failed to save engine: %s", err)
	}
	kk := slices.Collect(maps.Keys(s.data))
	slices.Sort(kk)
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
	dec := gob.NewDecoder(r)
	if err := dec.Decode(s.eng); err != nil {
		return fmt.Errorf("read: failed to load engine: %s", err)
	}
	if len(s.data) != 0 {
		s.data = New().data
	}
	var c cell
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
