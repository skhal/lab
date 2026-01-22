// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

import (
	"fmt"
	"sync"
)

type policy interface {
	Next() int
}

// B is a load balancer configured with a round-robin policy on a pool of
// workers.
type B struct {
	requests chan *Request
	pool     *RoundRobinPool
	once     *sync.Once
}

// New constructs [B] with provided number of workers.
func New(numWorkers int) *B {
	return &B{
		requests: make(chan *Request),
		pool:     newRoundRobinPool(numWorkers),
		once:     new(sync.Once),
	}
}

// Balance starts workers, load balances the requests, stops workers and waits
// for scheduled work to complete.
func (b *B) Balance() {
	b.pool.Start()
	b.balance()
	b.pool.Stop()
	b.pool.Wait()
}

func (b *B) balance() {
	for {
		req, ok := <-b.requests
		if !ok {
			break
		}
		b.pool.Dispatch(req)
		b.report()
	}
}

func (b *B) report() {
	stats := b.pool.Stats()
	fmt.Println(stats)
}

// Client creates a new client to send requests to the balancer.
func (b *B) Client() *Client {
	return newClient(b.requests)
}

// Stop closes the incoming requests channel. An attempt to send a request using
// one of the clients after the call to stop the balancer results in a runtime
// panic.
func (b *B) Stop() {
	// TODO: verity all clients are down
	b.once.Do(func() { close(b.requests) })
}
