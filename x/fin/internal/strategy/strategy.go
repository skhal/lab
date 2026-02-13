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

// CycleFunc is a single cycle to process a record, implementing some strategy.
type CycleFunc func(q Quote, prev, curr *pb.Record) Quote

type strategy struct {
	cf CycleFunc
}

// New createas a strategy, backed by the cycle function c.
func New(c CycleFunc) *strategy {
	return &strategy{c}
}

// Run executes the strategy on a set of rectors starting with a given balance.
// It returns the end balance.
func (s *strategy) Run(start fin.Cents, market []*pb.Record) fin.Cents {
	var prev *pb.Record
	q := Quote{Bal: start}
	for _, rec := range market {
		q = s.Cycle(q, prev, rec)
		prev = rec
	}
	return q.Total()
}

// Cycle executes one cycle of a strategy backed by the cycle function.
func (s *strategy) Cycle(q Quote, prev, curr *pb.Record) Quote {
	return s.cf(q, prev, curr)
}
