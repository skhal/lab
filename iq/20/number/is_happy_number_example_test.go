// Copyright 2025 Samvel Khalatyan. All rights reserved.

package number_test

import (
	"fmt"

	"github.com/skhal/lab/iq/20/number"
)

func ExampleIsHappyNumber() {
	fmt.Println(number.IsHappyNumber(1))
	fmt.Println(number.IsHappyNumber(2))
	// Output:
	// true
	// false
}
