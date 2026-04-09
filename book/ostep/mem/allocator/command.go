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
	"github.com/skhal/lab/book/ostep/mem/allocator/internal/sim"
	"github.com/skhal/lab/go/flags"
)

type command struct {
	heapBase  int
	heapSize  int
	alignment int
	numOps    int
	coalMode  heap.CoalesceMode
	allocMode heap.AllocateMode
	ops       []string
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
	fs.Var(newBoundedIntFlag(&cmd.numOps, 5, 50), "n", "number of random operations")
	fs.Var(newCoalesceModeFlag(&cmd.coalMode), "c", "coalesce mode")
	fs.Var(newAllocateModeFlag(&cmd.allocMode), "alloc", "allocate mode")
	fs.Var(newOperationListFlag(&cmd.ops), "ops", "list of operations: +N,-N")
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
	var simOps []sim.Option
	if cmd.ops != nil {
		simOps = append(simOps, sim.WithOps(cmd.ops))
	}
	sim := sim.NewSimulator(h, cmd.numOps, simOps...)
	return report.Generate(os.Stdout, report.Data{
		Heap: report.Heap{
			Base:      cmd.heapBase,
			Size:      cmd.heapSize,
			CoalMode:  cmd.coalMode.String(),
			AllocMode: cmd.allocMode.String(),
			Blocks:    blocks(h),
		},
		Trace: runSimulation(h, sim),
	})
}

func blocks(h *heap.Heap) iter.Seq[report.Block] {
	return func(yield func(report.Block) bool) {
		var stop bool
		h.Walk(func(sz, addr int, fl heap.BlockFlags) {
			if stop {
				return
			}
			b := report.Block{
				Size:      sz,
				Addr:      addr,
				Alloc:     fl.Allocated,
				AllocPrev: fl.AllocatedPrev,
			}
			if !yield(b) {
				stop = true
			}
		})
	}
}

func runSimulation(h *heap.Heap, sim *sim.Simulator) iter.Seq[report.Frame] {
	op := func() report.Frame {
		o := sim.Op()
		return &trace{
			name:   o.String(),
			err:    o.Run(),
			addrs:  sim.Allocated,
			blocks: func() iter.Seq[report.Block] { return blocks(h) },
		}
	}
	return func(yield func(report.Frame) bool) {
		for sim.Next() && yield(op()) {
			continue
		}
	}
}

type trace struct {
	name   string
	err    error
	addrs  func() []int
	blocks func() iter.Seq[report.Block]
}

// Operation returns the name of operation.
func (tr *trace) Operation() string {
	return tr.name
}

// Err return operation error, if any.
func (tr *trace) Err() error {
	return tr.err
}

// Addresses returns a list of allocated addresses.
func (tr *trace) Addresses() []int {
	return tr.addrs()
}

// Blocks returns a list of blocks in the heap.
func (tr *trace) Blocks() iter.Seq[report.Block] {
	return tr.blocks()
}
