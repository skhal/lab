// Copyright 2025 Samvel Khalatyan. All rights reserved.

package tests_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/go/tests"
)

func ExampleEquateFloat64_pass() {
	diff := cmp.Diff(3.4, 3.1, tests.EquateFloat64(0.1))
	if diff != "" {
		fmt.Print("3.4 !~= 3.1")
		return
	}
	fmt.Print("3.4 ~= 3.1")
	// Output:
	// 3.4 ~= 3.1
}

func ExampleEquateFloat64_fail() {
	diff := cmp.Diff(3.4, 3.1, tests.EquateFloat64(0.01))
	if diff != "" {
		fmt.Print("3.4 !~= 3.1")
		return
	}
	fmt.Print("3.4 ~= 3.1")
	// Output:
	// 3.4 !~= 3.1
}

func TestEquateFloat64(t *testing.T) {
	testcases := []struct {
		x, y      float64
		tolerance float64
		wantEqual bool
	}{
		// zero numbers
		{
			wantEqual: true,
		},
		{
			tolerance: 0.1,
			wantEqual: true,
		},
		{
			tolerance: -0.1,
			wantEqual: true,
		},
		// one non-zero, zero tolerance
		{
			x: 1.0,
		},
		{
			x: -1.0,
		},
		{
			y: 1.0,
		},
		{
			y: -1.0,
		},
		// one non-zero, non-zero tolerance
		{
			x:         1.0,
			tolerance: 2.0,
			wantEqual: true,
		},
		{
			x:         1.0,
			tolerance: -2.0,
			wantEqual: true,
		},
		{
			x:         -1.0,
			tolerance: 2.0,
			wantEqual: true,
		},
		{
			x:         -11.0,
			tolerance: -2.0,
			wantEqual: true,
		},
		{
			x:         1.0,
			tolerance: 1.9,
		},
		{
			x:         1.0,
			tolerance: -1.9,
		},
		{
			x:         -1.0,
			tolerance: 1.9,
		},
		{
			x:         -1.0,
			tolerance: -1.9,
		},
		{
			y:         1.0,
			tolerance: 2.0,
			wantEqual: true,
		},
		{
			y:         1.0,
			tolerance: -2.0,
			wantEqual: true,
		},
		{
			y:         -1.0,
			tolerance: 2.0,
			wantEqual: true,
		},
		{
			y:         -11.0,
			tolerance: -2.0,
			wantEqual: true,
		},
		{
			y:         1.0,
			tolerance: 1.9,
		},
		{
			y:         1.0,
			tolerance: -1.9,
		},
		{
			y:         -1.0,
			tolerance: 1.9,
		},
		{
			y:         -1.0,
			tolerance: -1.9,
		},
		// two non-zero, zero tolerance
		{
			x: 1.0,
			y: 2.0,
			tolerance: 0.7,
			wantEqual: true,
		},
		{
			x: 1.0,
			y: 2.0,
			tolerance: 0.6,
		},
		{
			x: -1.0,
			y: 2.0,
			tolerance: 7,
			wantEqual: true,
		},
		{
			x: -1.0,
			y: 2.0,
			tolerance: 5,
		},
		{
			x: 1.0,
			y: -2.0,
			tolerance: 7,
			wantEqual: true,
		},
		{
			x: 1.0,
			y: -2.0,
			tolerance: 5,
		},
		{
			x: -1.0,
			y: -2.0,
			tolerance: 0.7,
			wantEqual: true,
		},
		{
			x: -1.0,
			y: -2.0,
			tolerance: 0.6,
		},
	}
	for i, tc := range testcases {
		tc := tc
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			diff := cmp.Diff(tc.x, tc.y, tests.EquateFloat64(tc.tolerance))

			switch tc.wantEqual {
			case true:
				if diff != "" {
					t.Errorf("cmp.Diff(%f, %f, tests.EquateFloat64(%f)) mismatch - got diff:\n%s", tc.x, tc.y, tc.tolerance, diff)
				}
			case false:
				if diff == "" {
					t.Errorf("cmp.Diff(%f, %f, tests.EquateFloat64(%f)) mismatch - want diff", tc.x, tc.y, tc.tolerance)
				}
			}
		})
	}
}
