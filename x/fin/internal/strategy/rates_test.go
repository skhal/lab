// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strategy_test

import (
	"math"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/x/fin/internal/pb"
	"github.com/skhal/lab/x/fin/internal/strategy"
)

func TestRateOfReturn(t *testing.T) {
	tests := []struct {
		name string
		prev *pb.Record
		curr *pb.Record
		want strategy.Rate
	}{
		{
			name: "no change",
			prev: newRecord(t, time.January, 100, 0),
			curr: newRecord(t, time.February, 100, 0),
			want: 1.,
		},
		{
			name: "gain",
			prev: newRecord(t, time.January, 100, 0),
			curr: newRecord(t, time.February, 105, 0),
			want: 1.05,
		},
		{
			name: "loss",
			prev: newRecord(t, time.January, 100, 0),
			curr: newRecord(t, time.February, 95, 0),
			want: 0.95,
		},
		{
			name: "no prev zero return",
			curr: newRecord(t, time.January, 95, 0),
		},
		{
			name: "no curr zero return",
			prev: newRecord(t, time.January, 95, 0),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := strategy.SPRateOfReturn(tc.prev, tc.curr)

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
		want strategy.Rate
	}{
		{
			name: "no dividend",
			rec:  newRecord(t, time.January, 100, 0),
		},
		{
			name: "dividend",
			rec:  newRecord(t, time.January, 100, 5),
			want: 0.05,
		},
		{
			name: "no sp composite zero return",
			rec:  newRecord(t, time.January, 0, 5),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := strategy.DivRateOfReturn(tc.rec)

			if !cmp.Equal(tc.want, r, floatComparer) {
				t.Errorf("RateOfDividend() = %.2f; want %.2f", r, tc.want)
				t.Log(tc.rec)
			}
		})
	}
}

func newRecord(t *testing.T, m time.Month, sp, div int32) *pb.Record {
	t.Helper()
	year := int32(2006)
	month := int32(m)
	return pb.Record_builder{
		Date: pb.Date_builder{
			Year:  &year,
			Month: &month,
		}.Build(),
		Quote: pb.Quote_builder{
			SpComposite: pb.Cents_builder{Cents: &sp}.Build(),
			Dividend:    pb.Cents_builder{Cents: &div}.Build(),
		}.Build(),
	}.Build()
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
