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
func Run(c fin.Cents, market []*pb.Record, r *strategy.Runner) (start, end fin.Balance) {
	if len(market) == 0 {
		return
	}
	start = fin.Balance{
		Date: newTime(market[0].GetDate()),
		Cash: c,
	}
	pos := fin.Position{
		Investment: c,
	}
	pos = r.Run(pos, market)
	d := newTime(market[len(market)-1].GetDate())
	end = fin.Balance{
		Date: nextMonth(d),
		Cash: pos.Total(),
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
