// Copyright 2025 Samvel Khalatyan. All rights reserved.

package intersection_test

import (
	"fmt"

	"github.com/skhal/lab/iq/14/intersection"
)

func ExampleFind() {
	shared := intersection.NewList(1, 2)
	l1 := intersection.NewList(3, 4).Append(shared)
	l2 := intersection.NewList(5).Append(shared)
	node := intersection.Find(l1, l2)
	fmt.Println(node)
	// Output:
	// [1 2]
}
