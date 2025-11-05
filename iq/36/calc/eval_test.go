// Copyright 2025 Samvel Khalatyan. All rights reserved.

package calc_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/go/tests"
	"github.com/skhal/lab/iq/36/calc"
)

func ExampleEval() {
	res, err := calc.Eval("3 + 4 - 2")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	// Output:
	// 5
}

func TestEval_valid(t *testing.T) {
	testcases := []struct {
		name    string
		s       string
		want    int
		wantErr error
	}{
		{
			name: "empty",
		},
		{
			name: "number",
			s:    "1",
			want: 1,
		},
		{
			name: "plus",
			s:    "1 + 2",
			want: 3,
		},
		{
			name: "minus",
			s:    "1 - 2",
			want: -1,
		},
		{
			name: "parenthesis number",
			s:    "(1)",
			want: 1,
		},
		{
			name: "parenthesis parenthesis number",
			s:    "((1))",
			want: 1,
		},
		{
			name: "parenthesis plus",
			s:    "(1 + 2)",
			want: 3,
		},
		{
			name: "parenthesis minus",
			s:    "(1 - 2)",
			want: -1,
		},
		{
			name: "plus parenthesis plus",
			s:    "1 + (3 + 2)",
			want: 6,
		},
		{
			name: "plus parenthesis minus",
			s:    "1 + (3 - 2)",
			want: 2,
		},
		{
			name: "minus parenthesis plus",
			s:    "1 - (3 + 2)",
			want: -4,
		},
		{
			name: "minus parenthesis minus",
			s:    "1 - (3 - 4)",
			want: 2,
		},
		{
			name: "parenthesis plus plus",
			s:    "(3 + 2) + 1",
			want: 6,
		},
		{
			name: "parenthesis minus plus",
			s:    "(3 - 2) + 1",
			want: 2,
		},
		{
			name: "parenthesis plus minus",
			s:    "(3 + 2) - 1",
			want: 4,
		},
		{
			name: "parenthesis minus minus",
			s:    "(3 - 2) - 4",
			want: -3,
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := calc.Eval(tc.s)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("calc.Eval(%q) = _, %v; want error %v", tc.s, err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got, tests.EquateFloat64(0.001)); diff != "" {
				t.Errorf("calc.Eval(%q) mismatch (-want, +got):\n%s", tc.s, diff)
			}
		})
	}
}
