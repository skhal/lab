// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
