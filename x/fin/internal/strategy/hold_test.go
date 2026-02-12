// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strategy_test

import (
	"math"
	"testing"
	"time"

	"github.com/skhal/lab/x/fin/internal/fin"
	"github.com/skhal/lab/x/fin/internal/pb"
	"github.com/skhal/lab/x/fin/internal/strategy"
	"github.com/skhal/lab/x/fin/internal/tests"
)

func TestHold_Run(t *testing.T) {
	invest := func(bal, spa, spb fin.Cents) fin.Cents {
		return fin.Cents(math.Floor(float64(bal) * float64(spb) / float64(spa)))
	}
	dividend := func(bal, sp, div fin.Cents) fin.Cents {
		return fin.Cents(math.Floor(float64(bal) * float64(div) / float64(sp)))
	}
	tests := []struct {
		name   string
		start  fin.Cents
		market []*pb.Record
		want   fin.Cents
	}{
		{
			name:  "no cycles",
			start: fin.Cents(123),
			want:  fin.Cents(123),
		},
		{
			name:  "one cycle div zero",
			start: fin.Cents(123),
			market: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 100, 0),
			},
			want: fin.Cents(123),
		},
		{
			name:  "one cycle div non-zero",
			start: fin.Cents(123),
			market: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 100, 1),
			},
			want: fin.Cents(math.Floor(123 * (1 + 1./100.))),
		},
		{
			name:  "two cycles div zero",
			start: fin.Cents(123),
			market: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 100, 0),
				tests.NewRecord(t, 2006, time.February, 110, 0),
			},
			want: func() fin.Cents {
				// cycle 1
				c := fin.Cents(123)
				d := dividend(c, 100, 0)
				// cycle 2
				c = invest(c, 100, 110) // update
				d += dividend(c, 110, 0)
				return c + d
			}(),
		},
		{
			name:  "two cycles div non-zero",
			start: fin.Cents(123),
			market: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 100, 1),
				tests.NewRecord(t, 2006, time.February, 110, 2),
			},
			want: func() fin.Cents {
				// cycle 1
				c := fin.Cents(123)
				d := dividend(c, 100, 1)
				// cycle 2
				c = invest(c, 100, 110) // update
				d += dividend(c, 110, 2)
				return c + d
			}(),
		},
		{
			name:  "three cycles div zero",
			start: fin.Cents(123),
			market: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 100, 0),
				tests.NewRecord(t, 2006, time.February, 110, 0),
				tests.NewRecord(t, 2006, time.March, 120, 0),
			},
			want: func() fin.Cents {
				// cycle 1
				c := fin.Cents(123)
				d := dividend(c, 100, 0)
				// cycle 2
				c = invest(c, 100, 110) // update
				d += dividend(c, 110, 0)
				// cycle 3
				c = invest(c, 110, 120) // update
				d += dividend(c, 120, 0)
				return c + d
			}(),
		},
		{
			name:  "three cycles div non-zero",
			start: fin.Cents(123),
			market: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 100, 1),
				tests.NewRecord(t, 2006, time.February, 110, 2),
				tests.NewRecord(t, 2006, time.March, 120, 3),
			},
			want: func() fin.Cents {
				// cycle 1
				c := fin.Cents(123)
				d := dividend(c, 100, 1)
				// cycle 2
				c = invest(c, 100, 110) // update
				d += dividend(c, 110, 2)
				// cycle 3
				c = invest(c, 110, 120) // update
				d += dividend(c, 120, 3)
				return c + d
			}(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := strategy.NewHold()

			got := s.Run(tc.start, tc.market)

			if got != tc.want {
				t.Errorf("NewHold().Run(%s, _) = %s; want %s", tc.start, got, tc.want)
				t.Log(tc.market)
			}
		})
	}
}
