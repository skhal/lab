// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compdb_test

import (
	"errors"
	"testing"

	"github.com/skhal/lab/x/compdb/internal/bazel"
	"github.com/skhal/lab/x/compdb/internal/compdb"
)

func TestGetSource(t *testing.T) {
	tests := []struct {
		name    string
		action  bazel.Action
		want    string
		wantErr error
	}{
		{
			name:    "empty",
			wantErr: compdb.ErrSource,
		},
		{
			name: "no source flag",
			action: bazel.Action{
				Arguments: []string{"a", "b", "c"},
			},
			wantErr: compdb.ErrSource,
		},
		{
			name: "missing source",
			action: bazel.Action{
				Arguments: []string{"a", "b", "-c"},
			},
			wantErr: compdb.ErrSource,
		},
		{
			name: "has source",
			action: bazel.Action{
				Arguments: []string{"a", "b", "-c", "test"},
			},
			want: "test",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			src, err := compdb.GetSource(&tc.action)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("compdb.GetSource() = _, %v; want %v", err, tc.wantErr)
				t.Log(tc.action)
			}
			if src != tc.want {
				t.Errorf("GetSource() = _, %s; want %s", src, tc.want)
				t.Log(tc.action)
			}
		})
	}
}
