// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ping_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/book/irex/intent/feature/ping"
	"github.com/skhal/lab/book/irex/pb"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestFulfill(t *testing.T) {
	tests := []struct {
		name   string
		intent *pb.PingIntent
		want   *pb.PingFeature
	}{
		{
			name:   "empty intent",
			intent: &pb.PingIntent{},
			want:   &pb.PingFeature{},
		},
		{
			name:   "non empty intent",
			intent: pb.PingIntent_builder{Message: new("test")}.Build(),
			want: pb.PingFeature_builder{
				Ping: new("test"),
				Pong: new("tset"),
			}.Build(),
		},
		{
			name:   "unicode intent",
			intent: pb.PingIntent_builder{Message: new("テスト")}.Build(),
			want: pb.PingFeature_builder{
				Ping: new("テスト"),
				Pong: new("トステ"),
			}.Build(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ping.Fulfill(tc.intent)

			if err != nil {
				t.Errorf("unexpected error '%v'", err)
			}
			if d := cmp.Diff(tc.want, got, protocmp.Transform()); d != "" {
				t.Errorf("mismatch (-want +got):\n%s", d)
			}
		})
	}
}
