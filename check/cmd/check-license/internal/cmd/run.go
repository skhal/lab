// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skhal/lab/check/cmd/check-license/internal/license"
)

// ErrFlags indicates error with input flags.
var ErrFlags = errors.New("invalid flags")

// Run executes the check-license command.
func Run() error {
	cmd := new(command)
	return cmd.Run()
}

type command struct {
	fix    bool   // fix license if missing
	holder string // attribute license to the holder
	files  []string
}

// Run parses flags and executes the check.
func (cmd *command) Run() error {
	if err := cmd.parseFlags(); err != nil {
		return err
	}
	return license.Run(cmd.files, cmd.fix, cmd.holder)
}

func (cmd *command) parseFlags() error {
	fs := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)
	fs.BoolVar(&cmd.fix, "fix", cmd.fix, "insert license if missing")
	fs.StringVar(&cmd.holder, "holder", cmd.holder, "license holder name")
	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}
	cmd.files = fs.Args()
	if cmd.fix && cmd.holder == "" {
		return fmt.Errorf("%w: missing -holder with -fix", ErrFlags)
	}
	return nil
}
