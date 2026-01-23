// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package feed

import (
	"fmt"

	"github.com/skhal/lab/x/feed/internal/pb"
)

// Subscribe subscribes to a feed. It generates a stream of feed items, or
// returns an error if generation fails.
func Subscribe(f *pb.Feed) (Feed, error) {
	s := newSubscription(f)
	return s.Subscribe()
}

type subscription struct {
	feed *pb.Feed
}

func newSubscription(f *pb.Feed) *subscription {
	return &subscription{
		feed: f,
	}
}

// Subscribe starts a stream of feed items or returns an error if it fails.
func (s *subscription) Subscribe() (Feed, error) {
	if !s.feed.GetSource().HasSource() {
		return nil, fmt.Errorf("subscribe %s: missing source", s.feed)
	}
	items, err := Fetch(s.feed.GetSource())
	if err != nil {
		return nil, err
	}
	stream := make(chan *Item)
	go func() {
		defer close(stream)
		for _, item := range items {
			stream <- item
		}
	}()
	return Feed(stream), nil
}

func (s *subscription) String() string {
	return s.feed.String()
}
