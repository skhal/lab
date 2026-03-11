// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/skhal/lab/check/cmd/check-license/internal/check"
	"github.com/skhal/lab/check/cmd/check-license/internal/fix"
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

var runTimeout = 100 * time.Millisecond

// Run parses flags and executes the check.
func (cmd *command) Run() error {
	if err := cmd.parseFlags(); err != nil {
		return err
	}
	ch := make(chan error)
	go func() {
		defer close(ch)
		ctx, cancel := context.WithTimeout(context.Background(), runTimeout)
		defer cancel()
		var wg sync.WaitGroup
		for _, f := range cmd.files {
			wg.Go(func() {
				select {
				case ch <- cmd.run(f):
				case <-ctx.Done():
				}
			})
		}
		wg.Wait()
	}()
	var ee []error
	for err := range ch {
		ee = append(ee, err)
	}
	return errors.Join(ee...)
}

func (cmd *command) parseFlags() error {
	fs := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)
	fs.BoolVar(&cmd.fix, "fix", cmd.fix, "insert license if missing, requires -holder if set")
	fs.StringVar(&cmd.holder, "holder", cmd.holder, "license holder name")
	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}
	cmd.files = fs.Args()
	if cmd.fix && cmd.holder == "" {
		return fmt.Errorf("%w: -fix requires -holder", ErrFlags)
	}
	return nil
}

func (cmd *command) run(file string) error {
	info, err := os.Stat(file)
	if err != nil {
		return err
	}
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	if isBinary(b) {
		// skip binary files
		return nil
	}
	switch cmd.fix {
	case true:
		if err := check.Run(b); err == nil {
			break
		}
		b, err = fix.Run(file, b, cmd.holder)
		if err != nil {
			return fmt.Errorf("%s: %s", file, err)
		}
		err = os.WriteFile(file, b, info.Mode())
		if err != nil {
			return fmt.Errorf("%s: %s", file, err)
		}
		// report file change via error
		return fmt.Errorf("%s: fixed", file)
	case false:
		if err := check.Run(b); err != nil {
			return fmt.Errorf("%s: %s", file, err)
		}
	}
	return nil
}

const nullByte = 0

func isBinary(b []byte) bool {
	return bytes.IndexByte(b, nullByte) != -1
}
