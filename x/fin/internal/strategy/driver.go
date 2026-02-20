// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strategy

import (
	"iter"
	"math"
	"slices"
	"time"

	"github.com/skhal/lab/x/fin/internal/fin"
	"github.com/skhal/lab/x/fin/internal/pb"
	"github.com/skhal/lab/x/fin/internal/ror"
)

// RebalanceFunc implements a position rebalance at the end of the month.
type RebalanceFunc func(fin.Position) fin.Position

// Drive runs balance balance through market records. It returns a sequence of
// balance changes.
func Drive(cash fin.Cents, recs []*pb.Record, reb ...RebalanceFunc) []fin.Balance {
	d := &driver{rebfns: reb}
	return d.drive(cash, recs)
}

type driver struct {
	last   *pb.Record
	rebfns []RebalanceFunc
}

// drive runs the balanace through market records. It is responsible for
// opening and closing a position using cash available in the balance, before
// running through market records.
func (d *driver) drive(cash fin.Cents, recs []*pb.Record) []fin.Balance {
	bals := make([]fin.Balance, 0, 2+len(recs)) // +2 for open, close
	for bal := range d.run(cash, recs) {
		bals = append(bals, bal)
	}
	return bals
}

func (d *driver) run(cash fin.Cents, recs []*pb.Record) iter.Seq[fin.Balance] {
	nextMonth := func(t time.Time) time.Time {
		var y, m, d int
		m = 1
		return t.AddDate(y, m, d)
	}
	return func(yield func(fin.Balance) bool) {
		bal := fin.Balance{
			Cash: cash,
		}
		if !yield(bal) {
			return
		}
		if len(recs) == 0 {
			return
		}
		bal = d.openPosition(bal)
		// do not yield a balance with open position
		for rec := range slices.Values(recs) {
			bal = d.process(bal, rec)
			if !yield(bal) {
				return
			}
		}
		bal = d.closePosition(bal)
		bal.Date = nextMonth(bal.Date)
		yield(bal)
	}
}

func (d *driver) openPosition(bal fin.Balance) fin.Balance {
	bal.Position = fin.Position{
		Investment: bal.Cash,
	}
	bal.Cash = 0
	return bal
}

func (d *driver) process(bal fin.Balance, rec *pb.Record) fin.Balance {
	defer func() { d.last = rec }()
	pos := d.update(bal.Position, rec)
	for fn := range slices.Values(d.rebfns) {
		pos = fn(pos)
	}
	return fin.Balance{
		Date:     newTime(rec.GetDate()),
		Cash:     bal.Cash,
		Position: pos,
	}
}

func (d *driver) update(pos fin.Position, rec *pb.Record) fin.Position {
	inv := d.returnOnInvestment(pos.Investment, rec)
	div := d.payDividend(pos.Investment, rec)
	return fin.Position{
		Investment: inv,
		Dividend:   pos.Dividend + div,
	}
}

func (d *driver) returnOnInvestment(c fin.Cents, rec *pb.Record) fin.Cents {
	if d.last == nil {
		return c
	}
	rate := ror.SPComposite(d.last, rec)
	return fin.Cents(math.Floor(float64(c) * float64(rate)))
}

func (d *driver) payDividend(c fin.Cents, rec *pb.Record) fin.Cents {
	rate := ror.Dividend(rec)
	return fin.Cents(math.Floor(float64(c) * float64(rate)))
}

func (d *driver) closePosition(bal fin.Balance) fin.Balance {
	bal.Cash += bal.Position.Total()
	bal.Position = fin.Position{}
	return bal
}

func newTime(date *pb.Date) time.Time {
	y := int(date.GetYear())
	m := time.Month(date.GetMonth())
	d := 1
	var (
		hh, mm, ss, ns int
		tz             = time.Local // local TZ is ok
	)
	return time.Date(y, m, d, hh, mm, ss, ns, tz)
}
