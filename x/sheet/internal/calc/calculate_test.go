// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package calc_test

import (
	"errors"
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
			got, err := calc.Calculate(tc.node)

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
			got, err := calc.Calculate(tc.node)

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
			got, err := calc.Calculate(tc.node)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Calculate() = _, %v; want %v", err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("Calculate() = %f, _; want %f", got, tc.want)
			}
		})
	}
}
