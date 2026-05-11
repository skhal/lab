// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compdb_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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

func TestMakeCommands(t *testing.T) {
	tests := []struct {
		name         string
		actions      []*bazel.Action
		wantCommands []*compdb.Command
		wantErr      error
	}{
		{
			name: "empty",
		},
		{
			name: "one action",
			actions: []*bazel.Action{
				{Arguments: []string{"a", "b", "-c", "fileA"}},
			},
			wantCommands: []*compdb.Command{
				{File: "fileA", Arguments: []string{"a", "b", "-c", "fileA"}},
			},
		},
		{
			name: "missing source",
			actions: []*bazel.Action{
				{Arguments: []string{"a", "b", "c"}},
			},
			wantErr: compdb.ErrSource,
		},
		{
			name: "two actions",
			actions: []*bazel.Action{
				{Arguments: []string{"a1", "b1", "-c", "fileA"}},
				{Arguments: []string{"a2", "b2", "-c", "fileB"}},
			},
			wantCommands: []*compdb.Command{
				{File: "fileA", Arguments: []string{"a1", "b1", "-c", "fileA"}},
				{File: "fileB", Arguments: []string{"a2", "b2", "-c", "fileB"}},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var aset *bazel.ActionSet
			if len(tc.actions) != 0 {
				aset = &bazel.ActionSet{Actions: tc.actions}
			}

			commands, err := compdb.MakeCommands(aset)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("unexpected error %v; want %v", err, tc.wantErr)
			}
			opts := []cmp.Option{
				cmpopts.IgnoreFields(compdb.Command{}, "Directory"),
			}
			if d := cmp.Diff(tc.wantCommands, commands, opts...); d != "" {
				t.Errorf("mismatch (-want +got):\n%s", d)
			}
		})
	}
}

func TestPrint(t *testing.T) {
	commands := []*compdb.Command{
		{File: "fileA", Directory: "dirA", Arguments: []string{"A1", "A2"}},
		{File: "fileB", Directory: "dirB", Arguments: []string{"B1", "B2"}},
	}
	var b strings.Builder

	err := compdb.Print(&b, commands)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	want := `[
  {
    "file": "fileA",
    "directory": "dirA",
    "arguments": [
      "A1",
      "A2"
    ]
  },
  {
    "file": "fileB",
    "directory": "dirB",
    "arguments": [
      "B1",
      "B2"
    ]
  }
]
`
	if d := cmp.Diff(want, b.String()); d != "" {
		t.Errorf("mismatch (-want +got):\n%s", d)
	}
}
