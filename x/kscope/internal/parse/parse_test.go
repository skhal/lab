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
	}
	testParser_Parse(t, tests)
}

func TestParser_expr(t *testing.T) {
	tests := []testCase{
		{
			name: "plus",
			text: "var x = 1 + 2",
			want: ast.Var{
				Name: "x",
				Val: ast.BinExpr{
					Op:    ast.BinOpPlus,
					Left:  ast.Number{Val: 1},
					Right: ast.Number{Val: 2},
				},
			},
		},
		{
			name:    "plus misses left operand",
			text:    "var x = + 2",
			wantErr: parse.ErrParse,
		},
		{
			name:    "plus misses right operand",
			text:    "var x = 1 +",
			wantErr: parse.ErrParse,
		},
		{
			name: "minus",
			text: "var x = 1 - 2",
			want: ast.Var{
				Name: "x",
				Val: ast.BinExpr{
					Op:    ast.BinOpMinus,
					Left:  ast.Number{Val: 1},
					Right: ast.Number{Val: 2},
				},
			},
		},
		{
			name: "multiply",
			text: "var x = 1 * 2",
			want: ast.Var{
				Name: "x",
				Val: ast.BinExpr{
					Op:    ast.BinOpMul,
					Left:  ast.Number{Val: 1},
					Right: ast.Number{Val: 2},
				},
			},
		},
		{
			name: "divide",
			text: "var x = 1 / 2",
			want: ast.Var{
				Name: "x",
				Val: ast.BinExpr{
					Op:    ast.BinOpDiv,
					Left:  ast.Number{Val: 1},
					Right: ast.Number{Val: 2},
				},
			},
		},
		// same order
		{
			name: "plus and plus",
			text: "var x = 1 + 2 + 3",
			want: ast.Var{
				Name: "x",
				Val: ast.BinExpr{
					Op:   ast.BinOpPlus,
					Left: ast.Number{Val: 1},
					Right: ast.BinExpr{
						Op:    ast.BinOpPlus,
						Left:  ast.Number{Val: 2},
						Right: ast.Number{Val: 3},
					},
				},
			},
		},
		{
			name: "plus and minus",
			text: "var x = 1 + 2 - 3",
			want: ast.Var{
				Name: "x",
				Val: ast.BinExpr{
					Op:   ast.BinOpPlus,
					Left: ast.Number{Val: 1},
					Right: ast.BinExpr{
						Op:    ast.BinOpMinus,
						Left:  ast.Number{Val: 2},
						Right: ast.Number{Val: 3},
					},
				},
			},
		},
		{
			name: "plus and multiply",
			text: "var x = 1 + 2 * 3",
			want: ast.Var{
				Name: "x",
				Val: ast.BinExpr{
					Op:   ast.BinOpPlus,
					Left: ast.Number{Val: 1},
					Right: ast.BinExpr{
						Op:    ast.BinOpMul,
						Left:  ast.Number{Val: 2},
						Right: ast.Number{Val: 3},
					},
				},
			},
		},
		{
			name: "plus and divide",
			text: "var x = 1 + 2 / 3",
			want: ast.Var{
				Name: "x",
				Val: ast.BinExpr{
					Op:   ast.BinOpPlus,
					Left: ast.Number{Val: 1},
					Right: ast.BinExpr{
						Op:    ast.BinOpDiv,
						Left:  ast.Number{Val: 2},
						Right: ast.Number{Val: 3},
					},
				},
			},
		},
		{
			name: "minus and plus",
			text: "var x = 1 - 2 + 3",
			want: ast.Var{
				Name: "x",
				Val: ast.BinExpr{
					Op:   ast.BinOpMinus,
					Left: ast.Number{Val: 1},
					Right: ast.BinExpr{
						Op:    ast.BinOpPlus,
						Left:  ast.Number{Val: 2},
						Right: ast.Number{Val: 3},
					},
				},
			},
		},
		{
			name: "minus and minus",
			text: "var x = 1 - 2 - 3",
			want: ast.Var{
				Name: "x",
				Val: ast.BinExpr{
					Op:   ast.BinOpMinus,
					Left: ast.Number{Val: 1},
					Right: ast.BinExpr{
						Op:    ast.BinOpMinus,
						Left:  ast.Number{Val: 2},
						Right: ast.Number{Val: 3},
					},
				},
			},
		},
		{
			name: "minus and multiply",
			text: "var x = 1 - 2 * 3",
			want: ast.Var{
				Name: "x",
				Val: ast.BinExpr{
					Op:   ast.BinOpMinus,
					Left: ast.Number{Val: 1},
					Right: ast.BinExpr{
						Op:    ast.BinOpMul,
						Left:  ast.Number{Val: 2},
						Right: ast.Number{Val: 3},
					},
				},
			},
		},
		{
			name: "minus and divide",
			text: "var x = 1 - 2 / 3",
			want: ast.Var{
				Name: "x",
				Val: ast.BinExpr{
					Op:   ast.BinOpMinus,
					Left: ast.Number{Val: 1},
					Right: ast.BinExpr{
						Op:    ast.BinOpDiv,
						Left:  ast.Number{Val: 2},
						Right: ast.Number{Val: 3},
					},
				},
			},
		},
		{
			name: "multiply and plus",
			text: "var x = 1 * 2 + 3",
			want: ast.Var{
				Name: "x",
				Val: ast.BinExpr{
					Op: ast.BinOpPlus,
					Left: ast.BinExpr{
						Op:    ast.BinOpMul,
						Left:  ast.Number{Val: 1},
						Right: ast.Number{Val: 2},
					},
					Right: ast.Number{Val: 3},
				},
			},
		},
		{
			name: "multiply and minus",
			text: "var x = 1 * 2 - 3",
			want: ast.Var{
				Name: "x",
				Val: ast.BinExpr{
					Op: ast.BinOpMinus,
					Left: ast.BinExpr{
						Op:    ast.BinOpMul,
						Left:  ast.Number{Val: 1},
						Right: ast.Number{Val: 2},
					},
					Right: ast.Number{Val: 3},
				},
			},
		},
		{
			name: "multiply and multiply",
			text: "var x = 1 * 2 * 3",
			want: ast.Var{
				Name: "x",
				Val: ast.BinExpr{
					Op: ast.BinOpMul,
					Left: ast.BinExpr{
						Op:    ast.BinOpMul,
						Left:  ast.Number{Val: 1},
						Right: ast.Number{Val: 2},
					},
					Right: ast.Number{Val: 3},
				},
			},
		},
		{
			name: "multiply and divide",
			text: "var x = 1 * 2 / 3",
			want: ast.Var{
				Name: "x",
				Val: ast.BinExpr{
					Op: ast.BinOpDiv,
					Left: ast.BinExpr{
						Op:    ast.BinOpMul,
						Left:  ast.Number{Val: 1},
						Right: ast.Number{Val: 2},
					},
					Right: ast.Number{Val: 3},
				},
			},
		},
		{
			name: "divide and plus",
			text: "var x = 1 / 2 + 3",
			want: ast.Var{
				Name: "x",
				Val: ast.BinExpr{
					Op: ast.BinOpPlus,
					Left: ast.BinExpr{
						Op:    ast.BinOpDiv,
						Left:  ast.Number{Val: 1},
						Right: ast.Number{Val: 2},
					},
					Right: ast.Number{Val: 3},
				},
			},
		},
		{
			name: "divide and minus",
			text: "var x = 1 / 2 - 3",
			want: ast.Var{
				Name: "x",
				Val: ast.BinExpr{
					Op: ast.BinOpMinus,
					Left: ast.BinExpr{
						Op:    ast.BinOpDiv,
						Left:  ast.Number{Val: 1},
						Right: ast.Number{Val: 2},
					},
					Right: ast.Number{Val: 3},
				},
			},
		},
		{
			name: "divide and multiply",
			text: "var x = 1 / 2 * 3",
			want: ast.Var{
				Name: "x",
				Val: ast.BinExpr{
					Op: ast.BinOpMul,
					Left: ast.BinExpr{
						Op:    ast.BinOpDiv,
						Left:  ast.Number{Val: 1},
						Right: ast.Number{Val: 2},
					},
					Right: ast.Number{Val: 3},
				},
			},
		},
		{
			name: "divide and divide",
			text: "var x = 1 / 2 / 3",
			want: ast.Var{
				Name: "x",
				Val: ast.BinExpr{
					Op: ast.BinOpDiv,
					Left: ast.BinExpr{
						Op:    ast.BinOpDiv,
						Left:  ast.Number{Val: 1},
						Right: ast.Number{Val: 2},
					},
					Right: ast.Number{Val: 3},
				},
			},
		},
		{
			name: "group plus",
			text: "var x = (1 + 2)",
			want: ast.Var{
				Name: "x",
				Val: ast.BinExpr{
					Op:    ast.BinOpPlus,
					Left:  ast.Number{Val: 1},
					Right: ast.Number{Val: 2},
				},
			},
		},
		{
			name: "group prioritizes",
			text: "var x = 1 * (2 + 3)",
			want: ast.Var{
				Name: "x",
				Val: ast.BinExpr{
					Op:   ast.BinOpMul,
					Left: ast.Number{Val: 1},
					Right: ast.BinExpr{
						Op:    ast.BinOpPlus,
						Left:  ast.Number{Val: 2},
						Right: ast.Number{Val: 3},
					},
				},
			},
		},
		{
			name: "lhs is identifier",
			text: "var x = x + 1",
			want: ast.Var{
				Name: "x",
				Val: ast.BinExpr{
					Op:    ast.BinOpPlus,
					Left:  ast.Ident{Name: "x"},
					Right: ast.Number{Val: 1},
				},
			},
		},
		{
			name: "rhs is identifier",
			text: "var x = 1 + x",
			want: ast.Var{
				Name: "x",
				Val: ast.BinExpr{
					Op:    ast.BinOpPlus,
					Left:  ast.Number{Val: 1},
					Right: ast.Ident{Name: "x"},
				},
			},
		},
	}
	testParser_Parse(t, tests)
}

