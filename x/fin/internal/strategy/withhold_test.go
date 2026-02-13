// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strategy_test

import (
	"testing"
	"time"

	"github.com/skhal/lab/x/fin/internal/fin"
	"github.com/skhal/lab/x/fin/internal/pb"
	"github.com/skhal/lab/x/fin/internal/strategy"
	"github.com/skhal/lab/x/fin/internal/tests"
)

func TestWithhold_Run(t *testing.T) {
	cycle := func(q strategy.Quote, _ *pb.Record) strategy.Quote {
		return q
	}
	tctc := []struct {
		name    string
		percent strategy.Percent
		start   fin.Cents
		market  []*pb.Record
		want    fin.Cents
	}{
		{
			name:  "zero percent no year",
			start: fin.Cents(100),
			market: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 111, 122, 133),
				tests.NewRecord(t, 2006, time.February, 211, 222, 233),
			},
			want: fin.Cents(100),
		},
		{
			name:  "zero percent with year",
			start: fin.Cents(100),
			market: []*pb.Record{
				tests.NewRecord(t, 2006, time.December, 111, 122, 133),
				tests.NewRecord(t, 2007, time.January, 211, 222, 233),
			},
			want: fin.Cents(100),
		},
		{
			name:    "one percent no year",
			percent: strategy.Percent(1),
			start:   fin.Cents(100),
			market: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 111, 122, 133),
				tests.NewRecord(t, 2006, time.February, 211, 222, 233),
			},
			want: fin.Cents(100),
		},
		{
			name:    "one percent with year",
			percent: strategy.Percent(1),
			start:   fin.Cents(100),
			market: []*pb.Record{
				tests.NewRecord(t, 2006, time.December, 111, 122, 133),
				tests.NewRecord(t, 2007, time.January, 211, 222, 233),
			},
			want: fin.Cents(99),
		},
	}
	for _, tc := range tctc {
		t.Run(tc.name, func(t *testing.T) {
			s := strategy.NewWithhold(strategy.New(CycleFunc(cycle)), tc.percent)

			got := s.Run(tc.start, tc.market)

			if got != tc.want {
				t.Errorf("NewHold().Run(%s, _) = %s; want %s", tc.start, got, tc.want)
				t.Log(tc.market)
			}
		})
	}
}

type CycleFunc func(strategy.Quote, *pb.Record) strategy.Quote

func (cf CycleFunc) Cycle(q strategy.Quote, rec *pb.Record) strategy.Quote {
	return cf(q, rec)
}
