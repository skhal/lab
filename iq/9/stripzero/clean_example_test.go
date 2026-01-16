// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stripzero_test

import (
	"fmt"

	"github.com/skhal/lab/iq/9/stripzero"
)

func Example() {
	matrix := [][]int{
		[]int{1, 1, 1, 1},
		[]int{1, 0, 1, 1},
		[]int{1, 1, 1, 1},
		[]int{1, 1, 1, 0},
	}
	stripzero.Clean(matrix)
	for _, nn := range matrix {
		fmt.Println(nn)
	}
	// Output:
	// [1 0 1 0]
	// [0 0 0 0]
	// [1 0 1 0]
	// [0 0 0 0]
}
