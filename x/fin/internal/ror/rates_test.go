// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ror_test

import (
	"math"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/x/fin/internal/pb"
	"github.com/skhal/lab/x/fin/internal/ror"
	"github.com/skhal/lab/x/fin/internal/tests"
)

func TestRateOfReturn(t *testing.T) {
	tests := []struct {
		name string
		prev *pb.Record
		curr *pb.Record
		want ror.Rate
	}{
		{
			name: "no change",
			prev: tests.NewRecord(t, 2006, time.January, 100, 0, 0),
			curr: tests.NewRecord(t, 2006, time.February, 100, 0, 0),
			want: 1.,
		},
		{
			name: "gain",
			prev: tests.NewRecord(t, 2006, time.January, 100, 0, 0),
			curr: tests.NewRecord(t, 2006, time.February, 105, 0, 0),
			want: 1.05,
		},
		{
			name: "loss",
			prev: tests.NewRecord(t, 2006, time.January, 100, 0, 0),
			curr: tests.NewRecord(t, 2006, time.February, 95, 0, 0),
			want: 0.95,
		},
		{
			name: "no prev zero return",
			curr: tests.NewRecord(t, 2006, time.January, 95, 0, 0),
			want: 1.,
		},
		{
			name: "no curr zero return",
			prev: tests.NewRecord(t, 2006, time.January, 95, 0, 0),
			want: 1.,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := ror.SPComposite(tc.prev, tc.curr)

			if !cmp.Equal(tc.want, r, floatComparer) {
				t.Errorf("RateOfReturn() = %.2f; want %.2f", r, tc.want)
				t.Log("prev:", tc.prev)
				t.Log("curr:", tc.curr)
			}
		})
	}
}

func TestRateOfDividend(t *testing.T) {
	tests := []struct {
		name string
		rec  *pb.Record
		want ror.Rate
	}{
		{
			name: "no dividend",
			rec:  tests.NewRecord(t, 2006, time.January, 100, 0, 0),
		},
		{
			name: "dividend",
			rec:  tests.NewRecord(t, 2006, time.January, 100, 5, 0),
			want: 0.05,
		},
		{
			name: "no sp composite zero return",
			rec:  tests.NewRecord(t, 2005, time.January, 0, 5, 0),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := ror.Dividend(tc.rec)

			if !cmp.Equal(tc.want, r, floatComparer) {
				t.Errorf("RateOfDividend() = %.2f; want %.2f", r, tc.want)
				t.Log(tc.rec)
			}
		})
	}
}

const cutoff = 0.01

var floatComparer = cmp.Comparer(func(x, y float64) bool {
	d := math.Abs(x - y)
	m := math.Abs(x+y) / 2.0
	if m == 0 {
		return d == 0
	}
	return d/m < cutoff
})
