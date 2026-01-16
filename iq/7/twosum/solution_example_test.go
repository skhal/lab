// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package twosum_test

import (
	"fmt"
	"sort"

	"github.com/skhal/lab/iq/7/twosum"
)

func Example() {
	indices := twosum.Find([]int{1, 2, 3}, 4)
	sort.Ints(indices)
	fmt.Println(indices)
	// Output:
	// [0 2]
}
