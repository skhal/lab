// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vm_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/skhal/lab/x/sheet/internal/ast"
	"github.com/skhal/lab/x/sheet/internal/vm"
)

func TestCompile_nonFormula(t *testing.T) {
	tests := []struct {
		name    string
		ast     ast.Node
		want    vm.InstructionsSet
		wantErr error
	}{
		{
			name: "empty",
		},
		{
			name: "integer",
			ast:  &ast.NumberNode{Number: "123"},
			want: vm.InstructionsSet{
				Instructions: []vm.Inst{
					newNumber(t, 123),
				},
			},
		},
		{
			name: "float",
			ast:  &ast.NumberNode{Number: "1.23"},
			want: vm.InstructionsSet{
				Instructions: []vm.Inst{
					newNumber(t, 1.23),
				},
			},
		},
		{
			name:    "invalid number",
			ast:     &ast.NumberNode{Number: "test"},
			wantErr: vm.ErrCompile,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := vm.Compile(tc.ast)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Compile() unexpected error: %v; want %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("Compile() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCompile_formula(t *testing.T) {
	tests := []struct {
		name    string
		ast     ast.Node
		want    vm.InstructionsSet
		wantErr error
	}{
		{
			name: "operator plus",
			ast: &ast.BinOpNode{
				Op:    "+",
				Left:  &ast.NumberNode{Number: "1"},
				Right: &ast.NumberNode{Number: "2"},
			},
			want: vm.InstructionsSet{
				Instructions: []vm.Inst{
					newNumber(t, 1),
					newNumber(t, 2),
					newBinOp(t, vm.BinOpPlus),
				},
			},
		},
		{
			name: "operator minus",
			ast: &ast.BinOpNode{
				Op:    "-",
				Left:  &ast.NumberNode{Number: "1"},
				Right: &ast.NumberNode{Number: "2"},
			},
			want: vm.InstructionsSet{
				Instructions: []vm.Inst{
					newNumber(t, 1),
					newNumber(t, 2),
					newBinOp(t, vm.BinOpMinus),
				},
			},
		},
		{
			name: "operator plus invalid left operand",
			ast: &ast.BinOpNode{
				Op:    "+",
				Left:  &ast.NumberNode{Number: "test"},
				Right: &ast.NumberNode{Number: "2"},
			},
			wantErr: vm.ErrCompile,
		},
		{
			name: "operator plus invalid right operand",
			ast: &ast.BinOpNode{
				Op:    "+",
				Left:  &ast.NumberNode{Number: "1"},
				Right: &ast.NumberNode{Number: "test"},
			},
			wantErr: vm.ErrCompile,
		},
		{
			name: "unsupported binary operator",
			ast: &ast.BinOpNode{
				Op:    "test",
				Left:  &ast.NumberNode{Number: "1"},
				Right: &ast.NumberNode{Number: "2"},
			},
			wantErr: vm.ErrCompile,
		},
		{
			name: "identifier",
			ast:  &ast.RefNode{Ref: "ABC123"},
			want: vm.InstructionsSet{
				Instructions: []vm.Inst{
					newRef(t, "ABC123"),
				},
			},
		},
		{
			name:    "range",
			ast:     &ast.RangeNode{From: "A1", To: "A3"},
			wantErr: vm.ErrCompile,
		},
		{
			name: "call no args",
			ast:  &ast.CallNode{Name: "SUM"},
			want: vm.InstructionsSet{
				Instructions: []vm.Inst{
					newCall(t, vm.Call{Func: vm.FuncSum}),
				},
			},
		},
		{
			name:    "unsupported call",
			ast:     &ast.CallNode{Name: "test call"},
			wantErr: vm.ErrCompile,
		},
		{
			name: "call with one literal arg",
			ast: &ast.CallNode{
				Name: "SUM",
				Args: []ast.Node{
					&ast.NumberNode{Number: "123"},
				},
			},
			want: vm.InstructionsSet{
				Instructions: []vm.Inst{
					newNumber(t, 123),
					newCall(t, vm.Call{Func: vm.FuncSum, Args: 1}),
				},
			},
		},
		{
			name: "call with invalid argument",
			ast: &ast.CallNode{
				Name: "SUM",
				Args: []ast.Node{
					&ast.NumberNode{Number: "test"},
				},
			},
			wantErr: vm.ErrCompile,
		},
		{
			name: "call with two literal args",
			ast: &ast.CallNode{
				Name: "SUM",
				Args: []ast.Node{
					&ast.NumberNode{Number: "1"},
					&ast.NumberNode{Number: "2"},
				},
			},
			want: vm.InstructionsSet{
				Instructions: []vm.Inst{
					newNumber(t, 1),
					newNumber(t, 2),
					newCall(t, vm.Call{Func: vm.FuncSum, Args: 2}),
				},
			},
		},
		{
			name: "call with one expression arg",
			ast: &ast.CallNode{
				Name: "SUM",
				Args: []ast.Node{
					&ast.BinOpNode{
						Op:    "+",
						Left:  &ast.NumberNode{Number: "1"},
						Right: &ast.NumberNode{Number: "2"},
					},
				},
			},
			want: vm.InstructionsSet{
				Instructions: []vm.Inst{
					newNumber(t, 1),
					newNumber(t, 2),
					newBinOp(t, vm.BinOpPlus),
					newCall(t, vm.Call{Func: vm.FuncSum, Args: 1}),
				},
			},
		},
		{
			name: "call with two expression args",
			ast: &ast.CallNode{
				Name: "SUM",
				Args: []ast.Node{
					&ast.BinOpNode{
						Op:    "+",
						Left:  &ast.NumberNode{Number: "1"},
						Right: &ast.NumberNode{Number: "2"},
					},
					&ast.BinOpNode{
						Op:    "-",
						Left:  &ast.NumberNode{Number: "3"},
						Right: &ast.NumberNode{Number: "4"},
					},
				},
			},
			want: vm.InstructionsSet{
				Instructions: []vm.Inst{
					newNumber(t, 1),
					newNumber(t, 2),
					newBinOp(t, vm.BinOpPlus),
					newNumber(t, 3),
					newNumber(t, 4),
					newBinOp(t, vm.BinOpMinus),
					newCall(t, vm.Call{Func: vm.FuncSum, Args: 2}),
				},
			},
		},
		{
			name: "call with range",
			ast: &ast.CallNode{
				Name: "SUM",
				Args: []ast.Node{
					&ast.RangeNode{From: "A1", To: "A3"},
				},
			},
			want: vm.InstructionsSet{
				Instructions: []vm.Inst{
					newRef(t, "A1"),
					newRef(t, "A2"),
					newRef(t, "A3"),
					newCall(t, vm.Call{Func: vm.FuncSum, Args: 3}),
				},
			},
		},
		{
			name: "call with invalid range",
			ast: &ast.CallNode{
				Name: "SUM",
				Args: []ast.Node{
					&ast.RangeNode{
						From: "A1",
						To:   "A", // must be A#
					},
				},
			},
			wantErr: vm.ErrCompile,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := vm.Compile(tc.ast)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Compile() unexpected error: %v; want %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("Compile() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func newNumber(t *testing.T, n float64) vm.Inst {
	t.Helper()
	return vm.Inst{Type: vm.InstTypeNumber, Number: n}
}

func newBinOp(t *testing.T, op vm.BinOp) vm.Inst {
	t.Helper()
	return vm.Inst{Type: vm.InstTypeBinOp, BinOp: op}
}

func newRef(t *testing.T, ref string) vm.Inst {
	t.Helper()
	return vm.Inst{Type: vm.InstTypeRef, Ref: ref}
}

func newCall(t *testing.T, c vm.Call) vm.Inst {
	t.Helper()
	return vm.Inst{Type: vm.InstTypeCall, Call: &c}
}
