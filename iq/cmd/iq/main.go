// Copyright 2025 Samvel Khalatyan. All rights reserved.

/*
Iq gives access to the interview questions registry.

Synopsis:

	iq

Without arguments, iq dumps a list of questions, sorted by ID.
*/
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/skhal/lab/iq/create"
	"github.com/skhal/lab/iq/info"
	"github.com/skhal/lab/iq/registry"
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
			fmt.Fprintf(buf, "Usage: %s [-f <file>] <command> [<args>]\n", os.Args[0])
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
	cfg := new(create.Config)
	cfg.RegisterFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}
	return create.Run(cfg, reg)
}

func runInfo(reg *registry.R, args []string) error {
	return info.Run(reg, args...)
}
