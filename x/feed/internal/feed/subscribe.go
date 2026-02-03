// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package feed

import (
	"errors"
	"fmt"
	"sync"

	"github.com/skhal/lab/x/feed/internal/pb"
)

// Feed is a stream of RSS, Atom, etc. feed items.
type Feed <-chan *Item

// Subscription is a client API to access streamed feed.
type Subscription interface {
	// Feed generates a stream of feed items. It returns an error if it fails to
	// create a stream.
	Feed() (Feed, error)
	// Close stops the subscription and closes the feed.
	Close() error
}

// Subscribe creates a feed subscription.
func Subscribe(f *pb.Feed) Subscription {
	return newSubscription(f)
}

type subscription struct {
	feed *pb.Feed
	done chan struct{}
}

func newSubscription(f *pb.Feed) *subscription {
	return &subscription{
		feed: f,
		done: make(chan struct{}),
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
	return s.streamItems(items), nil
}

// Close stops the subscription and closes the feed.
func (s *subscription) Close() error {
	close(s.done)
	return nil
}

func (s *subscription) streamItems(items []*Item) Feed {
	stream := make(chan *Item)
	go func() {
		defer close(stream)
		for _, item := range items {
			select {
			case stream <- item:
			case <-s.done:
				return
			}
		}
	}()
	return Feed(stream)
}

// String implements fmt.Stringer interface.
func (s *subscription) String() string {
	return s.feed.String()
}

// Merge multiplexes subscriptions into a single subscription.
func Merge(subs []Subscription) Subscription {
	return newMultiplexer(subs)
}

type multiplexer struct {
	subs []Subscription
}

func newMultiplexer(subs []Subscription) *multiplexer {
	return &multiplexer{
		subs: subs,
	}
}

// Feed multiplexes multiple subscription feeds into a single feed.
func (mux *multiplexer) Feed() (Feed, error) {
	stream := make(chan *Item)
	go func() {
		defer close(stream)
		var wg sync.WaitGroup
		defer wg.Wait()
		for _, sub := range mux.subs {
			feed, err := sub.Feed()
			if err != nil {
				continue
			}
			wg.Go(func() {
				for {
					item, ok := <-feed
					if !ok {
						break
					}
					stream <- item
				}
			})
		}
	}()
	return Feed(stream), nil
}

// Close stops multiplexed subscriptions. It returns a joined error from failed
// subscriptions.
func (mux *multiplexer) Close() error {
	var ee []error
	for _, s := range mux.subs {
		if err := s.Close(); err != nil {
			ee = append(ee, err)
		}
	}
	return errors.Join(ee...)
}
