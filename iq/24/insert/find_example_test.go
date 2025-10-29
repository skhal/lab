// Copyright 2025 Samvel Khalatyan. All rights reserved.

package insert_test

import (
	"fmt"

	"github.com/skhal/lab/iq/24/insert"
)

func ExampleFindInsertIndex() {
	nn := []int{1, 2, 3}
	fmt.Println(insert.FindInsertIndex(nn, 2))
	// Output:
	// 1
}
