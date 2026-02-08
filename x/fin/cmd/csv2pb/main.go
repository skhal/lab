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
	"fmt"
	"os"

	"github.com/skhal/lab/x/fin/internal/csv"
	"google.golang.org/protobuf/proto"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(files []string) error {
	if len(files) < 1 {
		return fmt.Errorf("missing csv input file")
	}
	if len(files) != 1 {
		return fmt.Errorf("need only one file")
	}
	return runFile(files[0])
}

func runFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()
	m, err := csv.Read(f)
	if err != nil {
		return err
	}
	b, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(b)
	return err
}
