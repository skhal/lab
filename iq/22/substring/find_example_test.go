// Copyright 2025 Samvel Khalatyan. All rights reserved.

package substring_test

import (
	"fmt"

	"github.com/skhal/lab/iq/22/substring"
)

func ExampleFind() {
	s := substring.Find("abcad")
	fmt.Println(s)
	// Output:
	// bcad
}

func ExampleFindFast() {
	s := substring.FindFast("abcad")
	fmt.Println(s)
	// Output:
	// bcad
}
