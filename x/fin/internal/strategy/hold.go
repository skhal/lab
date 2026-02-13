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
	last              *pb.Record
}

// NewHold createsa a hold strategy.
func NewHold(opts ...HoldOpt) *Runner {
	h := new(Hold)
	for _, opt := range opts {
		opt(h)
	}
	return New(h)
}

// HoldOpt is an option
type HoldOpt func(*Hold)

// HoldOptReinvestDiv turns on dividend re-investment in the Hold strategy.
func HoldOptReinvestDiv() HoldOpt {
	return func(s *Hold) {
		s.reinvestDividends = true
	}
}

// Cycle executes a single cycle of the hold strategy.
func (s *Hold) Cycle(q Quote, rec *pb.Record) Quote {
	bal := s.invest(q.Bal, rec)
	div := s.payDividend(q.Bal, rec)
	if s.reinvestDividends {
		bal += div
		div = 0
	} else {
		div += q.Div
	}
	s.last = rec
	return Quote{Bal: bal, Div: div}
}

func (s *Hold) invest(c fin.Cents, curr *pb.Record) fin.Cents {
	ror := SPRateOfReturn(s.last, curr)
	return fin.Cents(math.Floor(float64(c) * float64(ror)))
}

func (s *Hold) payDividend(c fin.Cents, rec *pb.Record) fin.Cents {
	ror := DivRateOfReturn(rec)
	return fin.Cents(math.Floor(float64(c) * float64(ror)))
}
