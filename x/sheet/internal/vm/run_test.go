// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vm_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/skhal/lab/x/sheet/internal/vm"
)

func TestRun_number(t *testing.T) {
	tests := []struct {
		name    string
		iset    []any
		want    float64
		wantErr error
	}{
		{
			name:    "empty",
			wantErr: vm.ErrRun,
		},
		{
			name: "number",
			iset: []any{
				&vm.Number{Val: 1.23},
			},
			want: 1.23,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			iset := &vm.InstructionsSet{Instructions: tc.iset}
			var refcal func(string) (float64, error)

			got, err := vm.Run(iset, refcal)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Run() unexpected error: %v; want %v", err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("Run() unexpected result %f; want %f", got, tc.want)
			}
		})
	}
}

func TestRun_operator(t *testing.T) {
	tests := []struct {
		name    string
		iset    []any
		want    float64
		wantErr error
	}{
		{
			name: "operator plus",
			iset: []any{
				&vm.Number{Val: 1},
				&vm.Number{Val: 2},
				&vm.BinOp{Op: vm.OpPlus},
			},
			want: 3,
		},
		{
			name: "operator minus",
			iset: []any{
				&vm.Number{Val: 1},
				&vm.Number{Val: 2},
				&vm.BinOp{Op: vm.OpMinus},
			},
			want: -1,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			iset := &vm.InstructionsSet{Instructions: tc.iset}
			var refcal func(string) (float64, error)

			got, err := vm.Run(iset, refcal)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Run() unexpected error: %v; want %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got, cmpopts.EquateApprox(0.001, 0)); diff != "" {
				t.Errorf("Run() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestRun_reference(t *testing.T) {
	tests := []struct {
		name    string
		iset    []any
		refs    map[string]float64
		want    float64
		wantErr error
	}{
		{
			name: "reference",
			iset: []any{
				&vm.Ref{Cell: "A1"},
			},
			refs: map[string]float64{
				"A1": 1,
			},
			want: 1,
		},
		{
			name: "missing reference",
			iset: []any{
				&vm.Ref{Cell: "A1"},
			},
			wantErr: vm.ErrRun,
		},
		{
			name: "operator plus left operand reference",
			iset: []any{
				&vm.Ref{Cell: "A1"},
				&vm.Number{Val: 2},
				&vm.BinOp{Op: vm.OpPlus},
			},
			refs: map[string]float64{
				"A1": 1,
			},
			want: 3,
		},
		{
			name: "operator plus right operand reference",
			iset: []any{
				&vm.Number{Val: 1},
				&vm.Ref{Cell: "A2"},
				&vm.BinOp{Op: vm.OpPlus},
			},
			refs: map[string]float64{
				"A2": 2,
			},
			want: 3,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			iset := &vm.InstructionsSet{Instructions: tc.iset}
			refcal := newTestRefCalculator(t, tc.refs)

			got, err := vm.Run(iset, refcal)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Run() unexpected error: %v; want %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got, cmpopts.EquateApprox(0.001, 0)); diff != "" {
				t.Errorf("Run() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestRun_call(t *testing.T) {
	tests := []struct {
		name    string
		iset    []any
		refs    map[string]float64
		want    float64
		wantErr error
	}{
		{
			name: "sum no args",
			iset: []any{
				&vm.Call{Func: vm.FuncSum},
			},
			want: 0,
		},
		{
			name: "sum one arg",
			iset: []any{
				&vm.Number{Val: 1},
				&vm.Call{Func: vm.FuncSum, Args: 1},
			},
			want: 1,
		},
		{
			name: "sum two args",
			iset: []any{
				&vm.Number{Val: 1},
				&vm.Number{Val: 2},
				&vm.Call{Func: vm.FuncSum, Args: 2},
			},
			want: 3,
		},
		{
			name: "sum one reference",
			iset: []any{
				&vm.Ref{Cell: "A1"},
				&vm.Call{Func: vm.FuncSum, Args: 1},
			},
			refs: map[string]float64{
				"A1": 1,
			},
			want: 1,
		},
		{
			name: "sum multiple refs",
			iset: []any{
				&vm.Ref{Cell: "A1"},
				&vm.Ref{Cell: "A2"},
				&vm.Call{Func: vm.FuncSum, Args: 2},
			},
			refs: map[string]float64{
				"A1": 1,
				"A2": 2,
			},
			want: 3,
		},
		{
			name: "sum reference and number",
			iset: []any{
				&vm.Ref{Cell: "A1"},
				&vm.Number{Val: 2},
				&vm.Call{Func: vm.FuncSum, Args: 2},
			},
			refs: map[string]float64{
				"A1": 1,
			},
			want: 3,
		},
		{
			name: "sum number and a reference",
			iset: []any{
				&vm.Number{Val: 2},
				&vm.Ref{Cell: "A1"},
				&vm.Call{Func: vm.FuncSum, Args: 2},
			},
			refs: map[string]float64{
				"A1": 1,
			},
			want: 3,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			iset := &vm.InstructionsSet{Instructions: tc.iset}
			refcal := newTestRefCalculator(t, tc.refs)

			got, err := vm.Run(iset, refcal)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Run() unexpected error: %v; want %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got, cmpopts.EquateApprox(0.001, 0)); diff != "" {
				t.Errorf("Run() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

type testRefCalculator struct {
	refs map[string]float64
}

func newTestRefCalculator(t *testing.T, refs map[string]float64) func(string) (float64, error) {
	t.Helper()
	refcal := &testRefCalculator{refs}
	return refcal.Calculate
}

func (rc *testRefCalculator) Calculate(id string) (float64, error) {
	n, ok := rc.refs[id]
	if !ok {
		return 0, fmt.Errorf("ref %s does not exist", id)
	}
	return n, nil
}
