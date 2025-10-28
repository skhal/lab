// Copyright 2025 Samvel Khalatyan. All rights reserved.

package substring_test

import (
	"fmt"

	"github.com/skhal/lab/iq/23/substring"
)

func ExampleFind() {
	fmt.Println(substring.Find("aabcad", 2))
	// Output:
	// aabca
}
