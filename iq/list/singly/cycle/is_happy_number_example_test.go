// Copyright 2025 Samvel Khalatyan. All rights reserved.

package cycle_test

import (
	"fmt"

	"github.com/skhal/lab/iq/list/singly/cycle"
)

func ExampleIsHappyNumber() {
	fmt.Println(cycle.IsHappyNumber(1))
	fmt.Println(cycle.IsHappyNumber(2))
	// Output:
	// true
	// false
}
