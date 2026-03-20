// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"iter"
	"math/rand/v2"
	"os"
	"path/filepath"

	"github.com/skhal/lab/book/ostep/mem/segment/internal/mem"
	"github.com/skhal/lab/book/ostep/mem/segment/internal/report"
	"github.com/skhal/lab/go/flags"
)

const (
	maxVirtAddr = mem.MaxVirtAddress
	maxSegSize  = 2 * mem.KB
)

type command struct {
	numAddresses   int
	virtAddrBounds mem.B

	// segmentA grows in positive direction from the beginning of virtual address
	// space.
	segmentA mem.Segment

	// segmentB grows in negative direction from the end of virtual address space.
	segmentB mem.Segment
}

func newCommand() *command {
	return &command{
		numAddresses:   5,
		virtAddrBounds: maxVirtAddr,
		segmentA: mem.Segment{
			Base:   1 * mem.KB,
			Bounds: 1 * mem.KB,
		},
		segmentB: mem.Segment{
			Base:      10 * mem.KB,
			Bounds:    2 * mem.KB,
			VirtBase:  maxVirtAddr,
			Direction: mem.DirNegative,
		},
	}
}

// Run parses flags and translates random virtual addresses.
func (cmd *command) Run() error {
	if err := cmd.parseFlags(); err != nil {
		return err
	}
	d := report.Data{
		VirtAddrBounds: cmd.virtAddrBounds,
		Segments:       []*mem.Segment{&cmd.segmentA, &cmd.segmentB},
		Translations:   cmd.translations(),
	}
	return report.Generate(os.Stdout, d)
}

func (cmd *command) parseFlags() error {
	fs := flags.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)

	fs.IntVar(&cmd.numAddresses, "n", cmd.numAddresses, "number of addresses")
	vmb := newAddrBoundsFlag(&cmd.virtAddrBounds, 1*mem.KB, maxVirtAddr)
	fs.Var(vmb, "vm-bounds", "virtual address bounds (KB)")
	fs.Var(newSegmentFlag(&cmd.segmentA, maxSegSize), "segA", "first segment")
	fs.Var(newSegmentFlag(&cmd.segmentB, maxSegSize), "segB", "second segment")

	if err := fs.ParseAndValidate(os.Args[1:]); err != nil {
		return err
	}
	if cmd.numAddresses < 1 {
		return fmt.Errorf("invalid flag -n: non-positive number %d", cmd.numAddresses)
	}

	// update [cmd.segmentB] virtual base to possible changes to vm-bounds.
	cmd.segmentB.VirtBase = cmd.virtAddrBounds
	return nil
}

func (cmd *command) translations() iter.Seq[report.Translation] {
	tr := mem.NewTranslator(cmd.segmentA, cmd.segmentB)
	translate := func(a mem.Address) *translation {
		t := &translation{
			virt: a,
		}
		t.phys, t.err = tr.Translate(a)
		return t
	}
	return func(yield func(report.Translation) bool) {
		for a := range cmd.addresses() {
			t := translate(a)
			if !yield(t) {
				break
			}
		}
	}
}

func (cmd *command) addresses() iter.Seq[mem.Address] {
	return func(yield func(mem.Address) bool) {
		for range cmd.numAddresses {
			segm := rand.IntN(2)
			b := rand.IntN(int(cmd.virtAddrBounds))
			virt := mem.MakeAddress(segm, mem.B(b))
			if !yield(virt) {
				break
			}
		}
	}
}

type translation struct {
	virt mem.Address
	phys mem.Address
	err  error
}

// Errors gives access to the translation error if any.
func (tr translation) Error() error {
	return tr.err
}

// Physical returns physical address from the translation. It is equal to zero
// value if translation failed and [translation.Error] returns non-nil error.
func (tr translation) Physical() mem.Address {
	return tr.phys
}

// Virtual returns the translation virtual address.
func (tr translation) Virtual() mem.Address {
	return tr.virt
}
