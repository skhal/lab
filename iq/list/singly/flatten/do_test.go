// Copyright 2025 Samvel Khalatyan. All rights reserved.

package flatten_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/iq/list/singly/flatten"
)

var toSlice = cmp.Transformer("ToSlice", func(n *flatten.Node) []int {
	return n.Slice()
})

func TestDo_empty(t *testing.T) {
	tree := flatten.NewTree()
	want := flatten.NewList()

	got := flatten.Do(tree)

	if diff := cmp.Diff(want, got, toSlice); diff != "" {
		t.Errorf("flatten.Do() mismatch (-want, +got):\n%s", diff)
		t.Logf("tree:\n%v", tree)
	}
}

func TestDo_flatOneItem(t *testing.T) {
	tree := flatten.NewTree(1)
	want := flatten.NewList(1)

	got := flatten.Do(tree)

	if diff := cmp.Diff(want, got, toSlice); diff != "" {
		t.Errorf("flatten.Do() mismatch (-want, +got):\n%s", diff)
		t.Logf("tree:\n%v", tree)
	}
}

func TestDo_flatTwoItems(t *testing.T) {
	tree := flatten.NewTree(1, 2)
	want := flatten.NewList(1, 2)

	got := flatten.Do(tree)

	if diff := cmp.Diff(want, got, toSlice); diff != "" {
		t.Errorf("flatten.Do() mismatch (-want, +got):\n%s", diff)
		t.Logf("tree:\n%v", tree)
	}
}

func TestDo_flatThreeItems(t *testing.T) {
	tree := flatten.NewTree(1, 2, 3)
	want := flatten.NewList(1, 2, 3)

	got := flatten.Do(tree)

	if diff := cmp.Diff(want, got, toSlice); diff != "" {
		t.Errorf("flatten.Do() mismatch (-want, +got):\n%s", diff)
		t.Logf("tree:\n%v", tree)
	}
}

func TestDo_l1OneL2One(t *testing.T) {
	tree := flatten.NewTree(1)
	tree.Get(1).SetChild(flatten.NewTree(2))
	want := flatten.NewList(1, 2)

	got := flatten.Do(tree)

	if diff := cmp.Diff(want, got, toSlice); diff != "" {
		t.Errorf("flatten.Do() mismatch (-want, +got):\n%s", diff)
		t.Logf("tree:\n%v", tree)
	}
}

func TestDo_l1TwoL2OneOnFirstL1(t *testing.T) {
	tree := flatten.NewTree(1, 2)
	tree.Get(1).SetChild(flatten.NewTree(3))
	want := flatten.NewList(1, 2, 3)

	got := flatten.Do(tree)

	if diff := cmp.Diff(want, got, toSlice); diff != "" {
		t.Errorf("flatten.Do() mismatch (-want, +got):\n%s", diff)
		t.Logf("tree:\n%v", tree)
	}
}

func TestDo_l1TwoL2OneOnSecondL1(t *testing.T) {
	tree := flatten.NewTree(1, 2)
	tree.Get(2).SetChild(flatten.NewTree(3))
	want := flatten.NewList(1, 2, 3)

	got := flatten.Do(tree)

	if diff := cmp.Diff(want, got, toSlice); diff != "" {
		t.Errorf("flatten.Do() mismatch (-want, +got):\n%s", diff)
		t.Logf("tree:\n%v", tree)
	}
}

func TestDo_l1TwoL2TwoOnFirstL1(t *testing.T) {
	tree := flatten.NewTree(1, 2)
	tree.Get(1).SetChild(flatten.NewTree(3, 4))
	want := flatten.NewList(1, 2, 3, 4)

	got := flatten.Do(tree)

	if diff := cmp.Diff(want, got, toSlice); diff != "" {
		t.Errorf("flatten.Do() mismatch (-want, +got):\n%s", diff)
		t.Logf("tree:\n%v", tree)
	}
}

func TestDo_l1TwoL2TwoOnSecondL1(t *testing.T) {
	tree := flatten.NewTree(1, 2)
	tree.Get(2).SetChild(flatten.NewTree(3, 4))
	want := flatten.NewList(1, 2, 3, 4)

	got := flatten.Do(tree)

	if diff := cmp.Diff(want, got, toSlice); diff != "" {
		t.Errorf("flatten.Do() mismatch (-want, +got):\n%s", diff)
		t.Logf("tree:\n%v", tree)
	}
}

func TestDo_l1TwoL2OnePerL1Item(t *testing.T) {
	tree := flatten.NewTree(1, 2)
	tree.Get(1).SetChild(flatten.NewTree(3))
	tree.Get(2).SetChild(flatten.NewTree(4))
	want := flatten.NewList(1, 2, 3, 4)

	got := flatten.Do(tree)

	if diff := cmp.Diff(want, got, toSlice); diff != "" {
		t.Errorf("flatten.Do() mismatch (-want, +got):\n%s", diff)
		t.Logf("tree:\n%v", tree)
	}
}
