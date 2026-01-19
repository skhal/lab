// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Balancer demonstrates a round-robin work load balancer. It runs -w W workers,
// -p P producers, each generating at most -n N requests.
//
// The producers send requests to run a function, and blocks until the function
// runs by waiting for messages on "done" channel.
//
// The balancer round-robins the requests to the workers.
//
// The worker runs the requested function and communicates to the client about
// job completion by closing the "done" channel.
//
// Each component - producer, balancer, worker - have built-in random delay
// to emulate network and work latency.
package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sync"
	"time"
)

var (
	// keep-sorted start
	hFlag = flag.Bool("h", false, "show help and exit")
	nFlag = flag.Int("n", 1, "number of requests per producer")
	pFlag = flag.Int("p", 1, "number of producers")
	wFlag = flag.Int("w", 1, "number of workers")
	// keep-sorted end
)

func main() {
	flag.Parse()
	if *hFlag {
		flag.Usage()
		os.Exit(0)
	}
	b := newBalancer(*wFlag)
	run(b)
	b.Report(os.Stdout)
}

func run(b *Balancer) {
	bwait := b.Run()
	defer bwait()
	pwait := runProducers(*pFlag, *nFlag, b)
	pwait()
	b.Stop()
}

// SleepMillisecond sleeps current goroutine for up to msmax milliseconds.
func SleepMillisecond(msmax int) {
	n := 1 + rand.Intn(int(msmax)) // +1 to sleep for at least 1ms
	time.Sleep(time.Duration(n) * time.Millisecond)
}

// WaitFunc is a barrier, i.e., a synchronization mechanism to wait for a wait
// group to complete. It avoid explicit use of channels or pass around wait
// groups.
type WaitFunc func()

// runProducers starts nproducers producers with up to maxMessages to be sent
// using a per-producer balancer client.
func runProducers(nproducers, maxMessages int, b *Balancer) WaitFunc {
	var wg sync.WaitGroup
	for i := 1; i <= nproducers; i++ {
		p := newProducer(i, b.Client())
		wg.Go(func() { p.Produce(maxMessages) })
	}
	return func() {
		wg.Wait()
	}
}

// RunCompleter can run and mark complete a task.
type RunCompleter interface {
	Run()
	Complete()
}

// Sender can send a RunCompleter.
type Sender interface {
	Send(RunCompleter)
}

// Producer generate and sends a work request using a client.
type Producer struct {
	id     int
	sender Sender
}

func newProducer(id int, s Sender) *Producer {
	return &Producer{
		id:     id,
		sender: s,
	}
}

// Produce sends up to maxMessages requests to run work.
func (p *Producer) Produce(maxMessages int) {
	for n := 1 + rand.Intn(maxMessages); n > 0; n -= 1 {
		SleepMillisecond(10)
		p.send(func() {
			fmt.Printf("P.%d: work\n", p.id)
		})
	}
}

func (p *Producer) send(work func()) {
	req := newRequest(work)
	p.sender.Send(req)
	req.Wait()
}

// Request holds a function to run and a channel to report when the task
// completes.
type Request struct {
	fn   func()        // an RPC
	done chan struct{} // signal when RPC is done
	once *sync.Once
}

func newRequest(fn func()) *Request {
	return &Request{
		fn:   fn,
		done: make(chan struct{}),
		once: new(sync.Once),
	}
}

// Wait blocks until the request completes.
func (r *Request) Wait() {
	<-r.done
}

// Complete marks the request complete and unblocks waiters. It is safe to call
// Complete multiple times - onely first one has effect, all others are no-op.
func (r *Request) Complete() {
	r.once.Do(func() { close(r.done) })
}

// Run executes the requested work.
func (r *Request) Run() {
	r.fn()
}

// Worker listens for RunCompleter on the input channel and executes the work
// once a new one is posted. It marks the work complete at the end and blocks
// until a new work request arrives.
type Worker struct {
	id   int
	run  int
	ch   chan RunCompleter
	once *sync.Once
}

func newWorker(id int) *Worker {
	return &Worker{
		id:   id,
		ch:   make(chan RunCompleter),
		once: new(sync.Once),
	}
}

func (w *Worker) String() string {
	return fmt.Sprintf("worker %d run %d", w.id, w.run)
}

// Run listens on the input channel for work requests. It stops when the input
// channel closes.
func (w *Worker) Run() {
	for {
		rc, ok := <-w.ch
		if !ok {
			break
		}
		w.handle(rc)
	}
}

