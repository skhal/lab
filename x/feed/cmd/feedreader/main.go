// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
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

func main() {
	file := mustParseFlags()
	if err := run(file); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func mustParseFlags() string {
	fs := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() {
		w := fs.Output()
		fmt.Fprintf(w, "usage: %s -f file\n", fs.Name())
		fmt.Fprintln(w, "flags:")
		fs.PrintDefaults()
	}
	file := fs.String("f", "", "feeds file")
	err := fs.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fs.Usage()
		os.Exit(1)
	}
	if *file == "" {
		fmt.Fprintln(os.Stderr, "missing feeds file")
		fs.Usage()
		os.Exit(1)
	}
	return *file
}

func run(name string) error {
	feeds, err := load(name)
	if err != nil {
		return err
	}
	return read(feeds)
}

func load(name string) (*pb.FeedSet, error) {
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

func read(feeds *pb.FeedSet) error {
	var wg sync.WaitGroup
	defer wg.Wait()
	for _, f := range feeds.GetFeeds() {
		stream, err := feed.Subscribe(f)
		if err != nil {
			return err
		}
		wg.Go(func() {
			for item := range stream {
				// emulate delay
				time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
				fmt.Println(f.GetName(), ": ", printableItem(item))
			}
		})
	}
	return nil
}

type printableItem feed.Item

func (i printableItem) String() string {
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
