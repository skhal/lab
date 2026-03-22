// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/skhal/lab/book/ostep/mem/allocator/internal/heap"
	"github.com/skhal/lab/book/ostep/mem/allocator/internal/report"
	"github.com/skhal/lab/go/flags"
)

type command struct {
	heapBase int
	heapSize int
}

func newCommand() *command {
	return &command{
		heapBase: 1000,
		heapSize: 1000,
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
	fs.Var(newBoundedIntFlag(&cmd.heapBase, 0, 10000), "base", "heap base address")
	fs.Var(newBoundedIntFlag(&cmd.heapSize, 100, 10000), "size", "heab size")
	return fs.ParseAndValidate(os.Args[1:])
}

func (cmd *command) run() error {
	h := heap.New(cmd.heapBase, cmd.heapSize)
	return report.Generate(os.Stdout, report.Data{
		Heap: report.Heap{
			Base: cmd.heapBase,
			Size: cmd.heapSize,
			Free: freeBlocks(h),
		},
	})
}

func freeBlocks(h *heap.Heap) []report.Block {
	var bb []report.Block
	h.WalkFreeSpace(func(sz, addr int) bool {
		bb = append(bb, report.Block{
			Size: sz,
			Addr: addr,
		})
		return true
	})
	return bb
}