func TestParser_call(t *testing.T) {
	tests := []testCase{
		{
			name: "no args",
			text: "var x = test()",
			want: ast.Var{
				Name: "x",
				Val:  ast.Call{Name: "test"},
			},
		},
		{
			name:    "no right parenthesis",
			text:    "var x = test(",
			wantErr: parse.ErrParse,
		},
		{
			name: "one arg",
			text: "var x = test(1)",
			want: ast.Var{
				Name: "x",
				Val: ast.Call{
					Name: "test",
					Args: []ast.Node{ast.Number{Val: 1}},
				},
			},
		},
		{
			name:    "one arg no right parenthesis",
			text:    "var x = test(1",
			wantErr: parse.ErrParse,
		}, {
			name: "one arg expression",
			text: "var x = test(1 + 2)",
			want: ast.Var{
				Name: "x",
				Val: ast.Call{
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
		},
		{
			name: "two args",
			text: "var x = test(1, 2)",
			want: ast.Var{
				Name: "x",
				Val: ast.Call{
					Name: "test",
					Args: []ast.Node{
						ast.Number{Val: 1},
						ast.Number{Val: 2},
					},
				},
			},
		},
		{
			name:    "two args no right parenthesis",
			text:    "var x = test(1, 2",
			wantErr: parse.ErrParse,
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
			name: "one param",
			text: "def test(a) 1",
			want: ast.Func{
				Name:   "test",
				Params: []string{"a"},
				Body: []ast.Node{
					ast.Number{Val: 1},
				},
			},
		},
		{
			name: "two params",
			text: "def test(a, b) 1",
			want: ast.Func{
				Name:   "test",
				Params: []string{"a", "b"},
				Body: []ast.Node{
					ast.Number{Val: 1},
				},
			},
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

func TestParser_var(t *testing.T) {
	tests := []testCase{
		{
			name: "number",
			text: "var x = 1",
			want: ast.Var{
				Name: "x",
				Val:  ast.Number{Val: 1},
			},
		},
		{
			name:    "no identifier",
			text:    "var = 1",
			wantErr: parse.ErrParse,
		},
		{
			name:    "no assignment",
			text:    "var x 1",
			wantErr: parse.ErrParse,
		},
		{
			name:    "no value",
			text:    "var x =",
			wantErr: parse.ErrParse,
		},
	}
	testParser_Parse(t, tests)
}

func testParser_Parse(t *testing.T, tests []testCase) {
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parse.ParseExpr(tc.text)

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

func ExampleParseExpr() {
	const s = `
var a = 1
var b = a + 22
def c()
	a + b * 2
def d(x, y)
  x * c() + a
`
	n, err := parse.Parse(s)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(n)
	// Output:
	// var a = 1.0
	// var b = a + 22.0
	// def c()
	//   a + b * 2.0
	// def d(x, y)
	//   x * c() + a
}
