// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package feed

import (
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/mmcdole/gofeed"
	"github.com/skhal/lab/x/feed/internal/pb"
)

const blockSize = 5 // bock size to fetch items from files.

// Fetch fetches a feed and generates a stream of items. It returns an error if
// the feed source is not supported or there is an error in accessing the feed.
func Fetch(s *pb.Source) (Fetcher, error) {
	f, err := newFetcher(s)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Fetcher is responsible for getting items from the feed.
type Fetcher interface {
	// Fetch retrieves items from the feed.
	Fetch() ([]*Item, error)
}

func newFetcher(s *pb.Source) (Fetcher, error) {
	if s.HasFile() {
		f := newFileFetcher(s.GetFile())
		return newBlockFileFetcher(f, blockSize), nil
	}
	return nil, fmt.Errorf("subscribe %s: non-file sources are not supported", s)
}

type fileFetcher struct {
	file string
}

func newFileFetcher(name string) *fileFetcher {
	return &fileFetcher{file: name}
}

// Fetch retrieves feed items from the file.
func (f *fileFetcher) Fetch() ([]*Item, error) {
	name, err := expand(f.file)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return f.parse(file)
}

func (f *fileFetcher) parse(file *os.File) ([]*Item, error) {
	parser := gofeed.NewParser()
	feed, err := parser.Parse(file)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", f.file, err)
	}
	return transform(feed.Items), nil
}

func expand(name string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	expanded := strings.Replace(name, "~/", home+"/", 1)
	return filepath.Clean(expanded), nil
}

func transform(items []*gofeed.Item) []*Item {
	ii := make([]*Item, 0, len(items))
	for _, item := range items {
		ii = append(ii, &Item{
			Title:     item.Title,
			Updated:   item.UpdatedParsed,
			Published: item.PublishedParsed,
		})
	}
	return ii
}

type blockFileFetcher struct {
	fetcher Fetcher

	once  sync.Once
	items []*Item
	err   error

	size      int
	nextBlock func() (Block, bool)
}

func newBlockFileFetcher(f Fetcher, blockSize int) *blockFileFetcher {
	return &blockFileFetcher{
		fetcher: f,
		size:    blockSize,
	}
}

// Fetch generates a block of items from the list of items retrieved from the
// wrapped fetcher.
func (f *blockFileFetcher) Fetch() ([]*Item, error) {
	f.once.Do(func() {
		f.items, f.err = f.fetcher.Fetch()
		f.nextBlock, _ = iter.Pull(EqualSizeBlocks(f.size, len(f.items)))
	})
	if f.err != nil {
		return nil, f.err
	}
	return f.fetchItems()
}

func (f *blockFileFetcher) fetchItems() ([]*Item, error) {
	c, ok := f.nextBlock()
	if !ok {
		return nil, nil
	}
	return f.items[c.Low:c.High], nil
}
