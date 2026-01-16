// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reverse_test

import (
	"fmt"

	"github.com/skhal/lab/iq/12/reverse"
)

func Example() {
	list := reverse.NewList(1, 2, 3)
	list.Reverse()
	fmt.Println(list)
	// Output:
	// [3 2 1]
}
