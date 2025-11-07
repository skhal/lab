// Copyright 2025 Samvel Khalatyan. All rights reserved.

package midpoint_test

import (
	"fmt"

	"github.com/skhal/lab/iq/19/midpoint"
)

func ExampleFind_oddNumberOfItems() {
	list := &midpoint.Node{
		Val: 1,
		Next: &midpoint.Node{
			Val: 2,
			Next: &midpoint.Node{
				Val: 3,
			},
		},
	}
	node := midpoint.Find(list)
	fmt.Println(node.Val)
	// Output:
	// 2
}

func ExampleFind_evenNumberOfItems() {
	list := &midpoint.Node{
		Val: 1,
		Next: &midpoint.Node{
			Val: 2,
			Next: &midpoint.Node{
				Val: 3,
				Next: &midpoint.Node{
					Val: 4,
				},
			},
		},
	}
	node := midpoint.Find(list)
	fmt.Println(node.Val)
	// Output:
	// 3
}
