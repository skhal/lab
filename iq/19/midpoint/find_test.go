// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package midpoint_test

import (
	"fmt"
	"testing"

	"github.com/skhal/lab/iq/19/midpoint"
)

type list midpoint.Node

func (l *list) String() string {
	var nn []int
	for l != nil {
		nn = append(nn, l.Val)
		l = (*list)(l.Next)
	}
	return fmt.Sprint(nn)
}

func makeList(nn ...int) *midpoint.Node {
	var (
		head *midpoint.Node
		tail *midpoint.Node
	)
	for _, n := range nn {
		node := &midpoint.Node{Val: n}
		switch head {
		case nil:
			head = node
		default:
			tail.Next = node
		}
		tail = node
	}
	return head
}

func TestFind(t *testing.T) {
	tests := []struct {
		name string
		list *midpoint.Node
		want int
	}{
		{
			name: "one item",
			list: makeList(1),
			want: 1,
		},
		{
			name: "two items",
			list: makeList(1, 2),
			want: 2,
		},
		{
			name: "three items",
			list: makeList(1, 2, 3),
			want: 2,
		},
		{
			name: "four items",
			list: makeList(1, 2, 3, 4),
			want: 3,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := midpoint.Find(tc.list)

			if got.Val != tc.want {
				t.Errorf("midpoint.Find(%s) = %d; want %d", (*list)(tc.list), got.Val, tc.want)
			}
		})
	}
}
