// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sheet_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/go/tests"
	"github.com/skhal/lab/x/sheet/internal/sheet"
)

func TestSheet_Set(t *testing.T) {
	tt := []struct {
		name string
		text string
		want error
	}{
		{
			name: "number",
			text: "123",
		},
		{
			name: "text",
			text: "abc",
			want: sheet.ErrCell,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := sheet.New()
			cell := "A1"

			err := s.Set(cell, tc.text)

			if !errors.Is(err, tc.want) {
				t.Errorf("Set(_, %q) = %v; want %v", tc.text, err, tc.want)
			}
		})
	}
}

func TestSheet_VisitAll(t *testing.T) {
	tt := []struct {
		name  string
		cells []string
		want  []string
	}{
		{
			name: "empty",
		},
		{
			name:  "one cell",
			cells: []string{"A1"},
			want:  []string{"A1"},
		},
		{
			name:  "two cells in order",
			cells: []string{"A1", "A2"},
			want:  []string{"A1", "A2"},
		},
		{
			name:  "two cells out of order",
			cells: []string{"A2", "A1"},
			want:  []string{"A1", "A2"},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := sheet.New()
			for _, id := range tc.cells {
				s.Set(id, "123")
			}
			s.Calculate()

			var got []string
			s.VisitAll(func(id, _ string, _ float64) bool {
				got = append(got, id)
				return true
			})

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("VisitAll() mismatch (-want +got):\n%s", diff)
				t.Logf("keys: %v", tc.cells)
			}
		})
	}
}

func TestSheet_VisitAll_collectFew(t *testing.T) {
	tt := []struct {
		name  string
		cells []string
		size  int
		want  []string
	}{
		{
			name: "empty",
		},
		{
			name:  "one cell collect all",
			cells: []string{"A1"},
			size:  1,
			want:  []string{"A1"},
		},
		{
			name:  "two cells collect all",
			cells: []string{"A2", "A1"},
			size:  2,
			want:  []string{"A1", "A2"},
		},
		{
			name:  "two cells collect one",
			cells: []string{"A2", "A1"},
			size:  1,
			want:  []string{"A1"},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := sheet.New()
			for _, id := range tc.cells {
				s.Set(id, "123")
			}
			s.Calculate()

			var got []string
			s.VisitAll(func(id, _ string, _ float64) bool {
				got = append(got, id)
				return len(got) < tc.size
			})

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("VisitAll() mismatch (-want +got):\n%s", diff)
				t.Logf("keys: %v", tc.cells)
			}
		})
	}
}

func TestSheet_Calculate(t *testing.T) {
	tt := []struct {
		name    string
		cells   map[string]string
		want    map[string]float64
		wantErr error
	}{
		{
			name: "static cell",
			cells: map[string]string{
				"A1": "=123",
			},
			want: map[string]float64{
				"A1": 123,
			},
		},
		{
			name: "binary operator",
			cells: map[string]string{
				"A1": "=2 + 3",
			},
			want: map[string]float64{
				"A1": 5,
			},
		},
		{
			name: "invalid reference",
			cells: map[string]string{
				"A1": "=A2",
			},
			want: map[string]float64{
				"A1": 0,
			},
		},
		{
			name: "reference to static",
			cells: map[string]string{
				"A1": "=123",
				"B1": "=A1",
			},
			want: map[string]float64{
				"A1": 123,
				"B1": 123,
			},
		},
		{
			name: "reference to binary operator",
			cells: map[string]string{
				"A1": "=1 + 2",
				"B1": "=A1",
			},
			want: map[string]float64{
				"A1": 3,
				"B1": 3,
			},
		},
		{
			name: "multiple references calculate once",
			cells: map[string]string{
				"A1": "=1 + 2",
				"B1": "=A1",
				"B2": "=A1",
			},
			want: map[string]float64{
				"A1": 3,
				"B1": 3,
				"B2": 3,
			},
		},
		{
			name: "circular dependency in root",
			cells: map[string]string{
				"A1": "=A2",
				"A2": "=A1",
			},
			want: map[string]float64{
				"A1": 0,
				"A2": 0,
			},
			wantErr: sheet.ErrCell,
		},
		{
			name: "circular dependency in child",
			cells: map[string]string{
				"A1": "=A2",
				"A2": "=A3",
				"A3": "=A2",
			},
			want: map[string]float64{
				"A1": 0,
				"A2": 0,
				"A3": 0,
			},
			wantErr: sheet.ErrCell,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := sheet.New()
			for id, cell := range tc.cells {
				s.Set(id, cell)
			}

			err := s.Calculate()
			got := make(map[string]float64)
			s.VisitAll(func(id, _ string, n float64) bool {
				got[id] = n
				return true
			})

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Calculate() = %v; want %v", err, tc.wantErr)
			}
			opts := []cmp.Option{
				tests.EquateFloat64(0.001), // diff within 0.1%
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Errorf("Calculate() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
