// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ping fulfills the ping feature.
package ping

import (
	"strings"
	"unicode/utf8"

	"github.com/skhal/lab/book/irex/pb"
)

// Fulfilll reverses the message in the ping intent. The returned ping feature
// includes the requested message (ping) along with the reversed value (pong).
func Fulfill(msg *pb.PingIntent) (*pb.PingFeature, error) {
	s := msg.GetMessage()
	if s == "" {
		return &pb.PingFeature{}, nil
	}
	return pb.PingFeature_builder{
		Ping: &s,
		Pong: new(reverse(s)),
	}.Build(), nil
}

// reverse reverses a utf8 string.
func reverse(s string) string {
	var b strings.Builder
	for len(s) > 0 {
		r, n := utf8.DecodeLastRuneInString(s)
		if r == utf8.RuneError && n == 0 {
			// end-of-string
			break
		}
		if r == utf8.RuneError && n == 1 {
			// invalid encoding
			break
		}
		b.WriteRune(r)
		s = s[:len(s)-n]
	}
	return b.String()
}
