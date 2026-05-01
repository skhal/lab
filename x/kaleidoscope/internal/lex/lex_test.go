// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lex_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/x/kaleidoscope/internal/lex"
)

type testCase struct {
	wantErr error
	name    string
	s       string
	want    []lex.Token
}

func TestLex(t *testing.T) {
	tests := []testCase{
		{
			name: "empty",
		},
		{
			name:    "not number",
			s:       "abc",
			wantErr: lex.ErrScan,
		},
	}
	testLex(t, tests)
}

func TestLex_number(t *testing.T) {
	tests := []testCase{
		{
			name: "integer",
			s:    "123",
			want: []lex.Token{
				{Kind: lex.TokNum, Val: "123", Pos: lex.Position{0, 3}},
			},
		},
		{
			name: "float",
			s:    "1.2",
			want: []lex.Token{
				{Kind: lex.TokNum, Val: "1.2", Pos: lex.Position{0, 3}},
			},
		},
		{
			name: "float no fractional part",
			s:    "12.",
			want: []lex.Token{
				{Kind: lex.TokNum, Val: "12.", Pos: lex.Position{0, 3}},
			},
		},
		{
			name: "float no integral part",
			s:    ".12",
			want: []lex.Token{
				{Kind: lex.TokNum, Val: ".12", Pos: lex.Position{0, 3}},
			},
		},
		{
			name: "space prefix",
			s:    "\t 1.2",
			want: []lex.Token{
				{Kind: lex.TokNum, Val: "1.2", Pos: lex.Position{2, 5}},
			},
		},
		{
			name: "multiple",
			//  0    5    10  <- index
			//  |    |    |
			s: "1.2 3  4. .5",
			want: []lex.Token{
				{Kind: lex.TokNum, Val: "1.2", Pos: lex.Position{0, 3}},
				{Kind: lex.TokNum, Val: "3", Pos: lex.Position{4, 5}},
				{Kind: lex.TokNum, Val: "4.", Pos: lex.Position{7, 9}},
				{Kind: lex.TokNum, Val: ".5", Pos: lex.Position{10, 12}},
			},
		},
	}
	testLex(t, tests)
}

func testLex(t *testing.T, tests []testCase) {
	t.Helper()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var lx lex.Lexer

			got := slices.Collect(lx.Lex(tc.s))
			err := lx.Err()

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Lex(%q) error %v; want %v", tc.s, err, tc.wantErr)
			}
			if d := cmp.Diff(tc.want, got); d != "" {
				t.Errorf("Lex(%q) mismatch (-want +got):\n%s", tc.s, d)
			}
		})
	}
}
