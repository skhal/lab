// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package container_test

import (
	"fmt"

	"github.com/skhal/lab/iq/4/container"
)

func Example() {
	fmt.Println(container.Find([]int{1, 2, 3}))
	// Output:
	// 2
}
