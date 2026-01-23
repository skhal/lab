// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package feed

import (
	"fmt"

	"github.com/skhal/lab/x/feed/internal/pb"
)

type Feed <-chan *Item

// Subscription is a client API to access streamed feed.
type Subscription interface {
	// Feed generates a stream of feed items. It returns an error if it fails to
	// create a stream.
	Feed() (Feed, error)
}

// Subscribe creates a feed subscription.
func Subscribe(f *pb.Feed) Subscription {
	return newSubscription(f)
}

type subscription struct {
	feed *pb.Feed
}

func newSubscription(f *pb.Feed) *subscription {
	return &subscription{
		feed: f,
	}
}

// Feed starts a stream of feed items or returns an error if it fails.
func (s *subscription) Feed() (Feed, error) {
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
