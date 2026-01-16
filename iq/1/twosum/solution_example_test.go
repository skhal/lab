// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package twosum_test

import (
	"fmt"

	"github.com/skhal/lab/iq/1/twosum"
)

func Example() {
	fmt.Println(twosum.Find([]int{1, 2, 3, 4, 5}, 5))
	// Output:
	// [0 3]
}
