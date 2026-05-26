// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Irex runs commands for Irrational Exhuberance book.
//
// Synopsis
//
//	irex help
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type command struct {
	run         func([]string) error
	name        string
	description string
}

var commands = []command{
	{
		name:        "import",
		description: "import CSV data",
		run:         new(cmdImport).Run,
	},
	{
		name:        "serve",
		description: "run HTTP server",
		run:         new(cmdServe).Run,
	},
}

func init() {
	flag.CommandLine.Init(filepath.Base(os.Args[0]), flag.ExitOnError)
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "usage: %s command\n", flag.CommandLine.Name())
		fmt.Fprintln(w)
		fmt.Fprintln(w, "commands:")
		for _, cmd := range commands {
			fmt.Fprintf(w, "  %s\t%s\n", cmd.name, cmd.description)
		}
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	name, args := args[0], args[1:]
	cmd, ok := func(name string) (command, bool) {
		for _, cmd := range commands {
			if cmd.name == name {
				return cmd, true
			}
		}
		return command{}, false
	}(name)
	if !ok {
		fmt.Fprintf(os.Stderr, "invalid command: %s\n", name)
		flag.Usage()
		os.Exit(2)
	}
	if err := cmd.run(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
