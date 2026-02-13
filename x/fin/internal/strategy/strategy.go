// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strategy

import (
	"github.com/skhal/lab/x/fin/internal/fin"
	"github.com/skhal/lab/x/fin/internal/pb"
)

// Quote holds account balance and dividends payout.
type Quote struct {
	Bal fin.Cents // account balance
	Div fin.Cents // paid dividends
}

// Total gives a sum for the account balance and dividends.
func (b Quote) Total() fin.Cents {
	return b.Bal + b.Div
}

// Cycler is a strategy that can process a single cycle.
type Cycler interface {
	// Cycle processes a single market record.
	Cycle(Quote, *pb.Record) Quote
}

// Runner represents a strategy backed by Cycler.
type Runner struct {
	c Cycler
}

// New createas a strategy, backed by Cycler.
func New(c Cycler) *Runner {
	return &Runner{c}
}

// Run executes the strategy on a set of rectors starting with a given balance.
// It returns the end balance.
func (s *Runner) Run(start fin.Cents, market []*pb.Record) fin.Cents {
	q := Quote{Bal: start}
	for _, rec := range market {
		q = s.Cycle(q, rec)
	}
	return q.Total()
}

// Cycle executes one cycle of a strategy backed by the cycle function.
func (s *Runner) Cycle(q Quote, rec *pb.Record) Quote {
	return s.c.Cycle(q, rec)
}
