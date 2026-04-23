// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ErrParse means the was en error parsing data.
var ErrParse = errors.New("parse error")

// Parse runs lexical and syntacit analysis of s. It returns an AST node upon
// success or a non-nil error in case of failure.
func Parse(s string) (Node, error) {
	const (
		formulaPrefix = "="
	)
	parse := parseFormula
	s, ok := strings.CutPrefix(s, formulaPrefix)
	if !ok {
		parse = parseCell
	}
	n, err := parse(strings.TrimSpace(s))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrParse, err)
	}
	return n, nil
}

var (
	cellRx = regexp.MustCompile(`^\d+(?:\.\d*)?$`)   // numbers only
	refRx  = regexp.MustCompile(`^[[:upper:]]+\d+$`) // ref: ABC123
	callRx = regexp.MustCompile(`^[[:upper:]]+$`)    // f-n name: ABC
)

// parseCell parses a cell without formula.
func parseCell(s string) (Node, error) {
	if len(s) == 0 {
		return nil, fmt.Errorf("empty cell")
	}
	if !cellRx.MatchString(s) {
		return nil, fmt.Errorf("not a number %q", s)
	}
	return &NumberNode{Number: s}, nil
}

func parseFormula(s string) (Node, error) {
	if len(s) == 0 {
		return nil, fmt.Errorf("empty formula")
	}
	p := &formulaParser{}
	return p.Parse([]byte(s))
}
