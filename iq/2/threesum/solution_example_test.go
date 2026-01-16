// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package threesum_test

import (
	"fmt"

	"github.com/skhal/lab/iq/2/threesum"
)

func Example() {
	triplets := threesum.Find([]int{-2, -3, 5})
	for _, triplet := range triplets {
		fmt.Println(triplet)
	}
	// Output:
	// [-3 -2 5]
}
