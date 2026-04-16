// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parse_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/go/tests"
	"github.com/skhal/lab/x/sheet/internal/ast"
	"github.com/skhal/lab/x/sheet/internal/parse"
)

func TestParse_number(t *testing.T) {
	tt := []struct {
		name     string
		s        string
		wantNode *ast.NumberNode
		wantErr  error
	}{
		{
			name:    "empty",
			wantErr: parse.ErrParse,
		},
		{
			name:     "positive int",
			s:        "123",
			wantNode: &ast.NumberNode{Number: 123},
		},
		{
			name:     "positive float",
			s:        "1.23",
			wantNode: &ast.NumberNode{Number: 1.23},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parse.Parse(tc.s)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Parse(%q) = _, %v; want %v", tc.s, err, tc.wantErr)
			}
			diffOpts := []cmp.Option{
				tests.EquateFloat64(0.01), // equal within 1%
			}
			if diff := cmp.Diff(tc.wantNode, got, diffOpts...); diff != "" {
				t.Errorf("Parse(%q) mismatch (-want +got):\n%s", tc.s, diff)
			}
		})
	}
}

func TestParse_formula(t *testing.T) {
	tt := []struct {
		name     string
		s        string
		wantNode *ast.FormulaNode
		wantErr  error
	}{
		{
			name:    "empty",
			s:       "=",
			wantErr: parse.ErrParse,
		},
		{
			name:     "positive int",
			s:        "=123",
			wantNode: &ast.FormulaNode{Number: &ast.NumberNode{Number: 123}},
		},
		{
			name:     "positive float",
			s:        "=1.23",
			wantNode: &ast.FormulaNode{Number: &ast.NumberNode{Number: 1.23}},
		},
		{
			name:     "float without fractional part",
			s:        "=1.",
			wantNode: &ast.FormulaNode{Number: &ast.NumberNode{Number: 1.}},
		},
		{
			name:    "float must have integral part",
			s:       "=.1",
			wantErr: parse.ErrParse,
		},
		{
			name:    "multiple numbers",
			s:       "=1 2",
			wantErr: parse.ErrParse,
		},
		{
			name:    "unsupported token",
			s:       "=invalid",
			wantErr: parse.ErrParse,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parse.Parse(tc.s)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Parse(%q) = _, %v; want %v", tc.s, err, tc.wantErr)
			}
			diffOpts := []cmp.Option{
				tests.EquateFloat64(0.01), // equal within 1%
			}
			if diff := cmp.Diff(tc.wantNode, got, diffOpts...); diff != "" {
				t.Errorf("Parse(%q) mismatch (-want +got):\n%s", tc.s, diff)
			}
		})
	}
}
