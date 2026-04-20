// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parse_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/x/sheet/internal/ast"
	"github.com/skhal/lab/x/sheet/internal/parse"
)

func TestParse_nonFormula(t *testing.T) {
	tt := []struct {
		name     string
		s        string
		wantNode ast.Node
		wantErr  error
	}{
		{
			name:    "empty",
			wantErr: parse.ErrParse,
		},
		{
			name:    "not number",
			s:       "abc",
			wantErr: parse.ErrParse,
		},
		{
			name:     "integer",
			s:        "123",
			wantNode: &ast.NumberNode{Number: "123"},
		},
		{
			name:     "float",
			s:        "1.23",
			wantNode: &ast.NumberNode{Number: "1.23"},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parse.Parse(tc.s)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Parse(%q) = _, %v; want %v", tc.s, err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.wantNode, got); diff != "" {
				t.Errorf("Parse(%q) mismatch (-want +got):\n%s", tc.s, diff)
			}
		})
	}
}

func TestParse_formula(t *testing.T) {
	tt := []struct {
		name     string
		s        string
		wantNode ast.Node
		wantErr  error
	}{
		{
			name:    "empty",
			s:       "=",
			wantErr: parse.ErrParse,
		},
		{
			name:     "integer",
			s:        "=123",
			wantNode: &ast.NumberNode{Number: "123"},
		},
		{
			name:     "float",
			s:        "=1.23",
			wantNode: &ast.NumberNode{Number: "1.23"},
		},
		{
			name:     "float no fractional part",
			s:        "=1.",
			wantNode: &ast.NumberNode{Number: "1."},
		},
		{
			name:    "float must have integral part",
			s:       "=.1",
			wantErr: parse.ErrParse,
		},
		{
			name:    "unsupported token",
			s:       "=invalid",
			wantErr: parse.ErrParse,
		},
		{
			name: "operator plus",
			s:    "=1+2",
			wantNode: &ast.BinOpNode{
				Op:    "+",
				Left:  &ast.NumberNode{Number: "1"},
				Right: &ast.NumberNode{Number: "2"},
			},
		},
		{
			name:    "missing operator",
			s:       "=1 2",
			wantErr: parse.ErrParse,
		},
		{
			name:    "plus misses left operand",
			s:       "=+ 2",
			wantErr: parse.ErrParse,
		},
		{
			name:    "plus misses right operand",
			s:       "=1 +",
			wantErr: parse.ErrParse,
		},
		{
			name:    "invalid left operand",
			s:       "=1 + +",
			wantErr: parse.ErrParse,
		},
		{
			name:     "parentheses",
			s:        "=(1)",
			wantNode: &ast.NumberNode{Number: "1"},
		},
		{
			name:    "missing right parenthesis",
			s:       "=(1",
			wantErr: parse.ErrParse,
		},
		{
			name:    "missing left parenthesis",
			s:       "=1)",
			wantErr: parse.ErrParse,
		},
		{
			name:    "parentheses without expression",
			s:       "=()",
			wantErr: parse.ErrParse,
		},
		{
			name:     "nested parentheses",
			s:        "=((1))",
			wantNode: &ast.NumberNode{Number: "1"},
		},
		{
			name:    "nested parentheses unbalanced",
			s:       "=((1)",
			wantErr: parse.ErrParse,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parse.Parse(tc.s)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Parse(%q) = _, %v; want %v", tc.s, err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.wantNode, got); diff != "" {
				t.Errorf("Parse(%q) mismatch (-want +got):\n%s", tc.s, diff)
			}
		})
	}
}
