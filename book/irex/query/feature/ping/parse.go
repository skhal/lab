// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ping parses query as a ping intent.
package ping

import "github.com/skhal/lab/book/irex/pb"

// Parse parses ping-command parameters.
func Parse(params string) (*pb.PingIntent, error) {
	msg := pb.PingIntent_builder{
		Message: &params,
	}.Build()
	return msg, nil
}
