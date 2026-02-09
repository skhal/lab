// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// csv2pb converts Shiller market data from CSV to Protobuf format,
// https://shillerdata.com.
//
// Synopsis:
//
//	csv2pb /tmp/ie_data.csv > /tmp/ie_data.pb
//
// where ie_data.csv is exported "Data" table with S&P Composite index,
// Dividends, and other columns unmodified.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skhal/lab/x/fin/internal/csv"
	"github.com/skhal/lab/x/fin/internal/pb"
	"google.golang.org/protobuf/proto"
)

const permWrite = 0664

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ifile, ofile, err := parseFlags()
	if err != nil {
		return err
	}
	m, err := parseFile(ifile)
	if err != nil {
		return err
	}
	return writeFile(m, ofile)
}

func parseFlags() (ifile, ofile string, err error) {
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "usage: %s [flags] file\n", filepath.Base(os.Args[0]))
		fmt.Fprintln(out)
		fmt.Fprintf(out, "flags:\n")
		flag.PrintDefaults()
	}
	flag.StringVar(&ofile, "o", "", "output file")
	flag.Parse()
	if ofile == "" {
		err = errors.New("missing output file")
		return
	}
	if flag.NArg() != 1 {
		err = errors.New("missing input file")
		return
	}
	ifile = flag.Arg(0)
	return
}

func parseFile(file string) (*pb.Market, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()
	m, err := csv.Read(f)
	if err != nil {
		return nil, fmt.Errorf("%s:%s", file, err)
	}
	return m, nil
}

func writeFile(m *pb.Market, name string) error {
	b, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	return os.WriteFile(name, b, permWrite)
}
