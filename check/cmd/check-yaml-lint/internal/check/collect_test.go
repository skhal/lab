// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check_test

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/skhal/lab/check/cmd/check-yaml-lint/internal/check"
)

func TestCollect(t *testing.T) {
	tests := []struct {
		name       string
		fs         fs.FS
		files      check.Files
		wantGroups map[check.Base]check.Files
		wantFiles  check.Files
	}{
		{
			name: "empty files list",
			fs:   fstest.MapFS{},
		},
		{
			name:      "non-existent file",
			fs:        fstest.MapFS{},
			files:     []string{"a.yml"},
			wantFiles: []string{"a.yml"},
		},
		{
			name: "not ansible file",
			fs: fstest.MapFS{
				"a.yml": {},
			},
			files:     []string{"a.yml"},
			wantFiles: []string{"a.yml"},
		},
		{
			name: "ansible root project file",
			fs: fstest.MapFS{
				"ansible.cfg": {},
				"a.yml":       {},
			},
			files: []string{"a.yml"},
			wantGroups: map[check.Base]check.Files{
				{Path: ".", Kind: check.BaseProject}: []string{"a.yml"},
			},
		},
		{
			name: "ansible root project sub-file",
			fs: fstest.MapFS{
				"ansible.cfg": {},
				"A/a.yml":     {},
			},
			files: []string{"A/a.yml"},
			wantGroups: map[check.Base]check.Files{
				{Path: ".", Kind: check.BaseProject}: []string{"A/a.yml"},
			},
		},
		{
			name: "ansible project file",
			fs: fstest.MapFS{
				"A/ansible.cfg": {},
				"A/a.yml":       {},
			},
			files: []string{"A/a.yml"},
			wantGroups: map[check.Base]check.Files{
				{Path: "A", Kind: check.BaseProject}: []string{"a.yml"},
			},
		},
		{
			name: "ansible root collection file",
			fs: fstest.MapFS{
				"galaxy.yml": {},
				"a.yml":      {},
			},
			files: []string{"a.yml"},
			wantGroups: map[check.Base]check.Files{
				{Path: ".", Kind: check.BaseCollection}: []string{"a.yml"},
			},
		},
		{
			name: "ansible root collection sub-file",
			fs: fstest.MapFS{
				"galaxy.yml": {},
				"A/a.yml":    {},
			},
			files: []string{"A/a.yml"},
			wantGroups: map[check.Base]check.Files{
				{Path: ".", Kind: check.BaseCollection}: []string{"A/a.yml"},
			},
		},
		{
			name: "ansible collection file",
			fs: fstest.MapFS{
				"A/galaxy.yml": {},
				"A/a.yml":      {},
			},
			files: []string{"A/a.yml"},
			wantGroups: map[check.Base]check.Files{
				{Path: "A", Kind: check.BaseCollection}: []string{"a.yml"},
			},
		},
		{
			name: "ansible project and collection files",
			fs: fstest.MapFS{
				"A/ansible.cfg": {},
				"A/a.yml":       {},
				"B/galaxy.yml":  {},
				"B/b.yml":       {},
			},
			files: []string{"A/a.yml", "B/b.yml"},
			wantGroups: map[check.Base]check.Files{
				{Path: "A", Kind: check.BaseProject}:    []string{"a.yml"},
				{Path: "B", Kind: check.BaseCollection}: []string{"b.yml"},
			},
		},
		{
			name: "ansible and non-ansible files",
			fs: fstest.MapFS{
				"A/ansible.cfg": {},
				"A/a.yml":       {},
				"B/galaxy.yml":  {},
				"B/b.yml":       {},
				"C/c.yml":       {},
			},
			files: []string{"A/a.yml", "B/b.yml", "C/c.yml"},
			wantGroups: map[check.Base]check.Files{
				{Path: "A", Kind: check.BaseProject}:    []string{"a.yml"},
				{Path: "B", Kind: check.BaseCollection}: []string{"b.yml"},
			},
			wantFiles: []string{"C/c.yml"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			groups, files := check.Collect(tc.fs, tc.files)

			if d := cmp.Diff(tc.wantGroups, groups, cmpopts.EquateEmpty()); d != "" {
				t.Errorf("collected groups mismatch (-want +got):\n%s", d)
			}
			if d := cmp.Diff(tc.wantFiles, files, cmpopts.EquateEmpty()); d != "" {
				t.Errorf("not collected files mismatch (-want +got):\n%s", d)
			}
		})
	}
}
