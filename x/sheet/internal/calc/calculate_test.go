// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package calc_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/x/sheet/internal/ast"
	"github.com/skhal/lab/x/sheet/internal/calc"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name    string
		node    ast.Node
		wantErr error
		want    ast.Node
	}{
		{
			name: "number node",
			node: &ast.NumberNode{Number: 123},
			want: &ast.NumberNode{Number: 123},
		},
		{
			name: "formula with number",
			node: &ast.FormulaNode{Number: &ast.NumberNode{Number: 123}},
			want: &ast.FormulaNode{Number: &ast.NumberNode{Number: 123}, Result: 123},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := calc.Calculate(tc.node)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Calculate() = %v; want %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, tc.node); diff != "" {
				t.Errorf("Calculate() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
