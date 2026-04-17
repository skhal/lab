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
