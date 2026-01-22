// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

import "sync"

// Worker executes the requested function. It listens for the requests on the
// incoming channel and runs a function when available.
type Worker struct {
	requests chan *Request
	once     *sync.Once
	done     chan struct{}
}

func newWorker() *Worker {
	return &Worker{
		requests: make(chan *Request),
		once:     new(sync.Once),
		done:     make(chan struct{}),
	}
}

// Wait waits until the worker finishes processing data.
func (w *Worker) Wait() bool {
	_, ok := <-w.done
	return ok
}

// Run makes worker listen for incoming requests and execute the requested
// functions.
func (w *Worker) Run() {
	for {
		req, ok := <-w.requests
		if !ok {
			break
		}
		w.handle(req)
		w.reportDone()
	}
	w.closeDone()
}

func (w *Worker) handle(req *Request) {
	w.run(req)
	w.complete(req)
}

func (w *Worker) run(req *Request) {
	req.Run()
}

func (w *Worker) complete(req *Request) {
	req.Done()
}

func (w *Worker) reportDone() {
	select {
	default: // do not block in case nothing is listening for the feedback
	case w.done <- struct{}{}:
	}
}

func (w *Worker) closeDone() {
	close(w.done)
}

func (w *Worker) Send(req *Request) {
	w.requests <- req
}

func (w *Worker) Stop() {
	w.once.Do(func() { close(w.requests) })
}
