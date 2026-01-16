// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
