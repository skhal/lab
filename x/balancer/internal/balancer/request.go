// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

import "sync"

// Request holds a function to run and channel to communicate the completion.
type Request struct {
	fn   func()
	done chan struct{}
	once *sync.Once
}

// NewRequest creates a new request to run a function fn.
func NewRequest(fn func()) *Request {
	return &Request{
		fn:   fn,
		done: make(chan struct{}),
		once: new(sync.Once),
	}
}

// Done marks request completed, i.e., the requested function completed running.
func (req *Request) Done() {
	req.once.Do(func() { close(req.done) })
}

// Run executes the requested function.
func (req *Request) Run() {
	req.fn()
}
