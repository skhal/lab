// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/skhal/lab/book/irex/csvimport"
	"github.com/skhal/lab/book/irex/pb"
	"google.golang.org/protobuf/proto"
)

type cmdImport struct {
	fs        *flag.FlagSet
	output    string
	skipLines int
	scanLines int
}

// Run imports CSV data and writes it in binary format to specified file.
func (cmd *cmdImport) Run(args []string) error {
	cmd.init()
	if err := cmd.fs.Parse(args); err != nil {
		return err
	}
	cmdArgs := cmd.fs.Args()
	if err := cmd.validate(cmdArgs); err != nil {
		return err
	}
	return cmd.run(cmdArgs[0])
}

func (cmd *cmdImport) init() {
	name := fmt.Sprintf("%s import", flag.CommandLine.Name())
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	fs.Usage = func() {
		w := fs.Output()
		fmt.Fprintf(w, "usage: %s [-f file] file\n", fs.Name())
		fmt.Fprintln(w)
		fmt.Fprintln(w, "flags:")
		fs.PrintDefaults()
	}
	fs.StringVar(&cmd.output, "f", "", "output file")
	fs.IntVar(&cmd.skipLines, "lines-skip", 8, "number of lines to skip")
	fs.IntVar(&cmd.scanLines, "lines-scan", 0, "number of lines to scan including skip lines (0 no limit)")
	cmd.fs = fs
}

func (cmd *cmdImport) validate(args []string) error {
	if cmd.output == "" {
		err := fmt.Errorf("output is not set")
		return newUsageError(cmd.fs, err)
	}
	if cmd.skipLines < 0 {
		return fmt.Errorf("negative skip lines")
	}
	if cmd.scanLines < 0 {
		return fmt.Errorf("negative scan lines")
	}
	if len(args) == 0 {
		err := fmt.Errorf("missing input file")
		return newUsageError(cmd.fs, err)
	}
	if len(args) > 1 {
		err := fmt.Errorf("want one input file")
		return newUsageError(cmd.fs, err)
	}
	return nil
}

func (cmd *cmdImport) run(file string) error {
	quotes, err := cmd.read(file)
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	return cmd.write(quotes)
}

// read imports CSV data from a file.
func (cmd *cmdImport) read(file string) ([]*pb.Quote, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	opts := []csvimport.Opt{
		csvimport.WithSkipLines(cmd.skipLines),
		csvimport.WithScanLines(cmd.scanLines),
	}
	qq, err := csvimport.Import(csv.NewReader(f), opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", file, err)
	}
	return qq, nil
}

// write serializes quotes and writes the result to the output file. If the
// output is "-" (dash) then it writes data to standard output.
func (cmd *cmdImport) write(quotes []*pb.Quote) error {
	m := pb.Market_builder{Quotes: quotes}.Build()
	f, err := os.OpenFile(cmd.output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	for len(b) != 0 {
		n, err := f.Write(b)
		if err != nil {
			return err
		}
		b = b[n:]
	}
	return nil
}

type usageError struct {
	fs  *flag.FlagSet
	err error
}

func newUsageError(fs *flag.FlagSet, err error) *usageError {
	return &usageError{fs: fs, err: err}
}

// Error implements [builtin.error] interface.
func (e *usageError) Error() string {
	var s strings.Builder
	if e.err != nil {
		fmt.Fprintln(&s, e.err)
		fmt.Fprintln(&s)
	}
	if e.fs != nil {
		w := e.fs.Output()
		defer func() { e.fs.SetOutput(w) }()
		e.fs.SetOutput(&s)
		e.fs.Usage()
	}
	return s.String()
}
