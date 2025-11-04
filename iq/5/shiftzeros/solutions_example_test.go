// Copyright 2025 Samvel Khalatyan. All rights reserved.

package shiftzeros_test

import (
	"fmt"

	"github.com/skhal/lab/iq/5/shiftzeros"
)

func ExampleShift() {
	nn := []int{0, 1, 0, 3, 2}
	shiftzeros.Shift(nn)
	fmt.Printf("%v\n", nn)
	// Output:
	// [1 3 2 0 0]
}
