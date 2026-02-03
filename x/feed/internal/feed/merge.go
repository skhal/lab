// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package feed

import (
	"errors"
	"sync"
)

// Merge multiplexes subscriptions into a single subscription.
func Merge(subs []Subscription) Subscription {
	return newMultiplexer(subs)
}

type multiplexer struct {
	subs []Subscription
	feed Feed
	once sync.Once
}

func newMultiplexer(subs []Subscription) *multiplexer {
	return &multiplexer{
		subs: subs,
	}
}

// Feed multiplexes multiple subscription feeds into a single feed.
func (mux *multiplexer) Feed() (Feed, error) {
	mux.once.Do(func() { mux.merge() })
	return mux.feed, nil
}

func (mux *multiplexer) merge() {
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
	mux.feed = Feed(stream)
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
