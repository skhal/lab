// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"iter"
	"os"
	"path/filepath"

	"github.com/skhal/lab/book/ostep/mem/allocator/internal/heap"
	"github.com/skhal/lab/book/ostep/mem/allocator/internal/report"
	"github.com/skhal/lab/go/flags"
)

type command struct {
	heapBase  int
	heapSize  int
	alignment int
	numOps    int
	coalMode  heap.CoalesceMode
	allocMode heap.AllocateMode
}

func newCommand() *command {
	return &command{
		heapBase:  1000,
		heapSize:  1000,
		alignment: 1, // no alignment
		numOps:    5,
		coalMode:  heap.CoalesceModeNoop,
		allocMode: heap.AllocateModeFirstFit,
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
	fs.Var(newAlignmentFlag(&cmd.alignment), "align", "alignment, multiple of 2")
	fs.Var(newBoundedIntFlag(&cmd.numOps, 5, 25), "n", "number of random operations")
	fs.Var(newCoalesceModeFlag(&cmd.coalMode), "c", "coalesce mode")
	fs.Var(newAllocateModeFlag(&cmd.allocMode), "alloc", "allocate mode")
	return fs.ParseAndValidate(os.Args[1:])
}

func (cmd *command) run() error {
	opts := []heap.Option{
		// keep-sorted start
		heap.WithAlignment(cmd.alignment),
		heap.WithAllocator(cmd.allocMode),
		heap.WithCoalesce(cmd.coalMode),
		// keep-sorted end
	}
	h, err := heap.New(cmd.heapBase, cmd.heapSize, opts...)
	if err != nil {
		return err
	}
	sim := newSimulator(h, cmd.numOps)
	return report.Generate(os.Stdout, report.Data{
		Heap: report.Heap{
			Base:      cmd.heapBase,
			Size:      cmd.heapSize,
			CoalMode:  cmd.coalMode.String(),
			AllocMode: cmd.allocMode.String(),
			Blocks:    blocks(h),
		},
		Ops: trace(h, sim),
	})
}

func blocks(h *heap.Heap) []report.Block {
	var bb []report.Block
	h.Walk(func(sz, addr int, fl heap.BlockFlags) {
		b := report.Block{
			Size:      sz,
			Addr:      addr,
			Alloc:     fl.Allocated,
			AllocPrev: fl.AllocatedPrev,
		}
		bb = append(bb, b)
	})
	return bb
}

func trace(h *heap.Heap, sim *simulator) iter.Seq[report.Operation] {
	op := func() report.Operation {
		o := sim.Op()
		err := o.Run()
		return report.Operation{
			Name:      o.String(),
			Err:       err,
			Addresses: sim.Allocated(),
			Blocks:    blocks(h),
		}
	}
	return func(yield func(report.Operation) bool) {
		for sim.Next() && yield(op()) {
			continue
		}
	}
}
