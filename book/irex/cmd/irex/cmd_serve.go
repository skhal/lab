// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/skhal/lab/book/irex/serv"
)

const defaultAddress = ":8080"

type cmdServe struct {
	addr string
}

// Run parses flags and runs the web server.
func (cmd *cmdServe) Run(args []string) error {
	if err := cmd.parseFlags(args); err != nil {
		return err
	}
	return cmd.run()
}

func (cmd *cmdServe) parseFlags(args []string) error {
	name := fmt.Sprintf("%s serve", flag.CommandLine.Name())
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	fs.Usage = func() {
		w := fs.Output()
		fmt.Fprintf(w, "usage: %s [flags]\n", fs.Name())
		fmt.Fprintln(w)
		fmt.Fprintln(w, "flags:")
		fs.PrintDefaults()
	}
	fs.StringVar(&cmd.addr, "http", defaultAddress, "address to bind")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(fs.Args()) != 0 {
		err := fmt.Errorf("unexpected arguments: %s", strings.Join(fs.Args(), " "))
		return newUsageError(fs, err)
	}
	return nil
}

func (cmd *cmdServe) run() error {
	fmt.Printf("serve on %s\n", cmd.addr)
	s := &serv.Server{Address: cmd.addr}
	return s.Run()
}
