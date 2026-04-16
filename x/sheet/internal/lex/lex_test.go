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
			name: "not number",
			b:    []byte("abc"),
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
