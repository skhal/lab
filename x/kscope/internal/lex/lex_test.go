// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lex_test

import (
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/skhal/lab/x/kscope/internal/lex"
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
			name:    "invalid token",
			s:       "%", // pick any unused token character
			wantErr: lex.ErrScan,
		},
		{
			name: "left parenthesis",
			s:    "(",
			want: []lex.Token{
				{Kind: lex.TokLpar, Val: "("},
			},
		},
		{
			name: "right parenthesis",
			s:    ")",
			want: []lex.Token{
				{Kind: lex.TokRpar, Val: ")"},
			},
		},
		{
			name: "comma",
			s:    ",",
			want: []lex.Token{
				{Kind: lex.TokComma, Val: ","},
			},
		},
		{
			name: "assign",
			s:    "=",
			want: []lex.Token{
				{Kind: lex.TokAssign, Val: "="},
			},
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
				{Kind: lex.TokNum, Val: "123"},
			},
		},
		{
			name: "float",
			s:    "1.2",
			want: []lex.Token{
				{Kind: lex.TokNum, Val: "1.2"},
			},
		},
		{
			name: "float no fractional part",
			s:    "12.",
			want: []lex.Token{
				{Kind: lex.TokNum, Val: "12."},
			},
		},
		{
			name: "float no integral part",
			s:    ".12",
			want: []lex.Token{
				{Kind: lex.TokNum, Val: ".12"},
			},
		},
		{
			name: "space prefix",
			s:    "\t 1.2",
			want: []lex.Token{
				{Kind: lex.TokNum, Val: "1.2"},
			},
		},
		{
			name: "multiple",
			//  0    5    10  <- index
			//  |    |    |
			s: "1.2 3  4. .5",
			want: []lex.Token{
				{Kind: lex.TokNum, Val: "1.2"},
				{Kind: lex.TokNum, Val: "3"},
				{Kind: lex.TokNum, Val: "4."},
				{Kind: lex.TokNum, Val: ".5"},
			},
		},
	}
	testLex(t, tests)
}

func TestLex_binop(t *testing.T) {
	tests := []testCase{
		{
			name: "plus",
			s:    "+",
			want: []lex.Token{
				{Kind: lex.TokPlus, Val: "+"},
			},
		},
		{
			name: "minus",
			s:    "-",
			want: []lex.Token{
				{Kind: lex.TokMinus, Val: "-"},
			},
		},
		{
			name: "multiply",
			s:    "*",
			want: []lex.Token{
				{Kind: lex.TokMul, Val: "*"},
			},
		},
		{
			name: "divide",
			s:    "/",
			want: []lex.Token{
				{Kind: lex.TokDiv, Val: "/"},
			},
		},
	}
	testLex(t, tests)
}

func TestLex_ident(t *testing.T) {
	tests := []testCase{
		{
			name: "letter",
			s:    "a",
			want: []lex.Token{
				{Kind: lex.TokIdent, Val: "a"},
			},
		},
		{
			name: "letters",
			s:    "abc",
			want: []lex.Token{
				{Kind: lex.TokIdent, Val: "abc"},
			},
		},
		{
			name: "alnum",
			s:    "a1b2",
			want: []lex.Token{
				{Kind: lex.TokIdent, Val: "a1b2"},
			},
		},
		{
			name: "space prefix",
			s:    "\t a1",
			want: []lex.Token{
				{Kind: lex.TokIdent, Val: "a1"},
			},
		},
		{
			name: "multiple",
			//  0    5    10  <- index
			//  |    |    |
			s: "a a1 a2b a34b5",
			want: []lex.Token{
				{Kind: lex.TokIdent, Val: "a"},
				{Kind: lex.TokIdent, Val: "a1"},
				{Kind: lex.TokIdent, Val: "a2b"},
				{Kind: lex.TokIdent, Val: "a34b5"},
			},
		},
	}
	testLex(t, tests)
}

func TestLex_mix(t *testing.T) {
	tests := []testCase{
		{
			name: "number and identifier",
			s:    "1.2 a3",
			want: []lex.Token{
				{Kind: lex.TokNum, Val: "1.2"},
				{Kind: lex.TokIdent, Val: "a3"},
			},
		},
		{
			name: "comment",
			s:    "# test",
			want: []lex.Token{
				{Kind: lex.TokComment, Val: "# test"},
			},
		},
		{
			name: "comment multi-line",
			s:    "# test a\n# test b",
			want: []lex.Token{
				{Kind: lex.TokComment, Val: "# test a"},
				{Kind: lex.TokComment, Val: "# test b"},
			},
		},
		{
			name: "number and comment",
			s:    "123 # test",
			want: []lex.Token{
				{Kind: lex.TokNum, Val: "123"},
				{Kind: lex.TokComment, Val: "# test"},
			},
		},
	}
	testLex(t, tests)
}

func TestLex_commands(t *testing.T) {
	tests := []testCase{
		{
			name: "def",
			s:    "def",
			want: []lex.Token{
				{Kind: lex.TokDef, Val: "def"},
			},
		},
		{
			name: "extern",
			s:    "extern",
			want: []lex.Token{
				{Kind: lex.TokExtern, Val: "extern"},
			},
		},
		{
			name: "var",
			s:    "var",
			want: []lex.Token{
				{Kind: lex.TokVar, Val: "var"},
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

			seq, _ := lx.Lex(tc.s)
			got := slices.Collect(seq)
			err := lx.Err()

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Lex(%q) error %v; want %v", tc.s, err, tc.wantErr)
			}
			if d := cmp.Diff(tc.want, got, cmpopts.IgnoreUnexported(lex.Token{})); d != "" {
				t.Errorf("Lex(%q) mismatch (-want +got):\n%s", tc.s, d)
			}
		})
	}
}

func ExampleLexer_Lex() {
	s := `1 2
name 3`
	seq, positioner := (&lex.Lexer{}).Lex(s)
	for tk := range seq {
		fmt.Printf("%s %s\n", positioner.Pos(tk), tk)
	}
	// Output:
	// 1:1 number "1"
	// 1:3 number "2"
	// 2:1 identifier "name"
	// 2:6 number "3"
}
