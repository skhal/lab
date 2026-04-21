// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lex_test

import (
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/skhal/lab/x/sheet/internal/lex"
)

func TestLex(t *testing.T) {
	tests := []struct {
		name string
		b    []byte
		want []lex.Token
	}{
		{
			name: "empty",
		},
		{
			name: "number one digit",
			b:    []byte("1"),
			want: []lex.Token{
				{Type: lex.TokenNumber, Text: "1"},
			},
		},
		{
			name: "number two digits",
			b:    []byte("12"),
			want: []lex.Token{
				{Type: lex.TokenNumber, Text: "12"},
			},
		},
		{
			name: "number empty integral part",
			b:    []byte(".1"),
			want: []lex.Token{
				{Type: lex.TokenError, Err: lex.ErrLex},
			},
		},
		{
			name: "number empty fractional part",
			b:    []byte("1."),
			want: []lex.Token{
				{Type: lex.TokenNumber, Text: "1."},
			},
		},
		{
			name: "number non empty fractional part",
			b:    []byte("1.2"),
			want: []lex.Token{
				{Type: lex.TokenNumber, Text: "1.2"},
			},
		},
		{
			name: "invalid utf8 encoding",
			// https://cs.opensource.google/go/go/+/master:src/unicode/utf8/utf8_test.go;l=253;drc=925a3cdcd13472c8f78d51c9ce99a59e77d46eb4
			b: []byte("\x91\x80\x80\x80"),
			want: []lex.Token{
				{Type: lex.TokenError, Err: lex.ErrLex},
			},
		},
		{
			name: "operator plus",
			b:    []byte("+"),
			want: []lex.Token{
				{Type: lex.TokenPlus, Text: "+"},
			},
		},
		{
			name: "operator minus",
			b:    []byte("-"),
			want: []lex.Token{
				{Type: lex.TokenMinus, Text: "-"},
			},
		},
		{
			name: "left parenthesis",
			b:    []byte("("),
			want: []lex.Token{
				{Type: lex.TokenLpar, Text: "("},
			},
		},
		{
			name: "right parenthesis",
			b:    []byte(")"),
			want: []lex.Token{
				{Type: lex.TokenRpar, Text: ")"},
			},
		},
		{
			name: "identifier",
			b:    []byte("ABC123"),
			want: []lex.Token{
				{Type: lex.TokenIdent, Text: "ABC123"},
			},
		},
		{
			name: "identifier no digits",
			b:    []byte("ABC"),
			want: []lex.Token{
				{Type: lex.TokenIdent, Text: "ABC"},
			},
		},
		{
			name: "identifier mixed case",
			b:    []byte("AbC123"),
			want: []lex.Token{
				{Type: lex.TokenIdent, Text: "AbC123"},
			},
		},
		{
			name: "identifier mixed case no digits",
			b:    []byte("ABc"),
			want: []lex.Token{
				{Type: lex.TokenIdent, Text: "ABc"},
			},
		},
		{
			name: "comma",
			b:    []byte(","),
			want: []lex.Token{
				{Type: lex.TokenComma, Text: ","},
			},
		},
		{
			name: "range",
			b:    []byte("A1:A2"),
			want: []lex.Token{
				{Type: lex.TokenRange, Text: "A1:A2"},
			},
		},
		{
			name: "range misses second identifier",
			b:    []byte("A1:"),
			want: []lex.Token{
				{Type: lex.TokenError, Err: lex.ErrLex},
			},
		},
		{
			// a reference is a cell identifier, e.g. A123
			name: "range second identifier is not a reference",
			b:    []byte("A1:B"),
			want: []lex.Token{
				{Type: lex.TokenError, Err: lex.ErrLex},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := slices.Collect(lex.Lex(tc.b))

			if diff := cmp.Diff(tc.want, got, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Lex() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestLex_stop(t *testing.T) {
	tests := []struct {
		name string
		b    []byte
		want []lex.Token
	}{
		{
			name: "number",
			b:    []byte("1.2"),
			want: []lex.Token{
				{Type: lex.TokenNumber, Text: "1.2"},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got []lex.Token

			for tok := range lex.Lex(tc.b) {
				got = append(got, tok)
				break
			}

			if diff := cmp.Diff(tc.want, got, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Lex() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
