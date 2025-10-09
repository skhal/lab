// Copyright 2025 Samvel Khalatyan. All rights reserved.

package stripzero_test

import (
	"fmt"

	"github.com/skhal/lab/iq/mapset/stripzero"
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
