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
	s, e := runHoldStrategy(fetchLastN(m.GetRecords(), nmonth))
	report(s, e)
	return nil
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

func runHoldStrategy(recs []*pb.Record) (start, end sim.Quote) {
	return sim.Run(fin.Cents(100), recs, new(strategy.HoldReinvestDiv))
}

func report(s, e sim.Quote) {
	fmt.Printf("%s\n", s)
	fmt.Printf("%s\n", e)
	fmt.Printf("return: %.2f\n", float64(e.Balance)/float64(s.Balance))
}
