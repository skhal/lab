// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Finsim simulates financial market using Shiller data. It supports different
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
	"google.golang.org/protobuf/proto"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ifile, balance, nmonth, runners, err := parseFlags(createRegistry())
	if err != nil {
		return err
	}
	market, err := readFile(ifile)
	if err != nil {
		return err
	}
	fetchLastN := func(recs []*pb.Record, n int) []*pb.Record {
		n = max(len(recs)-n, 0)
		return recs[n:]
	}
	return runStrategies(runners, balance, fetchLastN(market.GetRecords(), nmonth))
}

func createRegistry() *registry {
	reg := newRegistry()
	mustRegister := func(nr *namedRunner) {
		if err := reg.Register(nr); err != nil {
			panic(err)
		}
	}
	mustRegister(Hold())
	mustRegister(HoldReinvest())
	mustRegister(Retain3Hold())
	mustRegister(Retain4Hold())
	mustRegister(Retain3HoldReinvest())
	mustRegister(Retain4HoldReinvest())
	return reg
}

func parseFlags(reg *registry) (file string, balance fin.Cents, months int, runners []*namedRunner, err error) {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "usage: %s [flags] file\n", filepath.Base(os.Args[0]))
		fmt.Fprintln(w)
		fmt.Fprintln(w, "flags:")
		flag.PrintDefaults()
	}
	bal := flag.Int("bal", 100, "initial balance in dollars")
	flag.IntVar(&months, "n", 12, "number of latest months to process")
	sflag := newStrategyListFlag(reg)
	flag.Var(sflag, "s", sflag.Help())
	flag.Parse()
	if *bal < 0 {
		err = errors.New("negative balance")
		return
	}
	balance = fin.Cents(*bal * 100) // bal is in dollars
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

func runStrategies(strategies []*namedRunner, balance fin.Cents, market []*pb.Record) error {
	infos := make([]*report.StrategyInfo, 0, len(strategies))
	for _, r := range strategies {
		infos = append(infos, r.Run(balance, market))
	}
	return report.Strategies(os.Stdout, infos)
}
