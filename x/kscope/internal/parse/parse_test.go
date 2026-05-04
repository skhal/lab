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
		{
			name: "multiply",
			text: "1 * 2",
			want: ast.BinExpr{
				Op:    ast.BinOpMul,
				Left:  ast.Number{Val: 1},
				Right: ast.Number{Val: 2},
			},
		},
		{
			name:    "multiply misses left operand",
			text:    "* 2",
			wantErr: parse.ErrParse,
		},
		{
			name:    "multiply misses right operand",
			text:    "1 *",
			wantErr: parse.ErrParse,
		},
		{
			name: "divide",
			text: "1 / 2",
			want: ast.BinExpr{
				Op:    ast.BinOpDiv,
				Left:  ast.Number{Val: 1},
				Right: ast.Number{Val: 2},
			},
		},
		{
			name:    "divide misses left operand",
			text:    "/ 2",
			wantErr: parse.ErrParse,
		},
		{
			name:    "divide misses right operand",
			text:    "1 /",
			wantErr: parse.ErrParse,
		},
		{
			name: "plus and plus",
			text: "1 + 2 + 3",
			want: ast.BinExpr{
				Op:   ast.BinOpPlus,
				Left: ast.Number{Val: 1},
				Right: ast.BinExpr{
					Op:    ast.BinOpPlus,
					Left:  ast.Number{Val: 2},
					Right: ast.Number{Val: 3},
				},
			},
		},
		{
			name: "plus and minus",
			text: "1 + 2 - 3",
			want: ast.BinExpr{
				Op:   ast.BinOpPlus,
				Left: ast.Number{Val: 1},
				Right: ast.BinExpr{
					Op:    ast.BinOpMinus,
					Left:  ast.Number{Val: 2},
					Right: ast.Number{Val: 3},
				},
			},
		},
		{
			name: "plus and multiply",
			text: "1 + 2 * 3",
			want: ast.BinExpr{
				Op:   ast.BinOpPlus,
				Left: ast.Number{Val: 1},
				Right: ast.BinExpr{
					Op:    ast.BinOpMul,
					Left:  ast.Number{Val: 2},
					Right: ast.Number{Val: 3},
				},
			},
		},
		{
			name: "plus and divide",
			text: "1 + 2 / 3",
			want: ast.BinExpr{
				Op:   ast.BinOpPlus,
				Left: ast.Number{Val: 1},
				Right: ast.BinExpr{
					Op:    ast.BinOpDiv,
					Left:  ast.Number{Val: 2},
					Right: ast.Number{Val: 3},
				},
			},
		},
		{
			name: "minus and plus",
			text: "1 - 2 + 3",
			want: ast.BinExpr{
				Op:   ast.BinOpMinus,
				Left: ast.Number{Val: 1},
				Right: ast.BinExpr{
					Op:    ast.BinOpPlus,
					Left:  ast.Number{Val: 2},
					Right: ast.Number{Val: 3},
				},
			},
		},
		{
			name: "minus and minus",
			text: "1 - 2 - 3",
			want: ast.BinExpr{
				Op:   ast.BinOpMinus,
				Left: ast.Number{Val: 1},
				Right: ast.BinExpr{
					Op:    ast.BinOpMinus,
					Left:  ast.Number{Val: 2},
					Right: ast.Number{Val: 3},
				},
			},
		},
		{
			name: "minus and multiply",
			text: "1 - 2 * 3",
			want: ast.BinExpr{
				Op:   ast.BinOpMinus,
				Left: ast.Number{Val: 1},
				Right: ast.BinExpr{
					Op:    ast.BinOpMul,
					Left:  ast.Number{Val: 2},
					Right: ast.Number{Val: 3},
				},
			},
		},
		{
			name: "minus and divide",
			text: "1 - 2 / 3",
			want: ast.BinExpr{
				Op:   ast.BinOpMinus,
				Left: ast.Number{Val: 1},
				Right: ast.BinExpr{
					Op:    ast.BinOpDiv,
					Left:  ast.Number{Val: 2},
					Right: ast.Number{Val: 3},
				},
			},
		},
		{
			name: "multiply and plus",
			text: "1 * 2 + 3",
			want: ast.BinExpr{
				Op: ast.BinOpPlus,
				Left: ast.BinExpr{
					Op:    ast.BinOpMul,
					Left:  ast.Number{Val: 1},
					Right: ast.Number{Val: 2},
				},
				Right: ast.Number{Val: 3},
			},
		},
		{
			name: "multiply and minus",
			text: "1 * 2 - 3",
			want: ast.BinExpr{
				Op: ast.BinOpMinus,
				Left: ast.BinExpr{
					Op:    ast.BinOpMul,
					Left:  ast.Number{Val: 1},
					Right: ast.Number{Val: 2},
				},
				Right: ast.Number{Val: 3},
			},
		},
		{
			name: "multiply and multiply",
			text: "1 * 2 * 3",
			want: ast.BinExpr{
				Op: ast.BinOpMul,
				Left: ast.BinExpr{
					Op:    ast.BinOpMul,
					Left:  ast.Number{Val: 1},
					Right: ast.Number{Val: 2},
				},
				Right: ast.Number{Val: 3},
			},
		},
		{
			name: "multiply and divide",
			text: "1 * 2 / 3",
			want: ast.BinExpr{
				Op: ast.BinOpDiv,
				Left: ast.BinExpr{
					Op:    ast.BinOpMul,
					Left:  ast.Number{Val: 1},
					Right: ast.Number{Val: 2},
				},
				Right: ast.Number{Val: 3},
			},
		},
		{
			name: "divide and plus",
			text: "1 / 2 + 3",
			want: ast.BinExpr{
				Op: ast.BinOpPlus,
				Left: ast.BinExpr{
					Op:    ast.BinOpDiv,
					Left:  ast.Number{Val: 1},
					Right: ast.Number{Val: 2},
				},
				Right: ast.Number{Val: 3},
			},
		},
		{
			name: "divide and minus",
			text: "1 / 2 - 3",
			want: ast.BinExpr{
				Op: ast.BinOpMinus,
				Left: ast.BinExpr{
					Op:    ast.BinOpDiv,
					Left:  ast.Number{Val: 1},
					Right: ast.Number{Val: 2},
				},
				Right: ast.Number{Val: 3},
			},
		},
		{
			name: "divide and multiply",
			text: "1 / 2 * 3",
			want: ast.BinExpr{
				Op: ast.BinOpMul,
				Left: ast.BinExpr{
					Op:    ast.BinOpDiv,
					Left:  ast.Number{Val: 1},
					Right: ast.Number{Val: 2},
				},
				Right: ast.Number{Val: 3},
			},
		},
		{
			name: "divide and divide",
			text: "1 / 2 / 3",
			want: ast.BinExpr{
				Op: ast.BinOpDiv,
				Left: ast.BinExpr{
					Op:    ast.BinOpDiv,
					Left:  ast.Number{Val: 1},
					Right: ast.Number{Val: 2},
				},
				Right: ast.Number{Val: 3},
			},
		},
	}
	testParser_Parse(t, tests)
}

