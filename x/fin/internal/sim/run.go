// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sim drives simulation of the market.
package sim

import (
	"fmt"
	"time"

	"github.com/skhal/lab/x/fin/internal/fin"
	"github.com/skhal/lab/x/fin/internal/pb"
)

// Strategy emulates an investment behavior.
type Strategy interface {
	// Run applies underlying strategy to the market. It begins with startBalance
	// and returns endBalance.
	Run(startBalance fin.Cents, market []*pb.Record) (endBalance fin.Cents)
}

// Run executes strategy s for market. It return the beginning and end balances.
func Run(c fin.Cents, market []*pb.Record, s Strategy) (start, end Quote) {
	if len(market) == 0 {
		return
	}
	start = Quote{
		Date:    newTime(market[0].GetDate()),
		Balance: c,
	}
	d := newTime(market[len(market)-1].GetDate())
	end = Quote{
		Date:    nextMonth(d),
		Balance: s.Run(c, market),
	}
	return start, end
}

// Quote is a the balance at the beginning of a given month.
// A value `2006 Jan: $1.23` means $1.23 balance as of Jan 1, 2006.
type Quote struct {
	Date    time.Time // first day of the month
	Balance fin.Cents // balance on the date
}

// String implements fmt.Stringer interface.
func (s Quote) String() string {
	d := s.Date.Format("2006 Jan")
	return fmt.Sprintf("%s: $%s", d, s.Balance)
}

func newTime(date *pb.Date) time.Time {
	y := int(date.GetYear())
	m := time.Month(date.GetMonth())
	d := 1
	var hh, mm, ss, ns int
	return time.Date(y, m, d, hh, mm, ss, ns, time.Local) // local TZ is ok
}

func nextMonth(t time.Time) time.Time {
	return t.AddDate(0, 1, 0)
}
