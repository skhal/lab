// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

// Client implements a mechanism to send requests to the balancer.
type Client struct {
	requests chan<- *Request
}

func newClient(rr chan<- *Request) *Client {
	return &Client{
		requests: rr,
	}
}

// Send submits a request to the load balancer.
func (c *Client) Send(req *Request) {
	c.requests <- req
}
