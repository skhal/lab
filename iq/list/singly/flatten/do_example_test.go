// Copyright 2025 Samvel Khalatyan. All rights reserved.

package flatten_test

import (
	"fmt"

	"github.com/skhal/lab/iq/list/singly/flatten"
)

func ExampleDo() {
	// Tree
	// L1  L2  L3
	// 1
	// 2 - 6 - 9
	//     7
	// 3
	// 4 - 8 - 10
	//         11
	// 5
	tree := flatten.NewTree(1, 2, 3, 4, 5)
	{
		l2 := flatten.NewTree(6, 7)
		l2.Get(6).SetChild(flatten.NewTree(9))
		tree.Get(2).SetChild(l2)
	}
	{
		l2 := flatten.NewTree(8)
		l2.Get(8).SetChild(flatten.NewTree(10, 11))
		tree.Get(4).SetChild(l2)
	}
	list := flatten.Do(tree)
	fmt.Println(list)
	// Output:
	// [1 2 3 4 5 6 7 8 9 10 11]
}
