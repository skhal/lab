// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package matrix_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/skhal/lab/iq/31/matrix"
)

type test struct {
	name string
	m    matrix.M
	n    int
	want bool
}

func TestMatrix(t *testing.T) {
	tests := []struct {
		name  string
		tests []test
	}{
		{
			name: "empty",
			tests: []test{
				{
					name: "empty matrix",
					n:    1,
				},
			},
		},
		{
			name: "row",
			tests: []test{
				// size 1
				{
					name: "size 1 hit",
					m: matrix.M{
						0: []int{2},
					},
					n:    2,
					want: true,
				},
				{
					name: "size 1 miss below",
					m: matrix.M{
						0: []int{2},
					},
					n: 1,
				},
				{
					name: "size 1 miss above",
					m: matrix.M{
						0: []int{2},
					},
					n: 3,
				},
				// size 2
				{
					name: "size 2 hit first",
					m: matrix.M{
						0: []int{2, 4},
					},
					n:    2,
					want: true,
				},
				{
					name: "size 2 hit second",
					m: matrix.M{
						0: []int{2, 4},
					},
					n:    2,
					want: true,
				},
				{
					name: "size 2 miss below",
					m: matrix.M{
						0: []int{2, 4},
					},
					n: 1,
				},
				{
					name: "size 2 miss between",
					m: matrix.M{
						0: []int{2, 4},
					},
					n: 3,
				},
				{
					name: "size 2 miss above",
					m: matrix.M{
						0: []int{2, 4},
					},
					n: 5,
				},
				// size 3
				{
					name: "size 3 hit first",
					m: matrix.M{
						0: []int{2, 4, 6},
					},
					n:    2,
					want: true,
				},
				{
					name: "size 3 hit second",
					m: matrix.M{
						0: []int{2, 4, 6},
					},
					n:    2,
					want: true,
				},
				{
					name: "size 3 hit third",
					m: matrix.M{
						0: []int{2, 4, 6},
					},
					n:    6,
					want: true,
				},
				{
					name: "size 3 miss below",
					m: matrix.M{
						0: []int{2, 4, 6},
					},
					n: 1,
				},
				{
					name: "size 3 miss between first and second",
					m: matrix.M{
						0: []int{2, 4, 6},
					},
					n: 3,
				},
				{
					name: "size 3 miss between second and third",
					m: matrix.M{
						0: []int{2, 4, 6},
					},
					n: 5,
				},
				{
					name: "size 3 miss above",
					m: matrix.M{
						0: []int{2, 4, 6},
					},
					n: 7,
				},
			},
		},
		{
			name: "column",
			tests: []test{
				// size 1 is covered in the row tests
				// size 2
				{
					name: "size 2 miss below",
					m: matrix.M{
						0: []int{2},
						1: []int{4},
					},
					n: 1,
				},
				{
					name: "size 2 hit first",
					m: matrix.M{
						0: []int{2},
						1: []int{4},
					},
					n:    2,
					want: true,
				},
				{
					name: "size 2 miss between",
					m: matrix.M{
						0: []int{2},
						1: []int{4},
					},
					n: 3,
				},
				{
					name: "size 2 hit second",
					m: matrix.M{
						0: []int{2},
						1: []int{4},
					},
					n:    4,
					want: true,
				},
				{
					name: "size 2 miss above",
					m: matrix.M{
						0: []int{2},
						1: []int{4},
					},
					n: 3,
				},
				// size 3
				{
					name: "size 3 miss below",
					m: matrix.M{
						0: []int{2},
						1: []int{4},
						2: []int{6},
					},
					n: 1,
				},
				{
					name: "size 3 hit first",
					m: matrix.M{
						0: []int{2},
						1: []int{4},
						2: []int{6},
					},
					n:    2,
					want: true,
				},
				{
					name: "size 3 miss between first and second",
					m: matrix.M{
						0: []int{2},
						1: []int{4},
						2: []int{6},
					},
					n: 3,
				},
				{
					name: "size 3 hit second",
					m: matrix.M{
						0: []int{2},
						1: []int{4},
						2: []int{6},
					},
					n:    4,
					want: true,
				},
				{
					name: "size 3 miss between second and third",
					m: matrix.M{
						0: []int{2},
						1: []int{4},
						2: []int{6},
					},
					n: 5,
				},
				{
					name: "size 3 hit third",
					m: matrix.M{
						0: []int{2},
						1: []int{4},
						2: []int{6},
					},
					n:    6,
					want: true,
				},
				{
					name: "size 3 miss above",
					m: matrix.M{
						0: []int{2},
						1: []int{4},
						2: []int{6},
					},
					n: 7,
				},
			},
		},
		{
			name: "two rows",
			tests: []test{
				// col 1 is covered in the column tests
				// col 2
				{
					name: "two cols miss below",
					m: matrix.M{
						0: []int{2, 4},
						1: []int{6, 8},
					},
					n: 1,
				},
				{
					name: "two cols hit first",
					m: matrix.M{
						0: []int{2, 4},
						1: []int{6, 8},
					},
					n:    2,
					want: true,
				},
				{
					name: "two cols miss between first and second",
					m: matrix.M{
						0: []int{2, 4},
						1: []int{6, 8},
					},
					n: 3,
				},
				{
					name: "two cols hit second",
					m: matrix.M{
						0: []int{2, 4},
						1: []int{6, 8},
					},
					n:    4,
					want: true,
				},
				{
					name: "two cols miss between second and third",
					m: matrix.M{
						0: []int{2, 4},
						1: []int{6, 8},
					},
					n: 5,
				},
				{
					name: "two cols hit third",
					m: matrix.M{
						0: []int{2, 4},
						1: []int{6, 8},
					},
					n:    6,
					want: true,
				},
				{
					name: "two cols miss between third and fourth",
					m: matrix.M{
						0: []int{2, 4},
						1: []int{6, 8},
					},
					n: 7,
				},
				{
					name: "two cols hit fourth",
					m: matrix.M{
						0: []int{2, 4},
						1: []int{6, 8},
					},
					n:    8,
					want: true,
				},
				{
					name: "two cols miss above",
					m: matrix.M{
						0: []int{2, 4},
						1: []int{6, 8},
					},
					n: 7,
				},
			},
		},
		{
			name:  "three rows",
			tests: []test{},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) { testMatrix(t, tc.tests) })
	}
}

func testMatrix(t *testing.T, tests []test) {
	t.Helper()
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := matrix.Has(tc.m, tc.n)

			if tc.want != got {
				t.Errorf("matrix.Has(..., %d) got %v; want %v", tc.n, got, tc.want)
				t.Logf("Matrix:\n%s", matrixToString(tc.m))
			}
		})
	}
}

func matrixToString(m matrix.M) string {
	rows := make([]string, 0, len(m))
	for _, row := range m {
		rows = append(rows, fmt.Sprintf("%v", row))
	}
	return strings.Join(rows, "\n")
}
