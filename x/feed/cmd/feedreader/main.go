// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Feedreader implements an RSS, Atom, etc. feed reader using streaming feeds.
package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/skhal/lab/x/feed/internal/feed"
	"github.com/skhal/lab/x/feed/internal/pb"
	"google.golang.org/protobuf/encoding/prototext"
)

const defaultTimeout = 100 * time.Millisecond

type config struct {
	file    string
	timeout time.Duration
}

// Validate checks whether the configuration has meaningful values, i.e., the
// file is not empty, etc. It returns a description of failed configuration
// parameter.
func (cfg *config) Validate() error {
	if len(cfg.file) == 0 {
		return fmt.Errorf("file is not set")
	}
	return nil
}

func main() {
	cfg := mustParseFlags()
	if err := run(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func mustParseFlags() *config {
	fs := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() {
		w := fs.Output()
		fmt.Fprintf(w, "usage: %s -f file\n", fs.Name())
		fmt.Fprintln(w, "flags:")
		fs.PrintDefaults()
	}
	cfg := new(config)
	fs.StringVar(&cfg.file, "f", "", "feeds file")
	fs.DurationVar(&cfg.timeout, "timeout", defaultTimeout, "stop timeout")
	err := fs.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fs.Usage()
		os.Exit(1)
	}
	if err := cfg.Validate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		fs.Usage()
		os.Exit(1)
	}
	return cfg
}

func run(cfg *config) error {
	ctx, cancel := context.WithTimeoutCause(context.Background(), cfg.timeout, fmt.Errorf("timeout %v", cfg.timeout))
	defer cancel()
	feeds, err := readConfig(cfg.file)
	if err != nil {
		return err
	}
	if err := readFeeds(ctx, feeds); err != nil {
		return err
	}
	return ctx.Err()
}

func readConfig(name string) (*pb.FeedSet, error) {
	b, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	fset := new(pb.FeedSet)
	if err := prototext.Unmarshal(b, fset); err != nil {
		return nil, err
	}
	return fset, nil
}

func readFeeds(ctx context.Context, feeds *pb.FeedSet) error {
	f, err := subscribe(feeds).Feed()
	if err != nil {
		return err
	}
	readFeed(ctx, f)
	return nil
}

func readFeed(ctx context.Context, f feed.Feed) {
	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Go(func() {
		for stop := false; !stop; {
			select {
			case <-ctx.Done():
				stop = true
			case item, ok := <-f:
				if !ok {
					stop = true
					break
				}
				// emulate a delay
				time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
				fmt.Printf("%s\n", (*printableItem)(item))
			}
		}
	})
}

func subscribe(feeds *pb.FeedSet) feed.Subscription {
	subs := make([]feed.Subscription, 0, len(feeds.GetFeeds()))
	for _, f := range feeds.GetFeeds() {
		subs = append(subs, feed.Subscribe(f))
	}
	return feed.Merge(subs)
}

type printableItem feed.Item

// String implements fmt.Stringer interface
func (i *printableItem) String() string {
	var t *time.Time
	switch {
	case i.Published != nil:
		t = i.Published
	case i.Updated != nil:
		t = i.Updated
	}
	tstr := "N/A"
	if t != nil {
		tstr = t.Format(time.RFC822)
	}
	return fmt.Sprintf("[%s] %s", tstr, i.Title)
}
