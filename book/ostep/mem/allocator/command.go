// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/skhal/lab/book/ostep/mem/allocator/internal/report"
	"github.com/skhal/lab/go/flags"
)

type heap struct {
	base int
	size int
}

type command struct {
	heap heap
}

func newCommand() *command {
	return &command{
		heap: heap{
			base: 1000,
			size: 1000,
		},
	}
}

// Run parses flags and generates the allocation simulation report. It returns
// an error if flag parser or simulation fail.
func (cmd *command) Run() error {
	if err := cmd.parseFlags(); err != nil {
		return err
	}
	return cmd.run()
}

func (cmd *command) parseFlags() error {
	fs := flags.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)
	fs.Var(newBoundedIntFlag(&cmd.heap.base, 0, 10000), "base", "heap base address")
	fs.Var(newBoundedIntFlag(&cmd.heap.size, 100, 10000), "size", "heab size")
	return fs.ParseAndValidate(os.Args[1:])
}

func (cmd *command) run() error {
	return report.Generate(os.Stdout, report.Data{
		Heap: report.Heap{
			Base: cmd.heap.base,
			Size: cmd.heap.size,
		},
	})
}
