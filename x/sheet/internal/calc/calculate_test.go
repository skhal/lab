// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package calc_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/skhal/lab/x/sheet/internal/ast"
	"github.com/skhal/lab/x/sheet/internal/calc"
)

const invalidOperator = "~!@#$%^&*"

func TestCalculate(t *testing.T) {
	tests := []struct {
		name    string
		node    ast.Node
		want    float64
		wantErr error
	}{
		{
			name: "number node",
			node: &ast.NumberNode{Number: "123"},
			want: 123,
		},
		{
			name: "unsupported operator",
			node: &ast.BinOpNode{
				Op:    invalidOperator,
				Left:  &ast.NumberNode{Number: "1"},
				Right: &ast.NumberNode{Number: "2"},
			},
			wantErr: calc.ErrCalculate,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := calc.Calculate(tc.node, newTestRefCalculator(t, nil))

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Calculate() = _, %v; want %v", err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("Calculate() = %f, _; want %f", got, tc.want)
			}
		})
	}
}

func TestCalculate_plus(t *testing.T) {
	tests := []struct {
		name    string
		node    ast.Node
		want    float64
		wantErr error
	}{
		{
			name: "operator plus",
			node: &ast.BinOpNode{
				Op:    "+",
				Left:  &ast.NumberNode{Number: "1"},
				Right: &ast.NumberNode{Number: "2"},
			},
			want: 3,
		},
		{
			name: "operator plus no operands",
			node: &ast.BinOpNode{
				Op: "+",
			},
			wantErr: calc.ErrCalculate,
		},
		{
			name: "operator plus one operand",
			node: &ast.BinOpNode{
				Op:   "+",
				Left: &ast.NumberNode{Number: "1"},
			},
			wantErr: calc.ErrCalculate,
		},
		{
			name: "operator plus recurse first",
			node: &ast.BinOpNode{
				Op: "+",
				Left: &ast.BinOpNode{
					Op:    "+",
					Left:  &ast.NumberNode{Number: "1"},
					Right: &ast.NumberNode{Number: "2"},
				},
				Right: &ast.NumberNode{Number: "3"},
			},
			want: 6,
		},
		{
			name: "operator plus recurse first fails",
			node: &ast.BinOpNode{
				Op: "+",
				Left: &ast.BinOpNode{
					Op:    invalidOperator,
					Left:  &ast.NumberNode{Number: "1"},
					Right: &ast.NumberNode{Number: "2"},
				},
				Right: &ast.NumberNode{Number: "3"},
			},
			wantErr: calc.ErrCalculate,
		},
		{
			name: "operator plus recurse second",
			node: &ast.BinOpNode{
				Op:   "+",
				Left: &ast.NumberNode{Number: "1"},
				Right: &ast.BinOpNode{
					Op:    "+",
					Left:  &ast.NumberNode{Number: "2"},
					Right: &ast.NumberNode{Number: "3"},
				},
			},
			want: 6,
		},
		{
			name: "operator plus recurse second fails",
			node: &ast.BinOpNode{
				Op:   "+",
				Left: &ast.NumberNode{Number: "1"},
				Right: &ast.BinOpNode{
					Op:    invalidOperator,
					Left:  &ast.NumberNode{Number: "2"},
					Right: &ast.NumberNode{Number: "3"},
				},
			},
			wantErr: calc.ErrCalculate,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := calc.Calculate(tc.node, newTestRefCalculator(t, nil))

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Calculate() = _, %v; want %v", err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("Calculate() = %f, _; want %f", got, tc.want)
			}
		})
	}
}

