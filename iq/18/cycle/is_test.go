// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cycle_test

import (
	"fmt"
	"testing"

	"github.com/skhal/lab/iq/18/cycle"
)

type testNode cycle.Node

func (node *testNode) String() string {
	var nn []int
	visits := make(map[*testNode]struct{})
	for ; node != nil; node = (*testNode)(node.Next) {
		nn = append(nn, node.Val)
		if _, ok := visits[node]; ok {
			// cycle
			break
		}
		visits[node] = struct{}{}
	}
	return fmt.Sprint(nn)
}

func TestIs_empty(t *testing.T) {
	var node *cycle.Node

	got := cycle.Is(node)

	if want := false; got != want {
		t.Errorf("cycle.Is(%s) = %v; want %v", (*testNode)(node), got, want)
	}
}

func TestIs_oneItemNoCycle(t *testing.T) {
	node := &cycle.Node{Val: 1}

	got := cycle.Is(node)

	if want := false; got != want {
		t.Errorf("cycle.Is(%s) = %v; want %v", (*testNode)(node), got, want)
	}
}

func TestIs_oneItemCycle(t *testing.T) {
	node := &cycle.Node{Val: 1}
	node.Next = node

	got := cycle.Is(node)

	if want := true; got != want {
		t.Errorf("cycle.Is(%s) = %v; want %v", (*testNode)(node), got, want)
	}
}

func TestIs_twoItemsNoCycle(t *testing.T) {
	node := &cycle.Node{
		Val: 1,
		Next: &cycle.Node{
			Val: 2,
		},
	}

	got := cycle.Is(node)

	if want := false; got != want {
		t.Errorf("cycle.Is(%s) = %v; want %v", (*testNode)(node), got, want)
	}
}

func TestIs_twoItemsCycleToFirst(t *testing.T) {
	node := &cycle.Node{
		Val: 1,
		Next: &cycle.Node{
			Val: 2,
		},
	}
	node.Next.Next = node

	got := cycle.Is(node)

	if want := true; got != want {
		t.Errorf("cycle.Is(%s) = %v; want %v", (*testNode)(node), got, want)
	}
}

func TestIs_twoItemsCycleToSecond(t *testing.T) {
	node := &cycle.Node{
		Val: 1,
		Next: &cycle.Node{
			Val: 2,
		},
	}
	node.Next.Next = node.Next

	got := cycle.Is(node)

	if want := true; got != want {
		t.Errorf("cycle.Is(%s) = %v; want %v", (*testNode)(node), got, want)
	}
}
