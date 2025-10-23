// Copyright 2025 Samvel Khalatyan. All rights reserved.

package cycle_test

import (
	"fmt"

	"github.com/skhal/lab/iq/list/singly/cycle"
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
