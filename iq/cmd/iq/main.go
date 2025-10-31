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

	"github.com/skhal/lab/iq/info"
	"github.com/skhal/lab/iq/registry"
)

var commands = map[string]command{
	"info": func(r *registry.R) error { return info.Run(r) },
}

type command func(*registry.R) error

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
	return runCommand(reg, args)
}

func runCommand(reg *registry.R, args []string) error {
	cmdName := args[0]
	cmd, ok := commands[cmdName]
	if !ok {
		return fmt.Errorf("invalid command -- %s", cmdName)
	}
	return cmd(reg)
}
