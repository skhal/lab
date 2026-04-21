// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package calc_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/x/sheet/internal/calc"
)

func TestCellScanner_new(t *testing.T) {
	tests := []struct {
		name    string
		from    string
		to      string
		wantErr error
	}{
		{
			name:    "empty",
			wantErr: calc.ErrCellRange,
		},
		{
			name:    "empty from",
			to:      "A2",
			wantErr: calc.ErrCellRange,
		},
		{
			name:    "empty to",
			from:    "A1",
			wantErr: calc.ErrCellRange,
		},
		{
			name: "valid range",
			from: "A1",
			to:   "A3",
		},
		{
			name: "valid inverse range",
			from: "A3",
			to:   "A1",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := calc.NewCellScanner(tc.from, tc.to)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("NewCellScanner(%q, %q) = _, %v; want %v", tc.from, tc.to, err, tc.wantErr)
			}
		})
	}
}

func TestCellScanner_Range(t *testing.T) {
	tests := []struct {
		name string
		from string
		to   string
		want []string
	}{
		{
			name: "range",
			from: "A1",
			to:   "A3",
			want: []string{"A1", "A2", "A3"},
		},
		{
			name: "one cell range",
			from: "A1",
			to:   "A1",
			want: []string{"A1"},
		},
		{
			name: "range box",
			from: "A1",
			to:   "B3",
			want: []string{"A1", "A2", "A3", "B1", "B2", "B3"},
		},
		{
			name: "inverse range",
			from: "A3",
			to:   "A1",
			want: []string{"A1", "A2", "A3"},
		},
		{
			name: "inverse range box",
			from: "B3",
			to:   "A1",
			want: []string{"A1", "A2", "A3", "B1", "B2", "B3"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cr, err := calc.NewCellScanner(tc.from, tc.to)
			if err != nil {
				t.Fatalf("NewCellScanner(%q, %q) unexpected error: %s", tc.from, tc.to, err)
			}

			got := slices.Collect(cr.Scan())

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Range() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
