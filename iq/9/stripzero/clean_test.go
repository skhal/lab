// Copyright 2025 Samvel Khalatyan. All rights reserved.

package stripzero_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/iq/9/stripzero"
)

type matrix [][]int

func (m matrix) String() string {
	buf := new(bytes.Buffer)
	for _, nn := range m {
		fmt.Fprintln(buf, nn)
	}
	return buf.String()
}

func DeepCopy(t *testing.T, m matrix) matrix {
	t.Helper()
	if m == nil {
		return nil
	}
	buf := make([][]int, 0, len(m))
	for _, nn := range m {
		dup := make([]int, len(nn))
		copy(dup, nn)
		buf = append(buf, dup)
	}
	return matrix(buf)
}

var tests = []struct {
	name string
	m    matrix
	want matrix
}{
	{
		name: "empty",
	},
	// 1x1
	{
		name: "one by one non zero",
		m: [][]int{
			{1},
		},
		want: [][]int{
			{1},
		},
	},
	{
		name: "one by one zero",
		m: [][]int{
			{0},
		},
		want: [][]int{
			{0},
		},
	},
	// row
	{
		name: "row non zero",
		m: [][]int{
			{1, 2, 3},
		},
		want: [][]int{
			{1, 2, 3},
		},
	},
	{
		name: "row with zero",
		m: [][]int{
			{1, 2, 0},
		},
		want: [][]int{
			{0, 0, 0},
		},
	},
	// column
	{
		name: "col non zero",
		m: [][]int{
			{1},
			{2},
			{3},
		},
		want: [][]int{
			{1},
			{2},
			{3},
		},
	},
	{
		name: "col with zero",
		m: [][]int{
			{1},
			{0},
			{3},
		},
		want: [][]int{
			{0},
			{0},
			{0},
		},
	},
	// m-by-n matrix
	{
		name: "matrix zero in first row",
		m: [][]int{
			{1, 0, 1, 1},
			{1, 1, 1, 1},
			{1, 1, 1, 1},
		},
		want: [][]int{
			{0, 0, 0, 0},
			{1, 0, 1, 1},
			{1, 0, 1, 1},
		},
	},
	{
		name: "matrix two zeros in first row",
		m: [][]int{
			{1, 0, 0, 1},
			{1, 1, 1, 1},
			{1, 1, 1, 1},
		},
		want: [][]int{
			{0, 0, 0, 0},
			{1, 0, 0, 1},
			{1, 0, 0, 1},
		},
	},
	{
		name: "matrix zero in first col",
		m: [][]int{
			{1, 1, 1, 1},
			{1, 1, 1, 1},
			{0, 1, 1, 1},
		},
		want: [][]int{
			{0, 1, 1, 1},
			{0, 1, 1, 1},
			{0, 0, 0, 0},
		},
	},
	{
		name: "matrix two zeros in first col",
		m: [][]int{
			{1, 1, 1, 1},
			{0, 1, 1, 1},
			{0, 1, 1, 1},
		},
		want: [][]int{
			{0, 1, 1, 1},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
	},
	{
		name: "matrix non zero",
		m: [][]int{
			{1, 1, 1, 1},
			{1, 1, 1, 1},
			{1, 1, 1, 1},
		},
		want: [][]int{
			{1, 1, 1, 1},
			{1, 1, 1, 1},
			{1, 1, 1, 1},
		},
	},
	{
		name: "matrix one zero",
		m: [][]int{
			{1, 1, 1, 1},
			{1, 0, 1, 1},
			{1, 1, 1, 1},
		},
		want: [][]int{
			{1, 0, 1, 1},
			{0, 0, 0, 0},
			{1, 0, 1, 1},
		},
	},
	{
		name: "matrix two zero",
		m: [][]int{
			{1, 0, 1, 1},
			{1, 1, 1, 1},
			{1, 1, 1, 0},
		},
		want: [][]int{
			{0, 0, 0, 0},
			{1, 0, 1, 0},
			{0, 0, 0, 0},
		},
	},
}

func TestClean(t *testing.T) {
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			m := DeepCopy(t, tc.m)

			stripzero.Clean(m)

			if diff := cmp.Diff(tc.want, m); diff != "" {
				t.Errorf("stripzero.Clean(...) mismatch (-want, +got):\n%s", diff)
				t.Logf("Input:\n%s", tc.m)
			}
		})
	}
}
