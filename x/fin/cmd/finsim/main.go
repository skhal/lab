// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// finsim simulates financial market using Shiller data. It supports different
// strategies, e.g. hold the investment position and re-invest dividends.
//
// Synopsis:
//
//	finsim [-n months] data.pb
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skhal/lab/x/fin/internal/fin"
	"github.com/skhal/lab/x/fin/internal/pb"
	"github.com/skhal/lab/x/fin/internal/report"
	"github.com/skhal/lab/x/fin/internal/sim"
	"github.com/skhal/lab/x/fin/internal/strategy"
	"google.golang.org/protobuf/proto"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ifile, nmonth, err := parseFlags()
	if err != nil {
		return err
	}
	m, err := readFile(ifile)
	if err != nil {
		return err
	}
	fetchLastN := func(recs []*pb.Record, n int) []*pb.Record {
		n = max(len(recs)-n, 0)
		return recs[n:]
	}
	return runStrategies(fetchLastN(m.GetRecords(), nmonth))
}

func parseFlags() (file string, months int, err error) {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "usage: %s [flags] file\n", filepath.Base(os.Args[0]))
		fmt.Fprintln(w)
		fmt.Fprintln(w, "flags:")
		flag.PrintDefaults()
	}
	flag.IntVar(&months, "n", 12, "number of latest months to process")
	flag.Parse()
	if flag.NArg() != 1 {
		err = errors.New("missing input file")
		return
	}
	file = flag.Arg(0)
	return
}

func readFile(name string) (*pb.Market, error) {
	b, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	var m = new(pb.Market)
	if err := proto.Unmarshal(b, m); err != nil {
		return nil, err
	}
	return m, nil
}

func runStrategies(market []*pb.Record) error {
	infos := make([]*report.StrategyInfo, 0, len(strategies))
	for name, r := range strategies {
		fmt.Println("run: ", name)
		infos = append(infos, r.Run(market))
	}
	fmt.Println()
	return report.Strategies(os.Stdout, infos)
}

var strategies = make(map[string]*strategyRunner)

func init() {
	register := func(name, desc string, s sim.Strategy) {
		if r, ok := strategies[name]; ok {
			err := fmt.Errorf("strategy with name %s already exists: %s", name, r.Description())
			panic(err)
		}
		strategies[name] = newStrategyRunner(name, desc, s)
	}
	register("hold-collect-div", "hold s&p, collect dividends", strategy.NewHold())
	register("hold-reinvest-div", "hold s&p, reinvest dividends", strategy.NewHold(strategy.HoldOptReinvestDiv()))
}

type strategyRunner struct {
	name     string
	desc     string
	strategy sim.Strategy
}

func newStrategyRunner(name, desc string, s sim.Strategy) *strategyRunner {
	return &strategyRunner{
		name:     name,
		desc:     desc,
		strategy: s,
	}
}

// Name returns the strategy name.
func (r *strategyRunner) Name() string { return r.name }

// Description gices a strategy description.
func (r *strategyRunner) Description() string { return r.desc }

// Run executes strategy.
func (r *strategyRunner) Run(market []*pb.Record) *report.StrategyInfo {
	info := report.StrategyInfo{
		Name:        r.Name(),
		Description: r.Description(),
	}
	info.Start, info.End = sim.Run(fin.Cents(100), market, r.strategy)
	return &info
}
