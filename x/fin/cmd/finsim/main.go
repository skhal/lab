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
	var cmd simCommand
	if err := cmd.Run(NewRegistry()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type simCommand struct {
	file string
	data []*pb.Record

	balance fin.Cents
	months  int

	runners []*namedRunner
}

// Run executes the command.
func (cmd *simCommand) Run(reg *registry) error {
	if err := cmd.parseFlags(reg); err != nil {
		return err
	}
	if err := cmd.loadData(); err != nil {
		return err
	}
	return cmd.runStrategies()
}

func (cmd *simCommand) parseFlags(reg *registry) error {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "usage: %s [flags] file\n", filepath.Base(os.Args[0]))
		fmt.Fprintln(w)
		fmt.Fprintln(w, "flags:")
		flag.PrintDefaults()
	}
	bal := flag.Int("bal", 100, "initial balance in dollars")
	flag.IntVar(&cmd.months, "n", 12, "number of latest months to process")
	sflag := newStrategyListFlag(reg)
	flag.Var(sflag, "s", sflag.Help())
	flag.Parse()
	if flag.NArg() != 1 {
		return errors.New("missing input file")
	}
	if *bal < 0 {
		return errors.New("negative balance")
	}
	cmd.balance = fin.Cents(*bal * 100) // bal is in dollars
	cmd.file = flag.Arg(0)
	cmd.runners = sflag.Runners()
	return nil
}

func (cmd *simCommand) loadData() error {
	b, err := os.ReadFile(cmd.file)
	if err != nil {
		return err
	}
	var m = new(pb.Market)
	if err := proto.Unmarshal(b, m); err != nil {
		return err
	}
	fetchLastN := func(recs []*pb.Record, n int) []*pb.Record {
		n = max(len(recs)-n, 0)
		return recs[n:]
	}
	cmd.data = fetchLastN(m.GetRecords(), cmd.months)
	return nil
}

func (cmd *simCommand) runStrategies() error {
	infos := make([]*report.StrategyInfo, 0, len(cmd.runners))
	for _, r := range cmd.runners {
		infos = append(infos, r.Run(cmd.balance, cmd.data))
	}
	return report.Strategies(os.Stdout, infos)
}
