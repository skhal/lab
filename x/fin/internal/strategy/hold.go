// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strategy

import (
	"math"

	"github.com/skhal/lab/x/fin/internal/fin"
	"github.com/skhal/lab/x/fin/internal/pb"
	"github.com/skhal/lab/x/fin/internal/ror"
)

type hold struct {
	last *pb.Record
}

// Hold createsa a hold strategy.
func Hold() *Runner {
	return New(&hold{})
}

// Cycle implements [Cycler] interface.
func (s *hold) Cycle(pos fin.Position, rec *pb.Record) fin.Position {
	defer func() { s.last = rec }()
	bal := s.invest(pos.Investment, rec)
	div := s.payDividend(pos.Investment, rec) + pos.Dividend
	return fin.Position{Investment: bal, Dividend: div}
}

func (s *hold) invest(c fin.Cents, curr *pb.Record) fin.Cents {
	r := ror.SPComposite(s.last, curr)
	return fin.Cents(math.Floor(float64(c) * float64(r)))
}

func (s *hold) payDividend(c fin.Cents, rec *pb.Record) fin.Cents {
	r := ror.Dividend(rec)
	return fin.Cents(math.Floor(float64(c) * float64(r)))
}

type holdReinvest struct {
	hold *hold
}

// HoldReinvest creates a hold strategy with dividend reinvestment.
func HoldReinvest() *Runner {
	s := &holdReinvest{new(hold)}
	return New(s)
}

// Cycle implements [Cycler] interface.
func (s *holdReinvest) Cycle(pos fin.Position, rec *pb.Record) fin.Position {
	pos = s.hold.Cycle(pos, rec)
	pos.Investment += pos.Dividend
	pos.Dividend = 0
	return pos
}
