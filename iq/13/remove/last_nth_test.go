// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package remove_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/iq/13/remove"
)

func TestRemoveNth(t *testing.T) {
	tests := []struct {
		name string
		list *remove.List
		n    int
		want *remove.List
	}{
		{
			name: "empty",
			n:    1,
		},
		// size 1
		{
			name: "one item remove head",
			list: remove.NewList(1),
			n:    1,
			want: remove.NewList(),
		},
		{
			name: "one item insufficient items",
			list: remove.NewList(1),
			n:    2,
			want: remove.NewList(1),
		},
		// size 2
		{
			name: "two items remove head",
			list: remove.NewList(1, 2),
			n:    2,
			want: remove.NewList(2),
		},
		{
			name: "two items remove tail",
			list: remove.NewList(1, 2),
			n:    1,
			want: remove.NewList(1),
		},
		{
			name: "two items insufficient items",
			list: remove.NewList(1, 2),
			n:    3,
			want: remove.NewList(1, 2),
		},
		// size 3
		{
			name: "three items remove head",
			list: remove.NewList(1, 2, 3),
			n:    3,
			want: remove.NewList(2, 3),
		},
		{
			name: "three items remove tail",
			list: remove.NewList(1, 2, 3),
			n:    1,
			want: remove.NewList(1, 2),
		},
		{
			name: "three items remove first",
			list: remove.NewList(1, 2, 3),
			n:    2,
			want: remove.NewList(1, 3),
		},
		{
			name: "three items insufficient items",
			list: remove.NewList(1, 2, 3),
			n:    4,
			want: remove.NewList(1, 2, 3),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := tc.list.String()

			remove.LastNth(tc.list, tc.n)

			opts := []cmp.Option{
				cmp.Transformer("ToList", func(l *remove.List) []int {
					return l.Slice()
				}),
			}
			if diff := cmp.Diff(tc.want, tc.list, opts...); diff != "" {
				t.Errorf("remove.LastNth(%s, %d) mismatch (-want, +got):\n%s", buf, tc.n, diff)
			}
		})
	}
}
