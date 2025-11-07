// Copyright 2025 Samvel Khalatyan. All rights reserved.

package geomseq_test

import (
	"fmt"

	"github.com/skhal/lab/iq/11/geomseq"
)

func Example() {
	triplets := geomseq.Find([]int{3, 1, 2, 3, 9, 3, 27}, 3)
	for _, triplet := range triplets {
		fmt.Println(triplet)
	}
	// Output:
	// {1 3 4}
	// {0 4 6}
	// {3 4 6}
}
