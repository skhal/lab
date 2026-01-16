// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Iq gives access to the interview questions registry.

Synopsis:

	iq

Without arguments, iq dumps a list of questions, sorted by ID.
*/
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/skhal/lab/iq/cmd/iq/internal/create"
	"github.com/skhal/lab/iq/cmd/iq/internal/info"
	"github.com/skhal/lab/iq/cmd/iq/internal/registry"
)

var commands = map[string]command{
	// keep-sorted start block=yes
	"create": {
		Desc: "create a question",
		Run:  runCreate,
	},
	"info": {
		Desc: "display information for questions",
		Run:  runInfo,
	},
	// keep-sorted end
}

type command struct {
	Desc string
	Run  func(reg *registry.R, args []string) error
}

func init() {
	flag.Usage = func() {
		header := func() string {
			buf := new(bytes.Buffer)
			fmt.Fprintf(buf, "Usage: %s [-f <file>] <command> [<args>]\n", flag.CommandLine.Name())
			fmt.Fprintln(buf)
			fmt.Fprintln(buf, "Commands:")
			for name, cmd := range commands {
				fmt.Fprintf(buf, "  %-10s%s\n", name, cmd.Desc)
			}
			fmt.Fprintln(buf)
			fmt.Fprintln(buf, "Common flags:")
			return buf.String()
		}
		fmt.Fprint(flag.CommandLine.Output(), header())
		flag.PrintDefaults()
	}
}

func parseFlags() {
	flag.Parse()
	if args := flag.Args(); len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	registryConfig := new(registry.Config)
	registryConfig.RegisterFlags(flag.CommandLine)
	parseFlags()
	if err := run(registryConfig, flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(registryConfig *registry.Config, args []string) error {
	reg, err := registry.Load(registryConfig)
	if err != nil {
		return err
	}
	cmdName := args[0]
	cmdArgs := args[1:]
	if err := runCommand(reg, cmdName, cmdArgs); err != nil {
		return fmt.Errorf("%s: %v", cmdName, err)
	}
	if reg.Updated() {
		return registry.Write(reg, registryConfig)
	}
	return nil
}

func runCommand(reg *registry.R, name string, args []string) error {
	cmd, ok := commands[name]
	if !ok {
		return fmt.Errorf("invalid command -- %s", name)
	}
	return cmd.Run(reg, args)
}

func runCreate(reg *registry.R, args []string) error {
	fs := flag.NewFlagSet("iq create", flag.ContinueOnError)
	fs.Usage = func() {
		header := func() string {
			buf := new(bytes.Buffer)
			fmt.Fprintf(buf, "Usage: %s <args>\n", fs.Name())
			fmt.Fprintln(buf)
			fmt.Fprintln(buf, "Arguments:")
			return buf.String()
		}
		fmt.Fprint(fs.Output(), header())
		fs.PrintDefaults()
	}
	cfg := new(create.Config)
	cfg.RegisterFlags(fs)
	if err := fs.Parse(args); errors.Is(err, flag.ErrHelp) {
		return nil
	} else if err != nil {
		return err
	}
	return create.Run(cfg, reg)
}

func runInfo(reg *registry.R, args []string) error {
	fs := flag.NewFlagSet("iq info", flag.ContinueOnError)
	fs.Usage = func() {
		header := func() string {
			buf := new(bytes.Buffer)
			fmt.Fprintf(buf, "Usage: %s <args>\n", fs.Name())
			fmt.Fprintln(buf)
			fmt.Fprintln(buf, "Arguments:")
			return buf.String()
		}
		fmt.Fprint(fs.Output(), header())
		fs.PrintDefaults()
	}
	cfg := new(info.Config)
	cfg.RegisterFlags(fs)
	if err := fs.Parse(args); errors.Is(err, flag.ErrHelp) {
		return nil
	} else if err != nil {
		return err
	}
	return info.Run(cfg, reg, fs.Args()...)
}
