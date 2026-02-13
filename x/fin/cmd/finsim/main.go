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
	"strings"

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
	reg := createRegistry()
	ifile, nmonth, runners, err := parseFlags(reg)
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
	return runStrategies(runners, fetchLastN(m.GetRecords(), nmonth))
}

func createRegistry() *registry {
	reg := newRegistry()
	mustRegister := func(nr *namedRunner) {
		if err := reg.Register(nr); err != nil {
			panic(err)
		}
	}
	mustRegister(HoldCoollect())
	mustRegister(HoldReinvest())
	mustRegister(Retain3HoldCollect())
	mustRegister(Retain4HoldCollect())
	mustRegister(Retain3HoldReinvest())
	mustRegister(Retain4HoldReinvest())
	return reg
}

func parseFlags(reg *registry) (file string, months int, runners []*namedRunner, err error) {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "usage: %s [flags] file\n", filepath.Base(os.Args[0]))
		fmt.Fprintln(w)
		fmt.Fprintln(w, "flags:")
		flag.PrintDefaults()
	}
	flag.IntVar(&months, "n", 12, "number of latest months to process")
	sflag := newStrategyListFlag(reg)
	flag.Var(sflag, "s", sflag.Help())
	flag.Parse()
	if flag.NArg() != 1 {
		err = errors.New("missing input file")
		return
	}
	file = flag.Arg(0)
	runners = sflag.Runners()
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

func runStrategies(strategies []*namedRunner, market []*pb.Record) error {
	infos := make([]*report.StrategyInfo, 0, len(strategies))
	for _, r := range strategies {
		infos = append(infos, r.Run(market))
	}
	return report.Strategies(os.Stdout, infos)
}

type registry struct {
	runners map[string]*namedRunner
}

func newRegistry() *registry {
	return &registry{make(map[string]*namedRunner)}
}

// Get retrieves a strategy runner from the registry. It returns a boolean flag
// to indicate whether a runner with a given name is available.
func (reg *registry) Get(name string) (*namedRunner, bool) {
	r, ok := reg.runners[name]
	return r, ok
}

// Len returns the number of registered runners.
func (reg *registry) Len() int {
	return len(reg.runners)
}

// Register adds a strategy runner to the registry.
func (reg *registry) Register(r *namedRunner) error {
	if _, ok := reg.runners[r.Name()]; ok {
		return fmt.Errorf("duplicate runner %s", r.Name())
	}
	reg.runners[r.Name()] = r
	return nil
}

// Walk applies f to every registered strategy. The callback may return false
// to stop the iteration short.
func (reg *registry) Walk(f func(*namedRunner) bool) {
	for _, r := range reg.runners {
		if !f(r) {
			break
		}
	}
}

type namedRunner struct {
	name   string
	desc   string
	runner *strategy.Runner
}

// Name returns the strategy name.
func (nr *namedRunner) Name() string { return nr.name }

// Description gices a strategy description.
func (nr *namedRunner) Description() string { return nr.desc }

// Run executes strategy.
func (nr *namedRunner) Run(market []*pb.Record) *report.StrategyInfo {
	info := report.StrategyInfo{
		Name:        nr.Name(),
		Description: nr.Description(),
	}
	info.Start, info.End = sim.Run(fin.Cents(100), market, nr.runner)
	return &info
}

type strategyListFlag struct {
	reg     *registry
	runners []*namedRunner

	seen map[string]bool
	set  bool
}

func newStrategyListFlag(reg *registry) *strategyListFlag {
	runners := make([]*namedRunner, 0, reg.Len())
	reg.Walk(func(r *namedRunner) bool {
		runners = append(runners, r)
		return true
	})
	return &strategyListFlag{
		reg:     reg,
		runners: runners,
		seen:    make(map[string]bool),
	}
}

// Help generates a help message for the flag.
func (f *strategyListFlag) Help() string {
	names := make([]string, 0, f.reg.Len())
	f.reg.Walk(func(r *namedRunner) bool {
		names = append(names, r.Name())
		return true
	})
	opts := strings.Join(names, "\n")
	return fmt.Sprintf("comma-separated list of strategies to run:\n%s\n", opts)
}

// Runners returns a list of registered runners.
func (f *strategyListFlag) Runners() []*namedRunner {
	return f.runners
}

// Set implements flag.Value interface.
func (f *strategyListFlag) Set(s string) error {
	var runners []*namedRunner
	for name := range strings.SplitSeq(s, ",") {
		r, ok := f.reg.Get(name)
		if !ok {
			return fmt.Errorf("unsupported strategy %s", name)
		}
		if f.seen[name] {
			return fmt.Errorf("duplicate strategy %s", name)
		}
		f.seen[name] = true
		runners = append(runners, r)
	}
	if !f.set {
		f.set = true
		f.runners = f.runners[:0]
	}
	f.runners = append(f.runners, runners...)
	return nil
}

// String implements flag.Valaue interface.
func (f *strategyListFlag) String() string {
	names := make([]string, 0, len(f.runners))
	for _, r := range f.runners {
		names = append(names, r.Name())
	}
	return strings.Join(names, ",")
}

// HoldCoollectDiv creates a strategy to hold SP composite index and collect
// dividends.
func HoldCoollect() *namedRunner {
	return &namedRunner{
		name:   "hold-collect-div",
		desc:   "hold s&p, collect dividends",
		runner: strategy.Hold(),
	}
}

// HoldReinvestDiv creates a strategy to hold SP composite index and reinvest
// dividend payouts into the index.
func HoldReinvest() *namedRunner {
	return &namedRunner{
		name:   "hold-reinvest-div",
		desc:   "hold s&p, reinvest dividends",
		runner: strategy.HoldReinvest(),
	}
}

// Retain3HoldCollectDiv creates a strategy to retain 3% every year from
// [HoldCollectDiv] strategy.
func Retain3HoldCollect() *namedRunner {
	return &namedRunner{
		name:   "retain-3-hold-collect-div",
		desc:   "retain 3% yearly, hold s&p, collect dividends",
		runner: strategy.Retain(strategy.Percent(3), strategy.Hold()),
	}
}

// Retain4HoldCollectDiv creates a strategy to retain 4% every year from
// [HoldCollectDiv] strategy.
func Retain4HoldCollect() *namedRunner {
	return &namedRunner{
		name:   "retain-4-hold-collect-div",
		desc:   "retain 4% yearly, hold s&p, collect dividends",
		runner: strategy.Retain(strategy.Percent(4), strategy.Hold()),
	}
}

// Retain3HoldReinvestDiv creates a strategy to retain 3% every year from
// [HoldReinvestDiv] strategy.
func Retain3HoldReinvest() *namedRunner {
	return &namedRunner{
		name:   "retain-3-hold-reinvest-div",
		desc:   "retain 3% yearly, hold s&p, reinvest dividends",
		runner: strategy.Retain(strategy.Percent(3), strategy.HoldReinvest()),
	}
}

// Retain4HoldReinvestDiv creates a strategy to retain 4% every year from
// [HoldReinvestDiv] strategy.
func Retain4HoldReinvest() *namedRunner {
	return &namedRunner{
		name:   "retain-4-hold-reinvest-div",
		desc:   "retain 4% yearly, hold s&p, reinvest dividends",
		runner: strategy.Retain(strategy.Percent(4), strategy.HoldReinvest()),
	}
}
