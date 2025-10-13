// Copyright 2025 Samvel Khalatyan. All rights reserved.

package singly_test

import (
	"fmt"

	"github.com/skhal/lab/iq/list/singly"
)

func Example() {
	list := singly.NewList(1, 2, 3)
	list.Reverse()
	fmt.Println(list)
	// Output:
	// [3 2 1]
}
