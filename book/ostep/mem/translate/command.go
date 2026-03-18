// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultBase   = 16 * 1024 // 16KB
	defaultBounds = 4 * 1024  // 4KB
)

var (
	// ErrBase means base address has invalid value, e.g. negative.
	ErrBase = errors.New("invalid base")

	// ErrBounds means address bounds have invalid value, e.g. negative.
	ErrBounds = errors.New("invalid bounds")
)

type command struct {
	base   int
	bounds int
}

// Run executes the command
func (c *command) Run() error {
	if err := c.parseFlags(); err != nil {
		return err
	}
	return errors.New("not implemented")
}

func (c *command) parseFlags() error {
	fs := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)
	fs.IntVar(&c.base, "base", defaultBase, "base of the physical address")
	fs.IntVar(&c.bounds, "bounds", defaultBounds, "size of process address space")
	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}
	if c.base < 0 {
		return fmt.Errorf("%w: negative value %d", ErrBase, c.base)
	}
	if c.bounds < 0 {
		return fmt.Errorf("%w: negative value %d", ErrBounds, c.bounds)
	}
	return nil
}
