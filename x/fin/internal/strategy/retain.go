// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strategy

import (
	"math"
	"time"

	"github.com/skhal/lab/x/fin/internal/fin"
	"github.com/skhal/lab/x/fin/internal/pb"
)

// retain withholds [retain.rate] value every year from underlying strategy
// [retain.c] in January.
type retain struct {
	c    Cycler
	rate float64
	last *pb.Record
}

// Percent is a value in the range [0, 100].
type Percent int

// Retain creates a withholder for a strategy c at annual percent.
func Retain(p Percent, c Cycler) *Runner {
	h := &retain{
		c:    c,
		rate: float64(p) / 100.,
	}
	return New(h)
}

// Cycle retains percent at the beginning of the year and runs wrapped
// strategy.
func (s *retain) Cycle(pos fin.Position, rec *pb.Record) fin.Position {
	defer func() { s.last = rec }()
	if s.shouldRetain(rec) {
		pos = s.withhold(pos)
	}
	return s.c.Cycle(pos, rec)
}

func (s *retain) shouldRetain(rec *pb.Record) bool {
	// the strategy must run for at least one cycle
	if s.last == nil {
		return false
	}
	return rec.GetDate().GetMonth() == int32(time.January)
}

func (s *retain) withhold(pos fin.Position) fin.Position {
	applyTo := func(c fin.Cents) fin.Cents {
		return c - fin.Cents(math.Floor(float64(c)*s.rate))
	}
	return fin.Position{
		Investment: applyTo(pos.Investment),
		Dividend:   applyTo(pos.Dividend),
	}
}
