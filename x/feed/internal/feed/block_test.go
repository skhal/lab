// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package feed_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/x/feed/internal/feed"
)

func ExampleEqualSizeBlocks() {
	chunks := slices.Collect(feed.EqualSizeBlocks(3, 7))
	fmt.Println(chunks)
	// Output:
	// [{0 3} {3 6} {6 7}]
}

func TestEqualSizeBlocks(t *testing.T) {
	tests := []struct {
		name  string
		size  int
		limit int
		want  []feed.Block
	}{
		{
			name: "empty",
		},
		{
			name:  "size less than limit",
			size:  2,
			limit: 3,
			want:  []feed.Block{{0, 2}, {2, 3}},
		},
		{
			name:  "size equal to limit",
			size:  3,
			limit: 3,
			want:  []feed.Block{{0, 3}},
		},
		{
			name:  "size larger then limit",
			size:  4,
			limit: 3,
			want:  []feed.Block{{0, 3}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			chunks := slices.Collect(feed.EqualSizeBlocks(tc.size, tc.limit))

			if diff := cmp.Diff(tc.want, chunks); diff != "" {
				t.Errorf("feed.EqualSizeBlocks(%d, %d) mismatch (-want,+got):\n%s", tc.size, tc.limit, diff)
			}
		})
	}
}
