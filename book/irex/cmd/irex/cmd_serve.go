// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/skhal/lab/book/irex/market"
	"github.com/skhal/lab/book/irex/pb"
	"github.com/skhal/lab/book/irex/web"
	"google.golang.org/protobuf/proto"
)

const defaultAddress = ":8080"

type cmdServe struct {
	addr           string
	marketDataFile string // market data
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
	fs.StringVar(&cmd.marketDataFile, "f", "", "market data file")
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
	waitMarketServer, err := cmd.runMarketServer()
	if err != nil {
		return err
	}
	waitWebServer, err := cmd.runWebServer()
	if err != nil {
		return err
	}

	var (
		wg sync.WaitGroup
		ee = make([]error, 2)
	)
	const (
		idxErrMarketServer = iota
		idxErrWebServer
	)
	wg.Go(func() {
		if err := waitMarketServer(); err != nil {
			ee[idxErrMarketServer] = err
		}
	})
	wg.Go(func() {
		if err := waitWebServer(); err != nil {
			ee[idxErrWebServer] = err
		}
	})
	wg.Wait()

	return errors.Join(ee...)
}

type waitFunc func() error

func (cmd *cmdServe) runMarketServer() (waitFunc, error) {
	data, err := cmd.loadMarketData()
	if err != nil {
		return nil, err
	}
	s := market.NewServer(data)

	if err := s.Serve(); err != nil {
		return nil, err
	}
	fmt.Println("market server: started")

	var wg sync.WaitGroup
	wg.Go(func() {
		sigint := make(chan os.Signal, 1)
		defer close(sigint)

		signal.Notify(sigint, os.Interrupt)

		<-sigint
		signal.Stop(sigint)

		fmt.Println("market server: shutdown")
		s.Shutdown()
	})

	wait := func() error {
		wg.Wait()
		return s.Err()
	}

	return wait, nil
}

func (cmd *cmdServe) loadMarketData() (*pb.Market, error) {
	b, err := os.ReadFile(cmd.marketDataFile)
	if err != nil {
		return nil, err
	}
	data := new(pb.Market)
	if err := proto.Unmarshal(b, data); err != nil {
		return nil, fmt.Errorf("%s: %s", cmd.marketDataFile, err)
	}
	return data, nil
}

func (cmd *cmdServe) runWebServer() (waitFunc, error) {
	s := new(web.Server)

	if err := s.ListenAndServe(cmd.addr); err != nil {
		return nil, err
	}
	fmt.Printf("web server: started at %s\n", s.Addr())

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
