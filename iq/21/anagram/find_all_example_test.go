// Copyright 2025 Samvel Khalatyan. All rights reserved.

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