func TestCalculate_minus(t *testing.T) {
	tests := []struct {
		name    string
		node    ast.Node
		want    float64
		wantErr error
	}{
		{
			name: "operator minus",
			node: &ast.BinOpNode{
				Op:    "-",
				Left:  &ast.NumberNode{Number: "1"},
				Right: &ast.NumberNode{Number: "2"},
			},
			want: -1,
		},
		{
			name: "operator minus no operands",
			node: &ast.BinOpNode{
				Op: "-",
			},
			wantErr: calc.ErrCalculate,
		},
		{
			name: "operator minus one operand",
			node: &ast.BinOpNode{
				Op:   "-",
				Left: &ast.NumberNode{Number: "1"},
			},
			wantErr: calc.ErrCalculate,
		},
		{
			name: "operator minus recurse first",
			node: &ast.BinOpNode{
				Op: "-",
				Left: &ast.BinOpNode{
					Op:    "-",
					Left:  &ast.NumberNode{Number: "1"},
					Right: &ast.NumberNode{Number: "2"},
				},
				Right: &ast.NumberNode{Number: "3"},
			},
			want: -4,
		},
		{
			name: "operator minus recurse first fails",
			node: &ast.BinOpNode{
				Op: "-",
				Left: &ast.BinOpNode{
					Op:    invalidOperator,
					Left:  &ast.NumberNode{Number: "1"},
					Right: &ast.NumberNode{Number: "2"},
				},
				Right: &ast.NumberNode{Number: "3"},
			},
			wantErr: calc.ErrCalculate,
		},
		{
			name: "operator minus recurse second",
			node: &ast.BinOpNode{
				Op:   "-",
				Left: &ast.NumberNode{Number: "1"},
				Right: &ast.BinOpNode{
					Op:    "-",
					Left:  &ast.NumberNode{Number: "2"},
					Right: &ast.NumberNode{Number: "3"},
				},
			},
			want: 2,
		},
		{
			name: "operator minus recurse second fails",
			node: &ast.BinOpNode{
				Op:   "-",
				Left: &ast.NumberNode{Number: "1"},
				Right: &ast.BinOpNode{
					Op:    invalidOperator,
					Left:  &ast.NumberNode{Number: "2"},
					Right: &ast.NumberNode{Number: "3"},
				},
			},
			wantErr: calc.ErrCalculate,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := calc.Calculate(tc.node, newTestRefCalculator(t, nil))

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Calculate() = _, %v; want %v", err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("Calculate() = %f, _; want %f", got, tc.want)
			}
		})
	}
}

func TestCalculate_reference(t *testing.T) {
	tests := []struct {
		name    string
		node    ast.Node
		refs    map[string]testCell
		want    float64
		wantErr error
	}{
		{
			name: "reference",
			node: &ast.RefNode{
				Ref: "A1",
			},
			refs: map[string]testCell{
				"A1": {res: 123},
			},
			want: 123,
		},
		{
			name: "missing reference",
			node: &ast.RefNode{
				Ref: "A1",
			},
			wantErr: errTest,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := calc.Calculate(tc.node, newTestRefCalculator(t, tc.refs))

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Calculate() = _, %v; want %v", err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("Calculate() = %f, _; want %f", got, tc.want)
			}
		})
	}
}

func TestCalculate_call(t *testing.T) {
	tests := []struct {
		name    string
		node    ast.Node
		want    float64
		wantErr bool
	}{
		{
			name: "unsupported function",
			node: &ast.CallNode{
				Name: "UNSUPPORTED",
			},
			wantErr: true,
		},
		{
			name: "disabled function",
			node: &ast.CallNode{
				Name: "MIN",
				Args: []ast.Node{
					&ast.NumberNode{Number: "3"},
					&ast.NumberNode{Number: "2"},
				},
			},
			wantErr: true,
		},
		{
			name: "sum no args",
			node: &ast.CallNode{
				Name: "SUM",
			},
		},
		{
			name: "sum with static args",
			node: &ast.CallNode{
				Name: "SUM",
				Args: []ast.Node{
					&ast.NumberNode{Number: "123"},
				},
			},
			want: 123,
		},
		{
			name: "fail to calculate arg",
			node: &ast.CallNode{
				Name: "SUM",
				Args: []ast.Node{
					&ast.RefNode{Ref: "A1"},
				},
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := calc.Calculate(tc.node, newTestRefCalculator(t, nil))

			switch {
			case tc.wantErr && err == nil:
				t.Error("Calculate() missing error")
			case !tc.wantErr && err != nil:
				t.Errorf("Calculate() unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("Calculate() = %f, _; want %f", got, tc.want)
			}
		})
	}
}

var errTest = errors.New("test error")

type testCell struct {
	res float64
	err error
}

type testRefCalculator struct {
	refs map[string]testCell
}

func newTestRefCalculator(t *testing.T, refs map[string]testCell) *testRefCalculator {
	t.Helper()
	return &testRefCalculator{refs}
}

func (rc *testRefCalculator) Calculate(id string) (float64, error) {
	c, ok := rc.refs[id]
	if !ok {
		return 0, fmt.Errorf("%w: ref %s does not exist", errTest, id)
	}
	if c.err != nil {
		return 0, c.err
	}
	return c.res, nil
}
