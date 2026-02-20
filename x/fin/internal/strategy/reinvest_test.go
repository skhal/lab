// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strategy_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/x/fin/internal/fin"
	"github.com/skhal/lab/x/fin/internal/strategy"
)

func TestReinvestDividend(t *testing.T) {
	tests := []struct {
		name string
		pos  fin.Position
		want fin.Position
	}{
		{
			name: "zero value",
		},
		{
			name: "zero dividend",
			pos:  fin.Position{Investment: 123},
			want: fin.Position{Investment: 123},
		},
		{
			name: "nonzero dividend",
			pos:  fin.Position{Investment: 123, Dividend: 56},
			want: fin.Position{Investment: 123 + 56},
		},
		{
			name: "zero investment",
			pos:  fin.Position{Dividend: 123},
			want: fin.Position{Investment: 123},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := strategy.ReinvestDividend(tc.pos)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Reinvest(%s) mismatch (-want,+got):\n%s", tc.pos, diff)
			}
		})
	}
}
