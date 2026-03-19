// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"os"
	"path/filepath"

	"github.com/skhal/lab/book/ostep/mem/segment/internal/mem"
	"github.com/skhal/lab/go/flags"
)

type command struct {
	segments []mem.Segment
}

func newCommand() *command {
	return &command{
		segments: []mem.Segment{
			{Base: 5 * mem.KB, Bounds: 1 * mem.KB},
			{Base: 10 * mem.KB, Bounds: 2 * mem.KB},
		},
	}
}

// Run executes the command.
func (cmd *command) Run() error {
	if err := cmd.parseFlags(); err != nil {
		return err
	}
	return errors.New("not implemented")
}

func (cmd *command) parseFlags() error {
	fs := flags.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)
	fs.Var(newSegmentListFlag(&cmd.segments), "segments", "list of segments base:bounds in KB")
	return fs.ParseAndValidate(os.Args[1:])
}
