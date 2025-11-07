// Copyright 2025 Samvel Khalatyan. All rights reserved.

package chain_test

import (
	"fmt"

	"github.com/skhal/lab/iq/10/chain"
)

func Example() {
	nn := chain.Find([]int{7, 1, 8, 9, 2, 12})
	fmt.Println(nn)
	// Output:
	// [7 8 9]
}
