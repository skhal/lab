// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package feed

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

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
			items []*Item
			err   error
		)
		sleep := func() {
			time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		}
		for {
			items, err = f.Fetch()
			if err != nil {
				sleep()
				continue
			}
			for _, item := range items {
				select {
				case stream <- item:
				case errc := <-s.stop:
					errc <- err
					close(errc)
					return
				}
			}
			sleep()
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
