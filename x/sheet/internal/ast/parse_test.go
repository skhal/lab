// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/x/sheet/internal/ast"
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
			wantErr: ast.ErrParse,
		},
		{
			name:    "not number",
			s:       "abc",
			wantErr: ast.ErrParse,
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
			got, err := ast.Parse(tc.s)

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
			wantErr: ast.ErrParse,
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
			wantErr: ast.ErrParse,
		},
		{
			name:    "unsupported token",
			s:       "=invalid",
			wantErr: ast.ErrParse,
		},
		{
			name:    "missing operator",
			s:       "=1 2",
			wantErr: ast.ErrParse,
		},
		{
			name: "plus",
			s:    "=1+2",
			wantNode: &ast.BinOpNode{
				Op:    "+",
				Left:  &ast.NumberNode{Number: "1"},
				Right: &ast.NumberNode{Number: "2"},
			},
		},
		{
			name:    "plus no left operand",
			s:       "=+ 2",
			wantErr: ast.ErrParse,
		},
		{
			name:    "plus no right operand",
			s:       "=1 +",
			wantErr: ast.ErrParse,
		},
		{
			name:    "plus invalid right operand",
			s:       "=1 + +",
			wantErr: ast.ErrParse,
		},
		{
			name: "minus",
			s:    "=1-2",
			wantNode: &ast.BinOpNode{
				Op:    "-",
				Left:  &ast.NumberNode{Number: "1"},
				Right: &ast.NumberNode{Number: "2"},
			},
		},
		{
			name:    "minus no left operand",
			s:       "=- 2",
			wantErr: ast.ErrParse,
		},
		{
			name:    "minus no right operand",
			s:       "=1 -",
			wantErr: ast.ErrParse,
		},
		{
			name:    "minus invalid right operand",
			s:       "=1 - -",
			wantErr: ast.ErrParse,
		},
		{
			name: "multiply",
			s:    "=1*2",
			wantNode: &ast.BinOpNode{
				Op:    "*",
				Left:  &ast.NumberNode{Number: "1"},
				Right: &ast.NumberNode{Number: "2"},
			},
		},
		{
			name:    "multiply no left operand",
			s:       "=* 2",
			wantErr: ast.ErrParse,
		},
		{
			name:    "multiply no right operand",
			s:       "=1 *",
			wantErr: ast.ErrParse,
		},
		{
			name:    "multiply invalid right operand",
			s:       "=1 * *",
			wantErr: ast.ErrParse,
		},
		{
			name: "divide",
			s:    "=1/2",
			wantNode: &ast.BinOpNode{
				Op:    "/",
				Left:  &ast.NumberNode{Number: "1"},
				Right: &ast.NumberNode{Number: "2"},
			},
		},
		{
			name:    "divide no left operand",
			s:       "=/ 2",
			wantErr: ast.ErrParse,
		},
		{
			name:    "divide no right operand",
			s:       "=1 /",
			wantErr: ast.ErrParse,
		},
		{
			name:    "divide invalid right operand",
			s:       "=1 / /",
			wantErr: ast.ErrParse,
		},
		{
			name: "multiply then plus",
			s:    "=1 * 2 + 3",
			wantNode: &ast.BinOpNode{
				Op: "+",
				Left: &ast.BinOpNode{
					Op:    "*",
					Left:  &ast.NumberNode{Number: "1"},
					Right: &ast.NumberNode{Number: "2"},
				},
				Right: &ast.NumberNode{Number: "3"},
			},
		},
		{
			name: "multiply then parenthesized plus",
			s:    "=1 * (2 + 3)",
			wantNode: &ast.BinOpNode{
				Op:   "*",
				Left: &ast.NumberNode{Number: "1"},
				Right: &ast.BinOpNode{
					Op:    "+",
					Left:  &ast.NumberNode{Number: "2"},
					Right: &ast.NumberNode{Number: "3"},
				},
			},
		},
		{
			name: "plus then multiply",
			s:    "=1 + 2 * 3",
			wantNode: &ast.BinOpNode{
				Op:   "+",
				Left: &ast.NumberNode{Number: "1"},
				Right: &ast.BinOpNode{
					Op:    "*",
					Left:  &ast.NumberNode{Number: "2"},
					Right: &ast.NumberNode{Number: "3"},
				},
			},
		},
		{
			name: "parenthesized plus then multiply",
			s:    "=(1 + 2) * 3",
			wantNode: &ast.BinOpNode{
				Op: "*",
				Left: &ast.BinOpNode{
					Op:    "+",
					Left:  &ast.NumberNode{Number: "1"},
					Right: &ast.NumberNode{Number: "2"},
				},
				Right: &ast.NumberNode{Number: "3"},
			},
		},
		{
			name:     "parentheses",
			s:        "=(1)",
			wantNode: &ast.NumberNode{Number: "1"},
		},
		{
			name:    "missing right parenthesis",
			s:       "=(1",
			wantErr: ast.ErrParse,
		},
		{
			name:    "unbalanced parentheses",
			s:       "=(1 2",
			wantErr: ast.ErrParse,
		},
		{
			name:    "missing left parenthesis",
			s:       "=1)",
			wantErr: ast.ErrParse,
		},
		{
			name:    "parentheses without expression",
			s:       "=()",
			wantErr: ast.ErrParse,
		},
		{
			name:     "nested parentheses",
			s:        "=((1))",
			wantNode: &ast.NumberNode{Number: "1"},
		},
		{
			name:    "nested parentheses unbalanced",
			s:       "=((1)",
			wantErr: ast.ErrParse,
		},
		{
			name:     "identifier",
			s:        "=A123",
			wantNode: &ast.RefNode{Ref: "A123"},
		},
		{
			name:    "invalid identifier",
			s:       "=ABC123",
			wantErr: ast.ErrParse,
		},
		{
			name:    "identifier must has upper case",
			s:       "=abc123",
			wantErr: ast.ErrParse,
		},
		{
			name:     "call with empty args",
			s:        "=SUM()",
			wantNode: &ast.CallNode{Name: "SUM"},
		},
		{
			name:    "call with invalid name",
			s:       "=SUM123()",
			wantErr: ast.ErrParse,
		},
		{
			name:    "call without right parenthesis",
			s:       "=SUM123(",
			wantErr: ast.ErrParse,
		},
		{
			name: "call with one literal arg",
			s:    "=SUM(123)",
			wantNode: &ast.CallNode{
				Name: "SUM",
				Args: []ast.Node{
					&ast.NumberNode{Number: "123"},
				},
			},
		},
		{
			name: "call with one literal args",
			s:    "=SUM(1, 2)",
			wantNode: &ast.CallNode{
				Name: "SUM",
				Args: []ast.Node{
					&ast.NumberNode{Number: "1"},
					&ast.NumberNode{Number: "2"},
				},
			},
		},
		{
			name: "call with one expr args",
			s:    "=SUM(1 + 2)",
			wantNode: &ast.CallNode{
				Name: "SUM",
				Args: []ast.Node{
					&ast.BinOpNode{
						Op:    "+",
						Left:  &ast.NumberNode{Number: "1"},
						Right: &ast.NumberNode{Number: "2"},
					},
				},
			},
		},
		{
			name:    "call with invalid expr args",
			s:       "=SUM(1 2)",
			wantErr: ast.ErrParse,
		},
		{
			name:     "range",
			s:        "=A1:A3",
			wantNode: &ast.RangeNode{From: "A1", To: "A3"},
		},
		{
			name:    "range misses first identifier",
			s:       "=:A3",
			wantErr: ast.ErrParse,
		},
		{
			name:    "range misses second identifier",
			s:       "=A1:",
			wantErr: ast.ErrParse,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ast.Parse(tc.s)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Parse(%q) = _, %v; want %v", tc.s, err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.wantNode, got); diff != "" {
				t.Errorf("Parse(%q) mismatch (-want +got):\n%s", tc.s, diff)
			}
		})
	}
}
