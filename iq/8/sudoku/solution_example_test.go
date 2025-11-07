// Copyright 2025 Samvel Khalatyan. All rights reserved.

package sudoku_test

import (
	"fmt"

	"github.com/skhal/lab/iq/8/sudoku"
)

func Example() {
	b := sudoku.NewBoard(
		sudoku.Box{sudoku.A1, 1},
		sudoku.Box{sudoku.A5, 1},
	)
	fmt.Println(b.IsValid())
	// Output:
	// false
}
