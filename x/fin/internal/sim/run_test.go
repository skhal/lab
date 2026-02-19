// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sim_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/x/fin/internal/fin"
	"github.com/skhal/lab/x/fin/internal/pb"
	"github.com/skhal/lab/x/fin/internal/sim"
	"github.com/skhal/lab/x/fin/internal/strategy"
	"github.com/skhal/lab/x/fin/internal/tests"
)

func TestRun(t *testing.T) {
	cycle := func(pos fin.Position, _ *pb.Record) fin.Position {
		return fin.Position{Investment: -pos.Investment, Dividend: -pos.Dividend}
	}
	type want struct {
		start fin.Quote
		end   fin.Quote
	}
	tests := []struct {
		name   string
		bal    fin.Cents
		market []*pb.Record
		want   want
	}{
		{
			name: "empty do nothing",
		},
		{
			name: "one record call strategy",
			bal:  123,
			market: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 2, 3, 0),
			},
			want: want{
				start: fin.Quote{Date: newTime(t, 2006, time.January), Balance: 123},
				end:   fin.Quote{Date: newTime(t, 2006, time.February), Balance: -123},
			},
		},
		{
			name: "two records call strategy",
			bal:  123,
			market: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 2, 3, 0),
				tests.NewRecord(t, 2006, time.February, 2, 3, 0),
			},
			want: want{
				start: fin.Quote{Date: newTime(t, 2006, time.January), Balance: 123},
				end:   fin.Quote{Date: newTime(t, 2006, time.March), Balance: 123},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			start, end := sim.Run(tc.bal, tc.market, strategy.New(CycleFunc(cycle)))

			if diff := cmp.Diff(tc.want.start, start); diff != "" {
				t.Errorf("sim.Run() = start, _; mismatch (-want, +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.end, end); diff != "" {
				t.Errorf("sim.Run() = _, end; mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

type CycleFunc func(fin.Position, *pb.Record) fin.Position

func (cf CycleFunc) Cycle(pos fin.Position, rec *pb.Record) fin.Position {
	return cf(pos, rec)
}

func newTime(t *testing.T, year int, month time.Month) time.Time {
	t.Helper()
	d := 1
	var hh, mm, ss, ns int
	return time.Date(year, month, d, hh, mm, ss, ns, time.Local)
}
