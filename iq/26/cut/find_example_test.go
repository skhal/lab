// Copyright 2025 Samvel Khalatyan. All rights reserved.

package cut_test

import (
	"fmt"

	"github.com/skhal/lab/iq/26/cut"
)

func ExampleFind() {
	fmt.Println(cut.Find([]int{1, 3, 2, 4}, 2))
	// Output:
	// 2
}
