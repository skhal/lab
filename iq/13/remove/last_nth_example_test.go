// Copyright 2025 Samvel Khalatyan. All rights reserved.

package remove_test

import (
	"fmt"

	"github.com/skhal/lab/iq/13/remove"
)

func ExampleLastNth() {
	l := remove.NewList(1, 2, 3, 4, 5)
	remove.LastNth(l, 3)
	fmt.Println(l)
	// Output:
	// [1 2 4 5]
}
