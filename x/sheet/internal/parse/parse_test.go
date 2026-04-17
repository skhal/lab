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

func TestParse_number(t *testing.T) {
	tt := []struct {
		name     string
		s        string
		wantNode ast.Node
		wantErr  error
	}{
		{
			name: "empty",
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
			if diff := cmp.Diff(tc.wantNode, got); diff != "" {
				t.Errorf("Parse(%q) mismatch (-want +got):\n%s", tc.s, diff)
			}
		})
	}
}
