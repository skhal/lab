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
	tests := []struct {
		name   string
		start  fin.Position
		market []*pb.Record
		want   fin.Position
	}{
		{
			name:  "no cycles",
			start: fin.Position{Investment: 123},
			want:  fin.Position{Investment: 123},
		},
		{
			name:  "one cycle div zero",
			start: fin.Position{Investment: 123},
			market: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 100, 0, 0),
			},
			want: func() fin.Position {
				// cycle 1
				c, d := func(c fin.Cents) (inv, div fin.Cents) {
					return c, dividend(t, c, 100, 0)
				}(fin.Cents(123))
				return fin.Position{Investment: c, Dividend: d}
			}(),
		},
		{
			name:  "one cycle div non-zero",
			start: fin.Position{Investment: 123},
			market: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 100, 20, 0),
			},
			want: func() fin.Position {
				// cycle 1
				c, d := func(c fin.Cents) (inv, div fin.Cents) {
					return c, dividend(t, c, 100, 20)
				}(fin.Cents(123))
				return fin.Position{Investment: c, Dividend: d}
			}(),
		},
		{
			name:  "two cycles div zero",
			start: fin.Position{Investment: 123},
			market: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 100, 0, 0),
				tests.NewRecord(t, 2006, time.February, 125, 0, 0),
			},
			want: func() fin.Position {
				// cycle 1
				c, d := func(c fin.Cents) (inv, div fin.Cents) {
					return c, dividend(t, c, 100, 0)
				}(fin.Cents(123))
				// cycle 2
				c, d = invest(t, c, 100, 125), d+dividend(t, c, 125, 0)
				return fin.Position{Investment: c, Dividend: d}
			}(),
		},
		{
			name:  "two cycles div non-zero",
			start: fin.Position{Investment: 123},
			market: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 100, 20, 0),
				tests.NewRecord(t, 2006, time.February, 125, 40, 0),
			},
			want: func() fin.Position {
				// cycle 1
				c, d := func(c fin.Cents) (inv, div fin.Cents) {
					return c, dividend(t, c, 100, 20)
				}(fin.Cents(123))
				// cycle 2
				c, d = invest(t, c, 100, 125), d+dividend(t, c, 125, 40)
				return fin.Position{Investment: c, Dividend: d}
			}(),
		},
		{
			name:  "three cycles div zero",
			start: fin.Position{Investment: 123},
			market: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 100, 0, 0),
				tests.NewRecord(t, 2006, time.February, 125, 0, 0),
				tests.NewRecord(t, 2006, time.March, 150, 0, 0),
			},
			want: func() fin.Position {
				// cycle 1
				c, d := func(c fin.Cents) (inv, div fin.Cents) {
					return c, dividend(t, c, 100, 0)
				}(fin.Cents(123))
				// cycle 2
				c, d = invest(t, c, 100, 125), d+dividend(t, c, 125, 0)
				// cycle 3
				c, d = invest(t, c, 125, 150), d+dividend(t, c, 150, 0)
				return fin.Position{Investment: c, Dividend: d}
			}(),
		},
		{
			name:  "three cycles div non-zero",
			start: fin.Position{Investment: 123},
			market: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 100, 20, 0),
				tests.NewRecord(t, 2006, time.February, 125, 40, 0),
				tests.NewRecord(t, 2006, time.March, 150, 60, 0),
			},
			want: func() fin.Position {
				// cycle 1
				c, d := func(c fin.Cents) (inv, div fin.Cents) {
					return c, dividend(t, c, 100, 20)
				}(fin.Cents(123))
				// cycle 2
				c, d = invest(t, c, 100, 125), d+dividend(t, c, 125, 40)
				// cycle 3
				c, d = invest(t, c, 125, 150), d+dividend(t, c, 150, 60)
				return fin.Position{Investment: c, Dividend: d}
			}(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := strategy.Hold()

			got := s.Run(tc.start, tc.market)

			if got != tc.want {
				t.Errorf("Hold().Run(%s, _) = %s; want %s", tc.start, got, tc.want)
				t.Log(tc.market)
			}
		})
	}
}

