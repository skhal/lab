// Copyright 2025 Samvel Khalatyan. All rights reserved.

package singly_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/iq/list/singly"
)

func EquateLists() cmp.Option {
	return cmp.FilterValues(areLists, cmp.Comparer(compareLists))
}

func areLists(x, y interface{}) bool {
	_, okX := x.(singly.List)
	_, okY := y.(singly.List)
	return okX && okY
}

func compareLists(x, y interface{}) bool {
	listX := x.(singly.List)
	listY := y.(singly.List)
	nodeX, nodeY := listX.Head, listY.Head
	for {
		if nodeX == nil && nodeY == nil {
			break
		}
		if nodeX == nil || nodeY == nil {
			return false
		}
		if nodeX.Value != nodeY.Value {
			return false
		}
		nodeX = nodeX.Next
		nodeY = nodeY.Next
	}
	return true
}

func TestReverse(t *testing.T) {
	tests := []struct {
		name string
		list *singly.List
		want *singly.List
	}{
		{
			name: "empty",
			list: singly.NewList(),
			want: singly.NewList(),
		},
		{
			name: "one item",
			list: singly.NewList(1),
			want: singly.NewList(1),
		},
		{
			name: "two items",
			list: singly.NewList(1, 2),
			want: singly.NewList(2, 1),
		},
		{
			name: "three items",
			list: singly.NewList(1, 2, 3),
			want: singly.NewList(3, 2, 1),
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			listString := fmt.Sprint(tc.list)

			tc.list.Reverse()

			if diff := cmp.Diff(tc.want, tc.list, EquateLists()); diff != "" {
				t.Errorf("singly.List.Reverse() mismatch (-want, got):\n%s", diff)
				t.Logf("Input list:\n%s", listString)
			}
		})
	}
}
