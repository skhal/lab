// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package feed

import (
	"fmt"
	"sync"
	"time"

	"github.com/skhal/lab/x/feed/internal/pb"
)

const fetchBackoffDelay = 5 * time.Millisecond
const maxPendingSize = 10

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
	cfg  *pb.Feed
	feed Feed
	stop chan chan error
	once sync.Once
}

func newSubscription(f *pb.Feed) *subscription {
	return &subscription{
		cfg:  f,
		stop: make(chan chan error),
	}
}

// Feed starts a stream of feed items or returns an error if it fails.
func (s *subscription) Feed() (Feed, error) {
	var err error
	s.once.Do(func() { err = s.run() })
	return s.feed, err
}

func (s *subscription) run() error {
	if s.feed != nil {
		return fmt.Errorf("already running")
	}
	if !s.cfg.GetSource().HasSource() {
		return fmt.Errorf("subscribe %s: missing source", s.cfg)
	}
	fetcher, err := Fetch(s.cfg.GetSource())
	if err != nil {
		return err
	}
	s.feed = s.streamFeed(fetcher)
	return nil
}

func (s *subscription) streamFeed(f Fetcher) Feed {
	stream := make(chan *Item)
	go func() {
		defer close(stream)
		var (
			pending []*Item
			err     error
		)
		var nextFetch time.Time
		for {
			var (
				send chan<- *Item
				item *Item
			)
			if len(pending) > 0 {
				send = stream
				item = pending[0]
			}
			var (
				fetchTime  <-chan time.Time
				fetchDelay time.Duration
			)
			if nextFetch.After(time.Now()) {
				fetchDelay = time.Until(nextFetch)
			}
			if len(pending) < maxPendingSize {
				fetchTime = time.After(fetchDelay)
			}
			select {
			case errc := <-s.stop:
				errc <- err
				close(errc)
				return
			case <-fetchTime:
				var items []*Item
				items, err = f.Fetch()
				if err != nil {
					nextFetch = time.Now().Add(fetchBackoffDelay)
					break
				}
				pending = append(pending, items...)
			case send <- item:
				pending = pending[1:]
			}
		}
	}()
	return stream
}

// Close stops the subscription and closes the feed.
func (s *subscription) Close() error {
	err := make(chan error)
	s.stop <- err
	return <-err
}

// String implements fmt.Stringer interface.
func (s *subscription) String() string {
	return s.cfg.String()
}