func TestHoldReinvest_Run(t *testing.T) {
	tests := []struct {
		name   string
		start  fin.Position
		market []*pb.Record
		want   fin.Position
	}{
		{
			name:  "no cycles",
			start: fin.Position{Investment: 123},
			want:  fin.Position{Investment: 123},
		},
		{
			name:  "one cycle div zero",
			start: fin.Position{Investment: 123},
			market: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 100, 0, 0),
			},
			want: func() fin.Position {
				// cycle 1
				c, d := func(c fin.Cents) (inv, div fin.Cents) {
					return c + dividend(t, c, 100, 0), 0
				}(fin.Cents(123))
				return fin.Position{Investment: c, Dividend: d}
			}(),
		},
		{
			name:  "one cycle div non-zero",
			start: fin.Position{Investment: 123},
			market: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 100, 20, 0),
			},
			want: func() fin.Position {
				// cycle 1
				c, d := func(c fin.Cents) (inv, div fin.Cents) {
					return c + dividend(t, c, 100, 20), 0
				}(fin.Cents(123))
				return fin.Position{Investment: c, Dividend: d}
			}(),
		},
		{
			name:  "two cycles div zero",
			start: fin.Position{Investment: 123},
			market: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 100, 0, 0),
				tests.NewRecord(t, 2006, time.February, 125, 0, 0),
			},
			want: func() fin.Position {
				// cycle 1
				c, d := func(c fin.Cents) (inv, div fin.Cents) {
					return c + dividend(t, c, 100, 0), 0
				}(fin.Cents(123))
				// cycle 2
				c, d = invest(t, c, 100, 125)+dividend(t, c, 125, 0), 0
				return fin.Position{Investment: c, Dividend: d}
			}(),
		},
		{
			name:  "two cycles div non-zero",
			start: fin.Position{Investment: 123},
			market: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 100, 20, 0),
				tests.NewRecord(t, 2006, time.February, 125, 40, 0),
			},
			want: func() fin.Position {
				// cycle 1
				c, d := func(c fin.Cents) (inv, div fin.Cents) {
					return c + dividend(t, c, 100, 20), 0
				}(fin.Cents(123))
				// cycle 2
				c, d = invest(t, c, 100, 125)+dividend(t, c, 125, 40), 0
				return fin.Position{Investment: c, Dividend: d}
			}(),
		},
		{
			name:  "three cycles div zero",
			start: fin.Position{Investment: 123},
			market: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 100, 0, 0),
				tests.NewRecord(t, 2006, time.February, 125, 0, 0),
				tests.NewRecord(t, 2006, time.March, 150, 0, 0),
			},
			want: func() fin.Position {
				// cycle 1
				c, d := func(c fin.Cents) (inv, div fin.Cents) {
					return c + dividend(t, c, 100, 0), 0
				}(fin.Cents(123))
				// cycle 2
				c, d = invest(t, c, 100, 125)+dividend(t, c, 125, 0), 0
				// cycle 3
				c, d = invest(t, c, 125, 150)+dividend(t, c, 150, 0), 0
				return fin.Position{Investment: c, Dividend: d}
			}(),
		},
		{
			name:  "three cycles div non-zero",
			start: fin.Position{Investment: 123},
			market: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 100, 20, 0),
				tests.NewRecord(t, 2006, time.February, 125, 40, 0),
				tests.NewRecord(t, 2006, time.March, 150, 60, 0),
			},
			want: func() fin.Position {
				// cycle 1
				c, d := func(c fin.Cents) (inv, div fin.Cents) {
					return c + dividend(t, c, 100, 20), 0
				}(fin.Cents(123))
				// cycle 2
				c, d = invest(t, c, 100, 125)+dividend(t, c, 125, 40), 0
				// cycle 3
				c, d = invest(t, c, 125, 150)+dividend(t, c, 150, 60), 0
				return fin.Position{Investment: c, Dividend: d}
			}(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := strategy.HoldReinvest().Run(tc.start, tc.market)

			if got != tc.want {
				t.Errorf("HoldReinvest().Run(%s, _) = %s; want %s", tc.start, got, tc.want)
				t.Log(tc.market)
			}
		})
	}
}

func invest(t *testing.T, bal, spa, spb fin.Cents) fin.Cents {
	t.Helper()
	return fin.Cents(math.Floor(float64(bal) * float64(spb) / float64(spa)))
}

func dividend(t *testing.T, bal, sp, div fin.Cents) fin.Cents {
	t.Helper()
	return fin.Cents(math.Floor(float64(bal) * float64(div) / float64(sp)))
}
