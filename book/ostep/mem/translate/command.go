// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"iter"
	"os"
	"path/filepath"
)

const (
	defaultNum              = 5
	defaultBase             = 16 * 1024 // 16KB
	defaultBounds           = 3 * 1024  // 3KB
	defaultVirtAddressSpace = 4 * 1024  // 4KB
)

// ErrFlag means the flag value is invalid.
var ErrFlag = errors.New("invalid flag")

type command struct {
	num              int
	base             int
	bounds           int
	virtAddressSpace int
}

// Run executes the command
func (c *command) Run() error {
	if err := c.parseFlags(); err != nil {
		return err
	}
	return c.run()
}

func (c *command) parseFlags() error {
	fs := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)
	fs.IntVar(&c.num, "n", defaultNum, "number of addresses to generate")
	fs.IntVar(&c.base, "base", defaultBase, "base of the physical address")
	fs.IntVar(&c.bounds, "bounds", defaultBounds, "size of process address space")
	fs.IntVar(&c.virtAddressSpace, "virt-addr-space", defaultVirtAddressSpace, "virtual address space size")
	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}
	if c.num < 1 {
		return fmt.Errorf("%w: num %d: non-positive value", ErrFlag, c.num)
	}
	if c.base < 0 {
		return fmt.Errorf("%w: base %d: negative value", ErrFlag, c.base)
	}
	if c.bounds < 0 {
		return fmt.Errorf("%w: bounds %d: negative value", ErrFlag, c.bounds)
	}
	if c.virtAddressSpace < 1 {
		return fmt.Errorf("%w: virt-addr-space %d: non-positive value", ErrFlag, c.virtAddressSpace)
	}
	return nil
}

func (c *command) run() error {
	t := newTranslator(Address(c.base), Address(c.bounds))
	for virt := range c.genAddress() {
		phys, err := t.Translate(virt)
		if err != nil {
			fmt.Printf("virt: %s %v\n", virt, err)
			continue
		}
		fmt.Printf("virt: %s phys: %s\n", virt, phys)
	}
	return nil
}

func (c *command) genAddress() iter.Seq[Address] {
	return func(yield func(Address) bool) {
		g := newGenerator(c.num, Address(c.virtAddressSpace))
		for g.Next() && yield(g.Address()) {
			continue
		}
	}
}
