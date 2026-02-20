// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strategy_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/x/fin/internal/fin"
	"github.com/skhal/lab/x/fin/internal/pb"
	"github.com/skhal/lab/x/fin/internal/strategy"
	"github.com/skhal/lab/x/fin/internal/tests"
)

func TestDrive_norebalance(t *testing.T) {
	tt := []struct {
		name string
		cash fin.Cents
		recs []*pb.Record
		want []fin.Balance
	}{
		{
			name: "no records",
			cash: fin.Cents(123),
			want: []fin.Balance{
				{Cash: fin.Cents(123)},
			},
		},
		{
			name: "one record no dividend",
			cash: fin.Cents(123),
			recs: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 100, 0, 0),
			},
			want: []fin.Balance{
				{Cash: fin.Cents(123)},
				tests.NewBalance(t, 2006, time.January, 0, tests.NewPosition(t, 123, 0)),
				tests.NewBalance(t, 2006, time.February, 123),
			},
		},
		{
			name: "one record with dividend",
			cash: fin.Cents(123),
			recs: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 100, 20, 0),
			},
			want: []fin.Balance{
				{Cash: fin.Cents(123)},
				tests.NewBalance(t, 2006, time.January, 0, func() fin.Position {
					c, d := func(c fin.Cents) (inv, div fin.Cents) {
						return c, dividend(t, c, 100, 20)
					}(fin.Cents(123))
					return fin.Position{Investment: c, Dividend: d}
				}()),
				tests.NewBalance(t, 2006, time.February, func() int64 {
					c, d := func(c fin.Cents) (inv, div fin.Cents) {
						return c, dividend(t, c, 100, 20)
					}(fin.Cents(123))
					return int64(c + d)
				}()),
			},
		},
		{
			name: "two records no dividend",
			cash: fin.Cents(123),
			recs: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 100, 0, 0),
				tests.NewRecord(t, 2006, time.February, 125, 0, 0),
			},
			want: []fin.Balance{
				{Cash: fin.Cents(123)},
				tests.NewBalance(t, 2006, time.January, 0, tests.NewPosition(t, 123, 0)),
				tests.NewBalance(t, 2006, time.February, 0, func() fin.Position {
					// record 1
					c, d := func(c fin.Cents) (inv, div fin.Cents) {
						return c, 0
					}(fin.Cents(123))
					// record 2
					c, d = invest(t, c, 100, 125), 0
					return fin.Position{Investment: c, Dividend: d}
				}()),
				tests.NewBalance(t, 2006, time.March, func() int64 {
					// record 1
					c, d := func(c fin.Cents) (inv, div fin.Cents) {
						return c, 0
					}(fin.Cents(123))
					// record 2
					c, d = invest(t, c, 100, 125), 0
					return int64(c + d)
				}()),
			},
		},
		{
			name: "two records with dividend",
			cash: fin.Cents(123),
			recs: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 100, 20, 0),
				tests.NewRecord(t, 2006, time.February, 125, 40, 0),
			},
			want: []fin.Balance{
				{Cash: fin.Cents(123)},
				tests.NewBalance(t, 2006, time.January, 0, func() fin.Position {
					// record 1
					c, d := func(c fin.Cents) (inv, div fin.Cents) {
						return c, dividend(t, c, 100, 20)
					}(fin.Cents(123))
					return fin.Position{Investment: c, Dividend: d}
				}()),
				tests.NewBalance(t, 2006, time.February, 0, func() fin.Position {
					// record 1
					c, d := func(c fin.Cents) (inv, div fin.Cents) {
						return c, dividend(t, c, 100, 20)
					}(fin.Cents(123))
					// record 2
					c, d = invest(t, c, 100, 125), d+dividend(t, c, 125, 40)
					return fin.Position{Investment: c, Dividend: d}
				}()),
				tests.NewBalance(t, 2006, time.March, func() int64 {
					// record 1
					c, d := func(c fin.Cents) (inv, div fin.Cents) {
						return c, dividend(t, c, 100, 20)
					}(fin.Cents(123))
					// record 2
					c, d = invest(t, c, 100, 125), d+dividend(t, c, 125, 40)
					return int64(c + d)
				}()),
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			bals := strategy.Drive(tc.cash, tc.recs)

			if diff := cmp.Diff(tc.want, bals); diff != "" {
				t.Error("Drive() mismatch (-want,+got):\n", diff)
				tests.LogRecords(t, tc.recs)
			}
		})
	}
}

