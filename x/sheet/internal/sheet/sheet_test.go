// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sheet_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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

			err := s.Set("A1", tc.text)

			if !errors.Is(err, tc.want) {
				t.Errorf("Set(_, %q) = %v; want %v", tc.text, err, tc.want)
			}
		})
	}
}

func TestSheet_VisitAll(t *testing.T) {
	tt := []struct {
		name  string
		cells map[string]string
		want  map[string]string
	}{
		{
			name: "empty",
		},
		{
			name: "one cell",
			cells: map[string]string{
				"A1": "123",
			},
			want: map[string]string{
				"A1": "123",
			},
		},
		{
			name: "two cells",
			cells: map[string]string{
				"A1": "1",
				"A2": "2",
			},
			want: map[string]string{
				"A1": "1",
				"A2": "2",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := newSheet(t, tc.cells)
			s.Calculate()

			cells := collectCells(t, s)

			if diff := cmp.Diff(tc.want, cells, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("VisitAll() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSheet_VisitAll_collectFew(t *testing.T) {
	tt := []struct {
		name  string
		cells map[string]string
		size  int
		want  map[string]string
	}{
		{
			name: "empty",
		},
		{
			name: "one cell collect all",
			cells: map[string]string{
				"A1": "1",
			},
			size: 1,
			want: map[string]string{
				"A1": "1",
			},
		},
		{
			name: "two cells collect one",
			cells: map[string]string{
				"A1": "1",
				"A2": "2",
			},
			size: 1,
			want: map[string]string{
				"A1": "1",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := newSheet(t, tc.cells)
			s.Calculate()

			cells := make(map[string]string)
			s.VisitAll(func(id, val string, _ float64) bool {
				cells[id] = val
				return len(cells) < tc.size
			})

			if diff := cmp.Diff(tc.want, cells, cmpopts.EquateEmpty()); diff != "" {
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
			s := newSheet(t, tc.cells)

			err := s.Calculate()
			cells := collectCellResults(t, s)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Calculate() = %v; want %v", err, tc.wantErr)
			}
			opts := []cmp.Option{
				cmpopts.EquateEmpty(),
				tests.EquateFloat64(0.001), // diff within 0.1%
			}
			if diff := cmp.Diff(tc.want, cells, opts...); diff != "" {
				t.Errorf("Calculate() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSheet_Write(t *testing.T) {
	tests := []struct {
		name    string
		cells   map[string]string
		wantErr error
	}{
		{
			name: "empty",
		},
		{
			name: "one cell",
			cells: map[string]string{
				"A1": "123",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newSheet(t, tc.cells)
			var buf bytes.Buffer

			err := s.Write(&buf)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Write() = %v; want %v", err, tc.wantErr)
			}
		})
	}
}

func TestSheet_Read(t *testing.T) {
	tests := []struct {
		name    string
		cells   map[string]string
		wantErr error
		want    map[string]string
	}{
		{
			name: "empty",
		},
		{
			name: "one cell",
			cells: map[string]string{
				"A1": "123",
			},
			want: map[string]string{
				"A1": "123",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := writeSheet(t, tc.cells)
			s := sheet.New()

			err := s.Read(bytes.NewReader(b))
			cells := collectCells(t, s)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Read() = %v; want %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, cells, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("Read() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSheet_Read_resetsSheet(t *testing.T) {
	b := writeSheet(t, map[string]string{
		"A1": "123",
	})
	s := newSheet(t, map[string]string{
		"B1": "567",
	})
	want := map[string]string{
		"A1": "123",
	}

	err := s.Read(bytes.NewReader(b))
	cells := collectCells(t, s)

	if err != nil {
		t.Errorf("Read() unexpected error %v", err)
	}
	if diff := cmp.Diff(want, cells, cmpopts.EquateEmpty()); diff != "" {
		t.Errorf("Read() mismatch (-want +got):\n%s", diff)
	}
}

func BenchmarkSheet(b *testing.B) {
	for b.Loop() {
		s := sheet.New()
		s.Set("A1", "1")
		s.Set("A2", "2")
		s.Set("A3", "3")
		s.Set("A4", "4")
		s.Set("A5", "5")
		s.Set("B1", "=SUM(A1:A5, 6-7)")
		s.Calculate()
	}
}

func ExampleSheet() {
	s := sheet.New()
	// ignore-error start
	s.Set("A1", "1")
	s.Set("B1", "=SUM(A1:A5, 7-6)")
	s.Calculate()
	// ignore-error end
	s.VisitAll(func(id, val string, res float64) bool {
		fmt.Printf("%s %3.1f\t%s\n", id, res, val)
		return true
	})
	// Output:
	// A1 1.0	1
	// B1 2.0	=SUM(A1:A5, 7-6)
}

func ExampleSheet_writeRead() {
	b := func() []byte {
		s := sheet.New()
		s.Set("A1", "1")
		s.Set("B1", "=SUM(A1:A5, 7-6)")
		var b bytes.Buffer
		if err := s.Write(&b); err != nil {
			fmt.Println(err)
			return nil
		}
		return b.Bytes()
	}()
	s := sheet.New()
	if err := s.Read(bytes.NewReader(b)); err != nil {
		fmt.Println(err)
		return
	}
	if err := s.Calculate(); err != nil {
		fmt.Println(err)
		return
	}
	s.VisitAll(func(id, val string, res float64) bool {
		fmt.Printf("%s %3.1f\t%s\n", id, res, val)
		return true
	})
	// Output:
	// A1 1.0	1
	// B1 2.0	=SUM(A1:A5, 7-6)
}

func newSheet(t *testing.T, cells map[string]string) *sheet.Sheet {
	t.Helper()
	s := sheet.New()
	for id, val := range cells {
		if err := s.Set(id, val); err != nil {
			t.Fatalf("Set() unexpected error %v", err)
		}
	}
	return s
}

func writeSheet(t *testing.T, cells map[string]string) []byte {
	t.Helper()
	s := newSheet(t, cells)
	var buf bytes.Buffer
	if err := s.Write(&buf); err != nil {
		t.Fatalf("Write() unexpected error %v", err)
	}
	return buf.Bytes()
}

func collectCells(t *testing.T, s *sheet.Sheet) map[string]string {
	t.Helper()
	cells := make(map[string]string)
	s.VisitAll(func(id, val string, _ float64) bool {
		cells[id] = val
		return true
	})
	return cells
}

func collectCellResults(t *testing.T, s *sheet.Sheet) map[string]float64 {
	t.Helper()
	cells := make(map[string]float64)
	s.VisitAll(func(id, _ string, n float64) bool {
		cells[id] = n
		return true
	})
	return cells
}
