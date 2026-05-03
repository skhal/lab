// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parse_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/skhal/lab/x/kscope/internal/ast"
	"github.com/skhal/lab/x/kscope/internal/parse"
)

// diffFloatFractionPcent is a relative difference (RD) of two floating numbers.
// when RD is blow this value, the two numbers are considered equal.
const diffFloatFractionPcent = 0.001

type testCase struct {
	want    ast.Node
	wantErr error
	name    string
	text    string
}

func TestParser_Parse(t *testing.T) {
	tests := []testCase{
		{
			name: "empty",
		},
		{
			name: "number",
			text: "12.3",
			want: ast.Number{Val: 12.3},
		},
	}
	testParser_Parse(t, tests)
}

func TestParser_binop(t *testing.T) {
	tests := []testCase{
		{
			name: "plus",
			text: "1 + 2",
			want: ast.BinExpr{
				Op:    ast.BinOpPlus,
				Left:  ast.Number{Val: 1},
				Right: ast.Number{Val: 2},
			},
		},
		{
			name:    "plus misses left operand",
			text:    "+ 2",
			wantErr: parse.ErrParse,
		},
		{
			name:    "plus misses right operand",
			text:    "1 +",
			wantErr: parse.ErrParse,
		},
		{
			name: "minus",
			text: "1 - 2",
			want: ast.BinExpr{
				Op:    ast.BinOpMinus,
				Left:  ast.Number{Val: 1},
				Right: ast.Number{Val: 2},
			},
		},
		{
			name:    "minus misses left operand",
			text:    "- 2",
			wantErr: parse.ErrParse,
		},
		{
			name:    "minus misses right operand",
			text:    "1 -",
			wantErr: parse.ErrParse,
		},
	}
	testParser_Parse(t, tests)
}

func testParser_Parse(t *testing.T, tests []testCase) {
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parse.Parse(tc.text)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("unexpected error %v; want %v", err, tc.wantErr)
			}
			opts := []cmp.Option{
				cmpopts.EquateApprox(diffFloatFractionPcent, 0),
				cmpopts.EquateEmpty(),
			}
			if d := cmp.Diff(tc.want, got, opts...); d != "" {
				t.Errorf("mismatch (-want +got):\n%s", d)
				t.Logf("text:\n%s", tc.text)
			}
		})
	}
}

func ExampleParse() {
	const s = `
123
`
	n, err := parse.Parse(s)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(n)
	// Output:
	// 123.0
}
