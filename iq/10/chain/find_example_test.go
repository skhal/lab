// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chain_test

import (
	"fmt"

	"github.com/skhal/lab/iq/10/chain"
)

func Example() {
	nn := chain.Find([]int{7, 1, 8, 9, 2, 12})
	fmt.Println(nn)
	// Output:
	// [7 8 9]
}
