// Copyright 2025 Samvel Khalatyan. All rights reserved.

package twosum_test

import (
	"fmt"

	"github.com/skhal/lab/iq/twopointer/twosum"
)

func Example() {
	fmt.Println(twosum.Find([]int{1, 2, 3, 4, 5}, 5))
	// Output:
	// [0 3]
}
