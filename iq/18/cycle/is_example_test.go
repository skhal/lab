// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cycle_test

import (
	"fmt"

	"github.com/skhal/lab/iq/18/cycle"
)

func ExampleIs() {
	head := &cycle.Node{
		Val: 1,
		Next: &cycle.Node{
			Val: 2,
			Next: &cycle.Node{
				Val: 3,
			},
		},
	}
	head.Next.Next = head
	fmt.Println(cycle.Is(head))
	// Output:
	// true
}
