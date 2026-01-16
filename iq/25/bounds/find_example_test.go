// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bounds_test

import (
	"fmt"

	"github.com/skhal/lab/iq/25/bounds"
)

func ExampleFind_hit() {
	nn := []int{1, 2, 2, 3}
	fmt.Println(bounds.Find(nn, 2))
	// Output:
	// {1 2}
}

func ExampleFind_miss() {
	nn := []int{1, 2, 2, 3}
	fmt.Println(bounds.Find(nn, 4))
	// Output:
	// {-1 -1}
}
