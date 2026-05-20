// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package options_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/dp/proto/options"
	"github.com/skhal/lab/dp/proto/options/pb"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestClean(t *testing.T) {
	tests := []struct {
		name string
		msg  *pb.Store
		want *pb.Store
	}{
		{
			name: "no secret",
			msg:  pb.Store_builder{Public: new("show")}.Build(),
			want: pb.Store_builder{Public: new("show")}.Build(),
		},
		{
			name: "with secret",
			msg:  pb.Store_builder{Public: new("show"), Private: new("hide")}.Build(),
			want: pb.Store_builder{Public: new("show")}.Build(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			options.Clean(tc.msg)

			if d := cmp.Diff(tc.want, tc.msg, protocmp.Transform()); d != "" {
				t.Errorf("mismatch (-want +got):\n%s", d)
			}
		})
	}
}
