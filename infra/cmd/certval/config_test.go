// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/infra/cmd/certval/pb"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    *pb.Config
		wantErr error
	}{
		{
			name:    "missing file",
			path:    "/non-existent",
			wantErr: os.ErrNotExist,
		},
		{
			name:    "invalid config",
			path:    "testdata/invalid_config.txtpb",
			wantErr: ErrConfigUnmarshal,
		},
		{
			name: "valid config",
			path: "testdata/valid_config.txtpb",
			want: pb.Config_builder{
				Slack: pb.Slack_builder{
					Webhook: new("test hook"),
				}.Build(),
				Certificate: []*pb.Certificate{
					pb.Certificate_builder{
						Name: new("test certificate"),
					}.Build(),
				},
			}.Build(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseConfig(tc.path)

			if d := cmp.Diff(tc.want, got, protocmp.Transform()); d != "" {
				t.Errorf("mismatch (-want, +got):\n%s", d)
			}
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("unexpected error '%v'; want '%v'", err, tc.wantErr)
			}
		})
	}
}

func TestConfigError_Unwrap(t *testing.T) {
	want := errors.New("test error")

	got := errors.Unwrap(newConfigError("test/path", want))

	if got != want {
		t.Errorf("Unwrap got %s; want %s", got, want)
	}
}
