// Copyright 2025 Samvel Khalatyan. All rights reserved.

package intersection_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/iq/14/intersection"
)

func TestFind(t *testing.T) {
	shared := struct {
		one *intersection.Node
		two *intersection.Node
	}{
		one: intersection.NewList(1),
		two: intersection.NewList(1, 2),
	}
	tests := []struct {
		name string
		one  *intersection.Node
		two  *intersection.Node
		want *intersection.Node
	}{
		{
			name: "empty",
		},
		{
			name: "empty and one item no intersection",
			two:  intersection.NewList(20),
		},
		{
			name: "one item and empty no intersection",
			one:  intersection.NewList(10),
		},
		// --
		{
			name: "one item no intersection",
			one:  intersection.NewList(10),
			two:  intersection.NewList(20),
		},
		{
			name: "full intersection one item",
			one:  shared.one,
			two:  shared.one,
			want: shared.one,
		},
		{
			name: "full intersection two items",
			one:  shared.two,
			two:  shared.two,
			want: shared.two,
		},
		// --
		{
			name: "same size one item no intersection",
			one:  intersection.NewList(10, 1),
			two:  intersection.NewList(20, 1),
		},
		{
			name: "same size one item intersection one item",
			one:  intersection.NewList(10).Append(shared.one),
			two:  intersection.NewList(20).Append(shared.one),
			want: shared.one,
		},
		{
			name: "same size one item intersection two items",
			one:  intersection.NewList(10).Append(shared.two),
			two:  intersection.NewList(20).Append(shared.two),
			want: shared.two,
		},
		{
			name: "same size two items no intersection",
			one:  intersection.NewList(10, 11, 1),
			two:  intersection.NewList(20, 21, 1),
		},
		{
			name: "same size two items intersection one item",
			one:  intersection.NewList(10, 11).Append(shared.one),
			two:  intersection.NewList(20, 21).Append(shared.one),
			want: shared.one,
		},
		{
			name: "same size two items intersection two items",
			one:  intersection.NewList(10).Append(shared.two),
			two:  intersection.NewList(20).Append(shared.two),
			want: shared.two,
		},
		// ---
		{
			name: "sizes two and one no intersection",
			one:  intersection.NewList(10, 11, 1),
			two:  intersection.NewList(20, 1),
		},
		{
			name: "sizes two and one intersection one item",
			one:  intersection.NewList(10, 11).Append(shared.one),
			two:  intersection.NewList(20).Append(shared.one),
			want: shared.one,
		},
		{
			name: "sizes two and one intersection two items",
			one:  intersection.NewList(10, 11).Append(shared.two),
			two:  intersection.NewList(20).Append(shared.two),
			want: shared.two,
		},
		{
			name: "sizes one and two intersection one item",
			one:  intersection.NewList(10).Append(shared.one),
			two:  intersection.NewList(20, 21).Append(shared.one),
			want: shared.one,
		},
		{
			name: "sizes one and two intersection two items",
			one:  intersection.NewList(10).Append(shared.two),
			two:  intersection.NewList(20, 21).Append(shared.two),
			want: shared.two,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := intersection.Find(tc.one, tc.two)

			opts := []cmp.Option{
				cmp.Transformer("ToList", func(node *intersection.Node) []int {
					return node.ToList()
				}),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Errorf("intersection.Find(%s, %s) mismatch (-want, +got):\n%s", tc.one, tc.two, diff)
			}
		})
	}
}
