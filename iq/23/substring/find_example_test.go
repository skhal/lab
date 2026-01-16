// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package substring_test

import (
	"fmt"

	"github.com/skhal/lab/iq/23/substring"
)

func ExampleFind() {
	fmt.Println(substring.Find("aabcad", 2))
	// Output:
	// aabca
}
