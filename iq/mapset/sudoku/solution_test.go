// Copyright 2025 Samvel Khalatyan. All rights reserved.

package sudoku_test

import (
	"testing"

	"github.com/skhal/lab/iq/mapset/sudoku"
)

func TestSudoku_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		board *sudoku.Board
		want  bool
	}{
		{
			name:  "empty",
			board: sudoku.NewBoard(),
			want:  true,
		},
		// rows
		{
			name: "valid row 1",
			board: sudoku.NewBoard(
				sudoku.Box{sudoku.C1, 1},
				sudoku.Box{sudoku.A1, 2},
				sudoku.Box{sudoku.G1, 3},
			),
			want: true,
		},
		{
			name: "invalid row 1",
			board: sudoku.NewBoard(
				sudoku.Box{sudoku.C1, 1},
				sudoku.Box{sudoku.A1, 2},
				sudoku.Box{sudoku.G1, 1},
			),
		},
		{
			name: "valid row 7",
			board: sudoku.NewBoard(
				sudoku.Box{sudoku.A7, 9},
				sudoku.Box{sudoku.F7, 3},
				sudoku.Box{sudoku.I7, 8},
			),
			want: true,
		},
		{
			name: "invalid row 7",
			board: sudoku.NewBoard(
				sudoku.Box{sudoku.A7, 9},
				sudoku.Box{sudoku.F7, 3},
				sudoku.Box{sudoku.I7, 3},
			),
		},
		// columns
		{
			name: "valid col c",
			board: sudoku.NewBoard(
				sudoku.Box{sudoku.C1, 1},
				sudoku.Box{sudoku.C4, 2},
				sudoku.Box{sudoku.C9, 3},
			),
			want: true,
		},
		{
			name: "invalid col c",
			board: sudoku.NewBoard(
				sudoku.Box{sudoku.C1, 1},
				sudoku.Box{sudoku.C4, 2},
				sudoku.Box{sudoku.C9, 2},
			),
		},
		// blocks
		{
			name: "valid block 1",
			board: sudoku.NewBoard(
				sudoku.Box{sudoku.A1, 1},
				sudoku.Box{sudoku.A2, 2},
				sudoku.Box{sudoku.A3, 3},
				sudoku.Box{sudoku.B1, 4},
				sudoku.Box{sudoku.B2, 5},
				sudoku.Box{sudoku.B3, 6},
				sudoku.Box{sudoku.C1, 7},
				sudoku.Box{sudoku.C2, 8},
				sudoku.Box{sudoku.C3, 9},
			),
			want: true,
		},
		{
			name: "invalid block 1",
			board: sudoku.NewBoard(
				sudoku.Box{sudoku.A1, 1},
				sudoku.Box{sudoku.A2, 2},
				sudoku.Box{sudoku.A3, 3},
				sudoku.Box{sudoku.B1, 4},
				sudoku.Box{sudoku.B2, 5},
				sudoku.Box{sudoku.B3, 4},
				sudoku.Box{sudoku.C1, 7},
				sudoku.Box{sudoku.C2, 8},
				sudoku.Box{sudoku.C3, 9},
			),
		},
		{
			name: "valid block 2",
			board: sudoku.NewBoard(
				sudoku.Box{sudoku.D1, 1},
				sudoku.Box{sudoku.D2, 2},
				sudoku.Box{sudoku.D3, 3},
				sudoku.Box{sudoku.E1, 4},
				sudoku.Box{sudoku.E2, 5},
				sudoku.Box{sudoku.E3, 6},
				sudoku.Box{sudoku.F1, 7},
				sudoku.Box{sudoku.F2, 8},
				sudoku.Box{sudoku.F3, 9},
			),
			want: true,
		},
		{
			name: "invalid block 2",
			board: sudoku.NewBoard(
				sudoku.Box{sudoku.D1, 1},
				sudoku.Box{sudoku.D2, 2},
				sudoku.Box{sudoku.D3, 3},
				sudoku.Box{sudoku.E1, 4},
				sudoku.Box{sudoku.E2, 5},
				sudoku.Box{sudoku.E3, 4},
				sudoku.Box{sudoku.F1, 7},
				sudoku.Box{sudoku.F2, 8},
				sudoku.Box{sudoku.F3, 9},
			),
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := tc.board.IsValid()

			if got != tc.want {
				t.Errorf("board.IsValid() = %v; want %v\nBoard:\n%v", got, tc.want, tc.board)
			}
		})
	}
}
