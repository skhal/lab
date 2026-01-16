// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rotated_test

import (
	"fmt"

	"github.com/skhal/lab/iq/29/rotated"
)

func ExampleFind() {
	fmt.Println(rotated.Find([]int{4, 5, 1, 2, 3}, 1))
	// Output:
	// 2
}
