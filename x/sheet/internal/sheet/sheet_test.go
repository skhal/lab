// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sheet_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
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
			s.VisitAll(func(id string, _ float64) bool {
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
			s.VisitAll(func(id string, _ float64) bool {
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
