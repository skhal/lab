// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package geomseq_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/iq/11/geomseq"
)

func TestFind(t *testing.T) {
	tests := []struct {
		name string
		nn   []int
		r    geomseq.Ratio
		want []geomseq.Triplet
	}{
		{
			name: "empty",
			r:    2,
		},
		{
			name: "one item",
			nn:   []int{1},
			r:    2,
		},
		{
			name: "two items",
			nn:   []int{1, 2},
			r:    2,
		},
		{
			name: "thee items no seq",
			nn:   []int{1, 2, 3},
			r:    2,
		},
		{
			name: "thee items seq",
			nn:   []int{1, 2, 4},
			r:    2,
			want: []geomseq.Triplet{
				{0, 1, 2},
			},
		},
		{
			name: "four items no seq",
			nn:   []int{1, 2, 3, 5},
			r:    2,
		},
		{
			name: "four items two seq",
			nn:   []int{1, 2, 2, 4},
			r:    2,
			want: []geomseq.Triplet{
				{0, 1, 3},
				{0, 2, 3},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := geomseq.Find(tc.nn, tc.r)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("geomseq.Find(%v, %d) mismatch (-want, +got):\n%s", tc.nn, tc.r, diff)
			}
		})
	}
}
