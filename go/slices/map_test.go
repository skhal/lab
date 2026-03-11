// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/go/slices"
)

func ExampleMapFunc() {
	nn := slices.MapFunc([]int{1, 2, 3}, func(n int) int { return n * 2 })
	fmt.Println(nn)
	// Output:
	// [2 4 6]
}

func TestMapFunc(t *testing.T) {
	mapfn := strconv.Itoa
	tests := []struct {
		name string
		s    []int
		want []string
	}{
		{
			name: "nil",
		},
		{
			name: "empty",
			s:    []int{},
			want: []string{},
		},
		{
			name: "one item",
			s:    []int{1},
			want: []string{"1"},
		},
		{
			name: "two items",
			s:    []int{1, 2},
			want: []string{"1", "2"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := slices.MapFunc(tc.s, mapfn)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("MapFunc() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
