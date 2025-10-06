// Copyright 2025 Samvel Khalatyan. All rights reserved.

package container_test

import (
	"fmt"

	"github.com/skhal/lab/iq/twopointer/container"
)

func Example() {
	fmt.Println(container.Find([]int{1, 2, 3}))
	// Output:
	// 2
}
