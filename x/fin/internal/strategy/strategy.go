// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strategy

import (
	"github.com/skhal/lab/x/fin/internal/fin"
	"github.com/skhal/lab/x/fin/internal/pb"
)

// Cycler is a strategy that can process a single cycle.
type Cycler interface {
	// Cycle processes a single market record.
	Cycle(fin.Position, *pb.Record) fin.Position
}

// Runner represents a strategy backed by Cycler.
type Runner struct {
	Cycler
}

// New createas a strategy, backed by Cycler.
func New(c Cycler) *Runner {
	return &Runner{c}
}

// Run executes the strategy on a set of rectors starting with a given balance.
// It returns the end balance.
func (s *Runner) Run(pos fin.Position, market []*pb.Record) fin.Position {
	for _, rec := range market {
		pos = s.Cycle(pos, rec)
	}
	return pos
}
