// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"math/rand"
	"sync"
	"time"

	"github.com/skhal/lab/x/balancer/internal/balancer"
)

var (
	numWorkers  = flag.Int("w", 1, "number of workers")
	numRequests = flag.Int("r", 10, "number of requests")
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

func main() {
	flag.Parse()
	b := balancer.New(*numWorkers)
	run(b, *numRequests)
}

func run(b *balancer.B, numRequests int) {
	var wg sync.WaitGroup
	wg.Go(b.Balance)
	sendRequests(b, numRequests)
	b.Stop()
	wg.Wait()
}

func sendRequests(b *balancer.B, num int) {
	c := b.Client()
	for n := 0; n < num; n++ {
		req := balancer.NewRequest(task)
		c.Send(req)
	}
}

func task() {
	n := 1 + rand.Int63n(10)
	time.Sleep(time.Duration(n) * time.Millisecond)
}
