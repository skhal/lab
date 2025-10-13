// Copyright 2025 Samvel Khalatyan. All rights reserved.

package reverse_test

import (
	"fmt"

	"github.com/skhal/lab/iq/list/singly/reverse"
)

func Example() {
	list := reverse.NewList(1, 2, 3)
	list.Reverse()
	fmt.Println(list)
	// Output:
	// [3 2 1]
}
