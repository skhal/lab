// Copyright 2025 Samvel Khalatyan. All rights reserved.

package chain_test

import (
	"fmt"

	"github.com/skhal/lab/iq/mapset/chain"
)

func Example() {
	nn := chain.Find([]int{7, 1, 8, 9, 2, 12})
	fmt.Println(nn)
	// Output:
	// [7 8 9]
}
