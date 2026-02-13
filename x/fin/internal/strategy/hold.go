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

type hold struct {
	reinvestDiv bool
	last        *pb.Record
}

// Hold createsa a hold strategy.
func Hold() *Runner {
	return New(&hold{})
}

// HoldReinvest creates a hold strategy with dividend reinvestment.
func HoldReinvest() *Runner {
	return New(&hold{reinvestDiv: true})
}

// Cycle executes a single cycle of the hold strategy.
func (s *hold) Cycle(q Quote, rec *pb.Record) Quote {
	bal := s.invest(q.Bal, rec)
	div := s.payDividend(q.Bal, rec)
	if s.reinvestDiv {
		bal += div
		div = 0
	} else {
		div += q.Div
	}
	s.last = rec
	return Quote{Bal: bal, Div: div}
}

func (s *hold) invest(c fin.Cents, curr *pb.Record) fin.Cents {
	ror := SPRateOfReturn(s.last, curr)
	return fin.Cents(math.Floor(float64(c) * float64(ror)))
}

func (s *hold) payDividend(c fin.Cents, rec *pb.Record) fin.Cents {
	ror := DivRateOfReturn(rec)
	return fin.Cents(math.Floor(float64(c) * float64(ror)))
}
