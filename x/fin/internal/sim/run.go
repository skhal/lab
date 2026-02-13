// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sim

import (
	"time"

	"github.com/skhal/lab/x/fin/internal/fin"
	"github.com/skhal/lab/x/fin/internal/pb"
	"github.com/skhal/lab/x/fin/internal/strategy"
)

// Run executes strategy s for market. It return the beginning and end balances.
func Run(c fin.Cents, market []*pb.Record, s *strategy.Runner) (start, end fin.Quote) {
	if len(market) == 0 {
		return
	}
	start = fin.Quote{
		Date:    newTime(market[0].GetDate()),
		Balance: c,
	}
	d := newTime(market[len(market)-1].GetDate())
	end = fin.Quote{
		Date:    nextMonth(d),
		Balance: s.Run(c, market),
	}
	return start, end
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
