// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package anagram_test

import (
	"fmt"

	"github.com/skhal/lab/iq/21/anagram"
)

func ExampleFindAll() {
	for _, s := range anagram.FindAll("abaabc", "aba") {
		fmt.Println(s)
	}
	// Output:
	// aba
	// baa
	// aab
}
