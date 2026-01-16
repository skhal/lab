// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package palindrome_test

import (
	"fmt"

	"github.com/skhal/lab/iq/16/palindrome"
)

func ExampleIs() {
	fmt.Println(palindrome.Is(palindrome.NewList(1, 2, 1)))
	fmt.Println(palindrome.Is(palindrome.NewList(1, 2, 3)))
	// Output:
	// true
	// false
}