func (w *Worker) handle(rc RunCompleter) {
	SleepMillisecond(10)
	rc.Run()
	rc.Complete()
	w.run += 1
}

// Stop closes worker's input channel.
func (w *Worker) Stop() {
	w.once.Do(func() { close(w.ch) })
}

// Send sends a request to the worker.
func (w *Worker) Send(rc RunCompleter) {
	w.ch <- rc
}

// Stopper stops processing.
type Stopper interface {
	Stop()
}

// SendStopper can send requests and stop processing.
type SendStopper interface {
	Sender
	Stopper
}

// Policy returns index of the next worker to dispatch the request to.
type Policy interface {
	Next() int
}

// RoundRobinPolicy implements a round robin dispatch through N workers.
type RoundRobinPolicy struct {
	idx  int
	size int
}

func newRoundRobinPolicy(size int) *RoundRobinPolicy {
	return &RoundRobinPolicy{
		size: size,
	}
}

// Next picks the next worker out of a pool f N workers to send the request
// to using round robin logic.
func (p *RoundRobinPolicy) Next() int {
	i := p.idx
	p.idx = (p.idx + 1) % p.size
	return i
}

// Balancer dispatches send requests [RunCompleter] to workers [SendStopper].
// Request generators should use [Client] to send requests to the balancer.
type Balancer struct {
	workers []SendStopper
	policy  Policy

	rr   chan RunCompleter // incoming requests
	once *sync.Once

	clients []*Client // keep track of client for reporting

	dispatched int // number of dispatched requests
}

func newBalancer(size int) *Balancer {
	return &Balancer{
		workers: make([]SendStopper, 0, size),
		policy:  newRoundRobinPolicy(size),
		rr:      make(chan RunCompleter),
		once:    new(sync.Once),
	}
}

func (b *Balancer) String() string {
	return fmt.Sprintf("dispatched %d", b.dispatched)
}

func (b *Balancer) Report(w io.Writer) {
	fmt.Fprintln(w, "Dispatcher:")
	fmt.Fprintln(w, " ", b)
	fmt.Fprintln(w, "Clients:")
	for _, c := range b.clients {
		fmt.Fprintln(w, " ", c)
	}
	fmt.Fprintln(w, "Workers:")
	for _, wkr := range b.workers {
		fmt.Fprintln(w, " ", wkr)
	}
}

// Client is a balancer client to send requests.
type Client struct {
	id   int
	sent int
	send func(rc RunCompleter)
}

// Send submits a request to run a task to the balancer.
func (c *Client) Send(rc RunCompleter) {
	c.send(rc)
	c.sent += 1
}

func (c *Client) String() string {
	return fmt.Sprintf("client %d sent %d", c.id, c.sent)
}

// Client generates a new client.
func (b *Balancer) Client() *Client {
	c := &Client{
		id: len(b.clients) + 1, // +1 to account for the new client
		send: func(rc RunCompleter) {
			b.rr <- rc
		},
	}
	b.clients = append(b.clients, c)
	return c
}

// Stop stops balancer from listening for incoming requests and complete any
// requested work.
func (b *Balancer) Stop() {
	b.once.Do(func() { close(b.rr) })
}

// Run runs a balancer in a goroutine. It returns a wait function to block until
// the balancer finishes.
func (b *Balancer) Run() WaitFunc {
	wg := new(sync.WaitGroup)
	wg.Go(func() {
		wait := b.startWorkers()
		defer wait()
		b.balance()
		b.stopWorkers()
	})
	return func() {
		wg.Wait()
	}
}

func (b *Balancer) startWorkers() WaitFunc {
	wg := new(sync.WaitGroup)
	for len(b.workers) < cap(b.workers) {
		w := newWorker(len(b.workers) + 1)
		wg.Go(w.Run)
		b.workers = append(b.workers, w)
	}
	return func() {
		wg.Wait()
	}
}

func (b *Balancer) balance() {
	for {
		SleepMillisecond(10)
		rc, ok := <-b.rr
		if !ok {
			break
		}
		b.dispatch(rc)
	}
}

func (b *Balancer) stopWorkers() {
	for _, w := range b.workers {
		w.Stop()
	}
}

func (b *Balancer) dispatch(rc RunCompleter) {
	w := b.workers[b.policy.Next()]
	w.Send(rc)
	b.dispatched += 1
}
