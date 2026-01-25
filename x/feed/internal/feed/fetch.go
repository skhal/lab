// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package feed

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mmcdole/gofeed"
	"github.com/skhal/lab/x/feed/internal/pb"
)

// Fetch fetches a feed and generates a stream of items. It returns an error if
// the feed source is not supported or there is an error in accessing the feed.
func Fetch(s *pb.Source) ([]*Item, error) {
	f, err := newFetcher(s)
	if err != nil {
		return nil, err
	}
	return f.Fetch()
}

type fetcher interface {
	Fetch() ([]*Item, error)
}

func newFetcher(s *pb.Source) (fetcher, error) {
	if s.HasFile() {
		return newFileFetcher(s.GetFile()), nil
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
func (ftch *fileFetcher) Fetch() ([]*Item, error) {
	n, err := expand(ftch.file)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(n)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	p := gofeed.NewParser()
	feed, err := p.Parse(file)
	if err != nil {
		return nil, err
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

func transform(s []*gofeed.Item) []*Item {
	ii := make([]*Item, 0, len(s))
	for _, item := range s {
		ii = append(ii, &Item{
			Title:     item.Title,
			Updated:   item.UpdatedParsed,
			Published: item.PublishedParsed,
		})
	}
	return ii
}
