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

// HoldReinvestDiv implements a strategy to hold investment for the duration of
// market and re-invest dividends.
type HoldReinvestDiv struct{}

// Run executes the strategy.
func (s *HoldReinvestDiv) Run(c fin.Cents, recs []*pb.Record) fin.Cents {
	var prev *pb.Record
	for _, r := range recs {
		c = s.invest(c, prev, r)
		prev = r
	}
	return c
}

func (s *HoldReinvestDiv) invest(c fin.Cents, prev, curr *pb.Record) fin.Cents {
	return s.trade(c, prev, curr) + s.payDividend(c, curr)
}

func (s *HoldReinvestDiv) trade(c fin.Cents, prev, curr *pb.Record) fin.Cents {
	if prev == nil {
		return c
	}
	ror := SPRateOfReturn(prev, curr)
	return fin.Cents(math.Floor(float64(c) * float64(ror)))
}

func (s *HoldReinvestDiv) payDividend(c fin.Cents, rec *pb.Record) fin.Cents {
	ror := DivRateOfReturn(rec)
	return fin.Cents(math.Floor(float64(c) * float64(ror)))
}