func TestDrive_rebalance(t *testing.T) {
	doublePosition := func(pos fin.Position) fin.Position {
		return fin.Position{
			Investment: pos.Investment * 2,
			Dividend:   pos.Dividend * 2,
		}
	}
	tt := []struct {
		name string
		cash fin.Cents
		recs []*pb.Record
		rebf strategy.RebalanceFunc
		want []fin.Balance
	}{
		{
			name: "no records",
			cash: fin.Cents(123),
			rebf: doublePosition,
			want: []fin.Balance{
				{Cash: fin.Cents(123)},
			},
		},
		{
			name: "one record no dividend",
			cash: fin.Cents(123),
			recs: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 100, 0, 0),
			},
			rebf: doublePosition,
			want: []fin.Balance{
				{Cash: fin.Cents(123)},
				tests.NewBalance(t, 2006, time.January, 0, func() fin.Position {
					// record 1
					pos := tests.NewPosition(t, 123, 0)
					return doublePosition(pos)
				}()),
				tests.NewBalance(t, 2006, time.February, func() int64 {
					// record 1
					pos := tests.NewPosition(t, 123, 0)
					return int64(doublePosition(pos).Total())
				}()),
			},
		},
		{
			name: "one record with dividend",
			cash: fin.Cents(123),
			recs: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 100, 20, 0),
			},
			rebf: doublePosition,
			want: []fin.Balance{
				{Cash: fin.Cents(123)},
				tests.NewBalance(t, 2006, time.January, 0, func() fin.Position {
					// record 1
					c, d := func(c fin.Cents) (inv, div fin.Cents) {
						return c, dividend(t, c, 100, 20)
					}(fin.Cents(123))
					pos := fin.Position{Investment: c, Dividend: d}
					return doublePosition(pos)
				}()),
				tests.NewBalance(t, 2006, time.February, func() int64 {
					// record 1
					c, d := func(c fin.Cents) (inv, div fin.Cents) {
						return c, dividend(t, c, 100, 20)
					}(fin.Cents(123))
					pos := fin.Position{Investment: c, Dividend: d}
					return int64(doublePosition(pos).Total())
				}()),
			},
		},
		{
			name: "two records no dividend",
			cash: fin.Cents(123),
			recs: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 100, 0, 0),
				tests.NewRecord(t, 2006, time.February, 125, 0, 0),
			},
			rebf: doublePosition,
			want: []fin.Balance{
				{Cash: fin.Cents(123)},
				tests.NewBalance(t, 2006, time.January, 0, func() fin.Position {
					// record 1
					c, d := func(c fin.Cents) (inv, div fin.Cents) {
						return c, 0
					}(fin.Cents(123))
					pos := fin.Position{Investment: c, Dividend: d}
					return doublePosition(pos)
				}()),
				tests.NewBalance(t, 2006, time.February, 0, func() fin.Position {
					// record 1
					c, d := func(c fin.Cents) (inv, div fin.Cents) {
						return c, 0
					}(fin.Cents(123))
					pos := fin.Position{Investment: c, Dividend: d}
					pos = doublePosition(pos)
					// record 2
					c, d = invest(t, pos.Investment, 100, 125), 0
					pos = fin.Position{Investment: c, Dividend: d}
					return doublePosition(pos)
				}()),
				tests.NewBalance(t, 2006, time.March, func() int64 {
					// record 1
					c, d := func(c fin.Cents) (inv, div fin.Cents) {
						return c, 0
					}(fin.Cents(123))
					pos := fin.Position{Investment: c, Dividend: d}
					pos = doublePosition(pos)
					// record 2
					c, d = invest(t, pos.Investment, 100, 125), 0
					pos = fin.Position{Investment: c, Dividend: d}
					return int64(doublePosition(pos).Total())
				}()),
			},
		},
		{
			name: "two records with dividend",
			cash: fin.Cents(123),
			recs: []*pb.Record{
				tests.NewRecord(t, 2006, time.January, 100, 20, 0),
				tests.NewRecord(t, 2006, time.February, 125, 40, 0),
			},
			rebf: doublePosition,
			want: []fin.Balance{
				{Cash: fin.Cents(123)},
				tests.NewBalance(t, 2006, time.January, 0, func() fin.Position {
					// record 1
					c, d := func(c fin.Cents) (inv, div fin.Cents) {
						return c, dividend(t, c, 100, 20)
					}(fin.Cents(123))
					pos := fin.Position{Investment: c, Dividend: d}
					return doublePosition(pos)
				}()),
				tests.NewBalance(t, 2006, time.February, 0, func() fin.Position {
					// record 1
					c, d := func(c fin.Cents) (inv, div fin.Cents) {
						return c, dividend(t, c, 100, 20)
					}(fin.Cents(123))
					pos := fin.Position{Investment: c, Dividend: d}
					pos = doublePosition(pos)
					// record 2
					c, d = invest(t, pos.Investment, 100, 125), pos.Dividend+dividend(t, pos.Investment, 125, 40)
					pos = fin.Position{Investment: c, Dividend: d}
					return doublePosition(pos)
				}()),
				tests.NewBalance(t, 2006, time.March, func() int64 {
					// record 1
					c, d := func(c fin.Cents) (inv, div fin.Cents) {
						return c, dividend(t, c, 100, 20)
					}(fin.Cents(123))
					pos := fin.Position{Investment: c, Dividend: d}
					pos = doublePosition(pos)
					// record 2
					c, d = invest(t, pos.Investment, 100, 125), pos.Dividend+dividend(t, pos.Investment, 125, 40)
					pos = fin.Position{Investment: c, Dividend: d}
					return int64(doublePosition(pos).Total())
				}()),
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			bals := strategy.Drive(tc.cash, tc.recs, tc.rebf)

			if diff := cmp.Diff(tc.want, bals); diff != "" {
				t.Error("Drive() mismatch (-want,+got):\n", diff)
				tests.LogRecords(t, tc.recs)
			}
		})
	}
}
