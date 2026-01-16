// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lexiseq_test

import (
	"fmt"

	"github.com/skhal/lab/iq/6/lexiseq"
)

func Example_one() {
	s := "abcd"
	fmt.Println(lexiseq.Next(s))
	// Output:
	// abdc
}

func Example_two() {
	s := "dcba"
	fmt.Println(lexiseq.Next(s))
	// Output:
	// abcd
}
