// Copyright 2025 Samvel Khalatyan. All rights reserved.

package palindrome_test

import (
	"fmt"

	"github.com/skhal/lab/iq/16/palindrome"
)

func ExampleIs() {
	fmt.Println(palindrome.Is(palindrome.NewList(1, 2, 1)))
	fmt.Println(palindrome.Is(palindrome.NewList(1, 2, 3)))
	// Output:
	// true
	// false
}
