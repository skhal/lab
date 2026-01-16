// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
