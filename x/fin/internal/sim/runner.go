// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sim

import (
	"fmt"
	"iter"

	"github.com/skhal/lab/x/fin/internal/fin"
	"github.com/skhal/lab/x/fin/internal/pb"
	"github.com/skhal/lab/x/fin/internal/stat"
	"github.com/skhal/lab/x/fin/internal/strategy"
)

// Result is a output of a simulation. It holds the beginning balance and
// descriptive statistics about the end sample.
type Result struct {
	// Start is the beginning balance.
	Start fin.Cents
	// End holds descriptive statistics of the end balance of a sample of
	// multiple cycles, run through simulation.
	End stat.Description
}

type runner struct {
	rfs            []strategy.RebalanceFunc
	numCycles      int
	cycleLenMonths int
}

// NewRunner creates a runner to run simulations with cycles of given length
// in months and rebalancers.
func NewRunner(cycles, cycleLenMonths int, rfs []strategy.RebalanceFunc) *runner {
	return &runner{
		rfs:            rfs,
		numCycles:      cycles,
		cycleLenMonths: cycleLenMonths,
	}
}

// Run executes a simulation on data with starting balance of cash. It runs
// multiple simulations, each covering a cycle of market data, up to
// [runner.numCycles].
func (r *runner) Run(cash fin.Cents, data []*pb.Record) *Result {
	endCash := func(bals []fin.Balance) fin.Cents {
		lastBal := bals[len(bals)-1]
		return lastBal.Cash
	}
	var end []fin.Cents
	balsNum := 0
	for cycle := range Cycles(data, r.numCycles, r.cycleLenMonths) {
		bals := strategy.Drive(cash, cycle, r.rfs...)
		if balsNum == 0 {
			balsNum = len(bals)
		} else {
			if len(bals) != balsNum {
				panic(fmt.Errorf("unexpected number of balances %d, want %d", len(bals), balsNum))
			}
		}
		end = append(end, endCash(bals))
	}
	desc := stat.Describe(end)
	return &Result{
		Start: cash,
		End:   desc,
	}
}

const yearMonths = 12

// Cycles generates num cycles of length size from data, moving from the end:
//
//	data[size:], data[size-12:len(data)-12], etc.
func Cycles(data []*pb.Record, num, size int) iter.Seq[[]*pb.Record] {
	subtractOneYear := func(d []*pb.Record) []*pb.Record {
		if len(d) < yearMonths {
			return nil
		}
		return d[:len(d)-yearMonths]
	}
	cycle := func(d []*pb.Record, months int) []*pb.Record {
		if len(d) < months {
			return nil
		}
		return d[len(d)-months:]
	}
	return func(yield func([]*pb.Record) bool) {
		for c := 0; c < num; c += 1 {
			cd := cycle(data, size)
			if cd == nil {
				break
			}
			if !yield(cd) {
				break
			}
			data = subtractOneYear(data)
		}
	}
}
