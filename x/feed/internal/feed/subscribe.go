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
	feed chan *Item
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
	s.feed = make(chan *Item)
	go func() {
		defer close(s.feed)
		s.stream(fetcher)
	}()
	return nil
}

type fetchResult struct {
	items []*Item
	err   error
}

func (s *subscription) stream(f Fetcher) {
	var (
		pending []*Item
		err     error
	)
	next := func() (chan<- *Item, *Item) {
		if len(pending) == 0 {
			return nil, nil
		}
		return s.feed, pending[0]
	}
	var (
		fetchDone chan fetchResult
		nextFetch time.Time
	)
	fetchTime := func() <-chan time.Time {
		var (
			t     <-chan time.Time
			delay time.Duration
		)
		if fetchDone != nil {
			return t
		}
		if nextFetch.After(time.Now()) {
			delay = time.Until(nextFetch)
		}
		if len(pending) < maxPendingSize {
			t = time.After(delay)
		}
		return t
	}
	for {
		send, item := next()
		select {
		case errc := <-s.stop:
			errc <- err
			close(errc)
			return
		case <-fetchTime():
			fetchDone = make(chan fetchResult)
			go func() {
				defer close(fetchDone)
				items, err := f.Fetch()
				fetchDone <- fetchResult{items, err}
			}()
		case res := <-fetchDone:
			if res.err != nil {
				nextFetch = time.Now().Add(fetchBackoffDelay)
				break
			}
			pending = append(pending, res.items...)
			fetchDone = nil
		case send <- item:
			pending = pending[1:]
		}
	}
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
