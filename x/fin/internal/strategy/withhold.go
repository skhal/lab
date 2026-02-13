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

// Withhold removes a set fraction of balance every year. For example, a 3%
// withhold takes 3% every year on January. It wraps strategies, i.e. one may
// set withholding on any strategy.
type Withhold struct {
	c    Cycler
	rate float64
}

// Percent is a value in the range [0, 100].
type Percent int

// NewWithhold creates a withholder for a strategy c at annual percent.
func NewWithhold(c Cycler, annual Percent) *Runner {
	h := &Withhold{
		c:    c,
		rate: float64(annual) / 100.,
	}
	return New(h)
}

// Cycle withholds percent at the beginning of the year and runs wrapped
// strategy.
func (s *Withhold) Cycle(q Quote, prev, curr *pb.Record) Quote {
	if s.shouldWithhold(prev, curr) {
		q = s.withhold(q)
	}
	return s.c.Cycle(q, prev, curr)
}

func (s *Withhold) shouldWithhold(prev, curr *pb.Record) bool {
	// the strategy must run for at least one cycle
	if prev == nil {
		return false
	}
	return curr.GetDate().GetMonth() == int32(time.January)
}

func (s *Withhold) withhold(q Quote) Quote {
	applyTo := func(c fin.Cents) fin.Cents {
		return c - fin.Cents(math.Floor(float64(c)*s.rate))
	}
	return Quote{
		Bal: applyTo(q.Bal),
		Div: applyTo(q.Div),
	}
}
