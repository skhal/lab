// Copyright 2025 Samvel Khalatyan. All rights reserved.

package threesum_test

import (
	"fmt"

	"github.com/skhal/lab/iq/twopointer/threesum"
)

func Example() {
	triplets := threesum.Find([]int{-2, -3, 5})
	for _, triplet := range triplets {
		fmt.Println(triplet)
	}
	// Output:
	// {-3 -2 5}
}
