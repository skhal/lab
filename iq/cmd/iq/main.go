// Copyright 2025 Samvel Khalatyan. All rights reserved.

/*
Iq gives access to the interview questions registry.

Synopsis:

	iq

Without arguments, iq dumps a list of questions, sorted by ID.
*/
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/skhal/lab/iq/create"
	"github.com/skhal/lab/iq/info"
	"github.com/skhal/lab/iq/registry"
)

var commands = map[string]command{
	// keep-sorted start
	"create": runCreate,
	"info":   runInfo,
	// keep-sorted end
}

type command func(reg *registry.R, args []string) error

func main() {
	registryConfig := new(registry.Config)
	registryConfig.RegisterFlags(flag.CommandLine)
	flag.Parse()
	if err := run(registryConfig, flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(registryConfig *registry.Config, args []string) error {
	if len(args) == 0 {
		return errors.New("missing command")
	}
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
	return cmd(reg, args)
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
