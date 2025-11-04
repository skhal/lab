// Copyright 2025 Samvel Khalatyan. All rights reserved.

package lexiseq_test

import (
	"fmt"

	"github.com/skhal/lab/iq/6/lexiseq"
)

func Example_one() {
	s := "abcd"
	fmt.Println(lexiseq.Next(s))
	// Output:
	// abdc
}

func Example_two() {
	s := "dcba"
	fmt.Println(lexiseq.Next(s))
	// Output:
	// abcd
}
