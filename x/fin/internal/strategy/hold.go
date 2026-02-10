// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strategy

import (
	"math"

	"github.com/skhal/lab/x/fin/internal/fin"
	"github.com/skhal/lab/x/fin/internal/pb"
)

// Hold implements a strategy to hold investment. It has an option to re-invest
// dividends (off by default).
type Hold struct {
	reinvestDividends bool
}

// NewHold createsa a hold strategy.
func NewHold(opts ...HoldOpt) *Hold {
	s := new(Hold)
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// HoldOpt is an option
type HoldOpt func(*Hold)

// HoldOptReinvestDiv turns on dividend re-investment in the Hold strategy.
func HoldOptReinvestDiv() HoldOpt {
	return func(s *Hold) {
		s.reinvestDividends = true
	}
}

type quote struct {
	bal fin.Cents
	div fin.Cents
}

func (b quote) total() fin.Cents {
	return b.bal + b.div
}

// Run executes the strategy.
func (s *Hold) Run(c fin.Cents, market []*pb.Record) fin.Cents {
	var prev *pb.Record
	q := quote{bal: c}
	for _, rec := range market {
		q = s.cycle(q, prev, rec)
		prev = rec
	}
	return q.total()
}

func (s *Hold) cycle(q quote, prev, curr *pb.Record) quote {
	bal := s.invest(q.bal, prev, curr)
	div := s.payDividend(q.bal, curr)
	if s.reinvestDividends {
		bal += div
		div = 0
	} else {
		div += q.div
	}
	return quote{bal: bal, div: div}
}

func (s *Hold) invest(c fin.Cents, prev, curr *pb.Record) fin.Cents {
	ror := SPRateOfReturn(prev, curr)
	return fin.Cents(math.Floor(float64(c) * float64(ror)))
}

func (s *Hold) payDividend(c fin.Cents, rec *pb.Record) fin.Cents {
	ror := DivRateOfReturn(rec)
	return fin.Cents(math.Floor(float64(c) * float64(ror)))
}
