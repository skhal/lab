// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/skhal/lab/book/irex/web"
)

const defaultAddress = ":8080"

type cmdServe struct {
	addr string
}

// Run parses flags and runs the web server.
func (cmd *cmdServe) Run(args []string) error {
	if err := cmd.parseFlags(args); err != nil {
		return err
	}
	return cmd.run()
}

func (cmd *cmdServe) parseFlags(args []string) error {
	name := fmt.Sprintf("%s serve", flag.CommandLine.Name())
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	fs.Usage = func() {
		w := fs.Output()
		fmt.Fprintf(w, "usage: %s [flags]\n", fs.Name())
		fmt.Fprintln(w)
		fmt.Fprintln(w, "flags:")
		fs.PrintDefaults()
	}
	fs.StringVar(&cmd.addr, "http", defaultAddress, "address to bind")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(fs.Args()) != 0 {
		err := fmt.Errorf("unexpected arguments: %s", strings.Join(fs.Args(), " "))
		return newUsageError(fs, err)
	}
	return nil
}

func (cmd *cmdServe) run() error {
	wait, err := cmd.runWebServer()
	if err != nil {
		return err
	}
	return wait()
}

func (cmd *cmdServe) runWebServer() (func() error, error) {
	s := &web.Server{Address: cmd.addr}

	fmt.Printf("web server: start on %s\n", cmd.addr)
	if err := s.Run(); err != nil {
		return nil, err
	}

	var (
		wg  sync.WaitGroup
		err error
	)
	wg.Go(func() {
		sigint := make(chan os.Signal, 1)
		defer close(sigint)

		signal.Notify(sigint, os.Interrupt)

		<-sigint
		signal.Stop(sigint)

		fmt.Println()
		fmt.Println("web server: shutdown")
		err = s.Shutdown(context.Background())
	})

	wait := func() error {
		wg.Wait()
		if err := s.Err(); err != nil {
			return err
		}
		return err
	}

	return wait, nil
}
