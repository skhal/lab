// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package number_test

import (
	"fmt"

	"github.com/skhal/lab/iq/20/number"
)

func ExampleIsHappyNumber() {
	fmt.Println(number.IsHappyNumber(1))
	fmt.Println(number.IsHappyNumber(2))
	// Output:
	// true
	// false
}
