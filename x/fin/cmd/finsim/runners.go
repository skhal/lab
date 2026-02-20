// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/skhal/lab/x/fin/internal/fin"
	"github.com/skhal/lab/x/fin/internal/pb"
	"github.com/skhal/lab/x/fin/internal/report"
	"github.com/skhal/lab/x/fin/internal/strategy"
)

type namedRunner struct {
	name        string
	desc        string
	rebalancers []strategy.RebalanceFunc
}

// Name returns the strategy name.
func (nr *namedRunner) Name() string { return nr.name }

// Description gices a strategy description.
func (nr *namedRunner) Description() string { return nr.desc }

// Run executes strategy.
func (nr *namedRunner) Run(bal fin.Cents, market []*pb.Record) *report.StrategyInfo {
	info := report.StrategyInfo{
		Name:        nr.Name(),
		Description: nr.Description(),
	}
	bals := strategy.Drive(bal, market, nr.rebalancers...)
	info.Start = bals[0]
	info.End = bals[len(bals)-1]
	return &info
}

// HoldDiv creates a strategy to hold SP composite index and collect
// dividends.
func Hold() *namedRunner {
	return &namedRunner{
		name: "hold",
		desc: "hold s&p, collect dividends",
	}
}

// HoldReinvestDiv creates a strategy to hold SP composite index and reinvest
// dividend payouts into the index.
func HoldReinvest() *namedRunner {
	return &namedRunner{
		name: "hold-reinvest",
		desc: "hold s&p, reinvest dividends",
		rebalancers: []strategy.RebalanceFunc{
			strategy.ReinvestDividend,
		},
	}
}

// Retain3HoldDiv creates a strategy to retain 3% every year from
// [HoldDiv] strategy.
func Retain3Hold() *namedRunner {
	return &namedRunner{
		name: "retain-3-hold",
		desc: "retain 3% yearly, hold s&p, collect dividends",
		rebalancers: []strategy.RebalanceFunc{
			func(rate float64) strategy.RebalanceFunc {
				yw := strategy.YearlyWithdrawer{Rate: rate}
				return yw.Rebalance
			}(0.03),
		},
	}
}

// Retain4HoldDiv creates a strategy to retain 4% every year from
// [HoldDiv] strategy.
func Retain4Hold() *namedRunner {
	return &namedRunner{
		name: "retain-4-hold",
		desc: "retain 4% yearly, hold s&p, collect dividends",
		rebalancers: []strategy.RebalanceFunc{
			func(rate float64) strategy.RebalanceFunc {
				yw := strategy.YearlyWithdrawer{Rate: rate}
				return yw.Rebalance
			}(0.04),
		},
	}
}

// Retain3HoldReinvestDiv creates a strategy to retain 3% every year from
// [HoldReinvestDiv] strategy.
func Retain3HoldReinvest() *namedRunner {
	return &namedRunner{
		name: "retain-3-hold-reinvest",
		desc: "retain 3% yearly, hold s&p, reinvest dividends",
		rebalancers: []strategy.RebalanceFunc{
			strategy.ReinvestDividend,
			func(rate float64) strategy.RebalanceFunc {
				yw := strategy.YearlyWithdrawer{Rate: rate}
				return yw.Rebalance
			}(0.03),
		},
	}
}

// Retain4HoldReinvestDiv creates a strategy to retain 4% every year from
// [HoldReinvestDiv] strategy.
func Retain4HoldReinvest() *namedRunner {
	return &namedRunner{
		name: "retain-4-hold-reinvest",
		desc: "retain 4% yearly, hold s&p, reinvest dividends",
		rebalancers: []strategy.RebalanceFunc{
			strategy.ReinvestDividend,
			func(rate float64) strategy.RebalanceFunc {
				yw := strategy.YearlyWithdrawer{Rate: rate}
				return yw.Rebalance
			}(0.04),
		},
	}
}
