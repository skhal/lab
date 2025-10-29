// Copyright 2025 Samvel Khalatyan. All rights reserved.

package lower_test

import (
	"fmt"

	"github.com/skhal/lab/iq/27/lower"
)

func ExampleFind() {
	nn := []int{1, 3}
	n, _ := lower.Find(nn, 0)
	fmt.Println(n)
	n, _ = lower.Find(nn, 1)
	fmt.Println(n)
	// Output:
	// 1
	// 1
}