func TestParser_call(t *testing.T) {
	tests := []testCase{
		{
			name: "no args",
			text: "test()",
			want: ast.Call{
				Name: "test",
			},
		},
		{
			name:    "no right parenthesis",
			text:    "test(",
			wantErr: parse.ErrParse,
		},
		{
			name: "one arg",
			text: "test(1)",
			want: ast.Call{
				Name: "test",
				Args: []ast.Node{
					ast.Number{Val: 1},
				},
			},
		},
		{
			name:    "one arg no right parenthesis",
			text:    "test(1",
			wantErr: parse.ErrParse,
		},
		{
			name: "one arg expression",
			text: "test(1 + 2)",
			want: ast.Call{
				Name: "test",
				Args: []ast.Node{
					ast.BinExpr{
						Op:    ast.BinOpPlus,
						Left:  ast.Number{Val: 1},
						Right: ast.Number{Val: 2},
					},
				},
			},
		},
		{
			name: "two args",
			text: "test(1, 2)",
			want: ast.Call{
				Name: "test",
				Args: []ast.Node{
					ast.Number{Val: 1},
					ast.Number{Val: 2},
				},
			},
		},
		{
			name:    "two args no right parenthesis",
			text:    "test(1, 2",
			wantErr: parse.ErrParse,
		},
		{
			name: "two arg expressions",
			text: "test(1 + 2, 3 * 4)",
			want: ast.Call{
				Name: "test",
				Args: []ast.Node{
					ast.BinExpr{
						Op:    ast.BinOpPlus,
						Left:  ast.Number{Val: 1},
						Right: ast.Number{Val: 2},
					},
					ast.BinExpr{
						Op:    ast.BinOpMul,
						Left:  ast.Number{Val: 3},
						Right: ast.Number{Val: 4},
					},
				},
			},
		},
	}
	testParser_Parse(t, tests)
}

func TestParser_func(t *testing.T) {
	tests := []testCase{
		{
			name: "body is a number",
			text: "def test() 1",
			want: ast.Func{
				Name: "test",
				Body: []ast.Node{
					ast.Number{Val: 1},
				},
			},
		},
		{
			name:    "no identifier",
			text:    "def () 1",
			wantErr: parse.ErrParse,
		},
		{
			name:    "no left parenthesis",
			text:    "def test) 1",
			wantErr: parse.ErrParse,
		},
		{
			name:    "no right parenthesis",
			text:    "def test( 1",
			wantErr: parse.ErrParse,
		},
		{
			name:    "no body",
			text:    "def test()",
			wantErr: parse.ErrParse,
		},
		{
			name: "body is a binary expression",
			text: "def test() 1 + 2",
			want: ast.Func{
				Name: "test",
				Body: []ast.Node{
					ast.BinExpr{
						Op:    ast.BinOpPlus,
						Left:  ast.Number{Val: 1},
						Right: ast.Number{Val: 2},
					},
				},
			},
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
