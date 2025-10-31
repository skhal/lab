// Copyright 2025 Samvel Khalatyan. All rights reserved.

package matrix_test

import (
	"fmt"

	"github.com/skhal/lab/iq/31/matrix"
)

func ExampleHas() {
	m := matrix.M{
		0: []int{1, 2, 4},
		1: []int{4, 5, 7},
	}
	fmt.Println(matrix.Has(m, 3))
	fmt.Println(matrix.Has(m, 5))
	// Output:
	// false
	// true
}
