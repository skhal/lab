// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package substring_test

import (
	"fmt"

	"github.com/skhal/lab/iq/22/substring"
)

func ExampleFind() {
	s := substring.Find("abcad")
	fmt.Println(s)
	// Output:
	// bcad
}

func ExampleFindFast() {
	s := substring.FindFast("abcad")
	fmt.Println(s)
	// Output:
	// bcad
}
