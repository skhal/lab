// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/skhal/lab/check/cmd/check-yaml-lint/internal/check"
	"github.com/skhal/lab/go/tests"
)

func TestCheck(t *testing.T) {
	// Unset ANSIBLE_CONFIG inherited from the developer environment. Set it
	// first to restore the original value at the end of the test.
	t.Setenv("ANSIBLE_CONFIG", "")
	if err := os.Unsetenv("ANSIBLE_CONFIG"); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name     string
		helper   string
		projects map[check.Base]check.Files
		files    check.Files
		wantCmd  string
		wantArgs []string
		wantErr  error
	}{
		{
			name: "empty input",
		},
		// Ansible project
		{
			name:   "ansible project lint validate env",
			helper: "fakeAnsibleLintValidateEnv",
			projects: map[check.Base]check.Files{
				{Path: "A", Kind: check.BaseProject}: {"a.yml"},
			},
			wantCmd: "poetry",
			wantArgs: []string{
				"run",
				"ansible-lint",
				"--project-dir",
				"A",
				"--format",
				"brief",
				"-q",
				"--fix",
				"all",
				"--offline",
				"A/a.yml",
			},
		},
		{
			name:   "ansible project lint pass",
			helper: "fakeAnsibleLintExitOk",
			projects: map[check.Base]check.Files{
				{Path: "A", Kind: check.BaseProject}: {"a.yml"},
			},
			wantCmd: "poetry",
			wantArgs: []string{
				"run",
				"ansible-lint",
				"--project-dir",
				"A",
				"--format",
				"brief",
				"-q",
				"--fix",
				"all",
				"--offline",
				"A/a.yml",
			},
		},
		{
			name:   "ansible project lint fail",
			helper: "fakeAnsibleLintExitFail",
			projects: map[check.Base]check.Files{
				{Path: "A", Kind: check.BaseProject}: {"a.yml"},
			},
			wantCmd: "poetry",
			wantArgs: []string{
				"run",
				"ansible-lint",
				"--project-dir",
				"A",
				"--format",
				"brief",
				"-q",
				"--fix",
				"all",
				"--offline",
				"A/a.yml",
			},
			wantErr: check.ErrCheckProject,
		},
		// Ansible collection
		{
			name:   "ansible collection lint validate env",
			helper: "fakeAnsibleLintValidateNoConfig",
			projects: map[check.Base]check.Files{
				{Path: "A", Kind: check.BaseCollection}: {"a.yml"},
			},
			wantCmd: "poetry",
			wantArgs: []string{
				"run",
				"ansible-lint",
				"--project-dir",
				"A",
				"--format",
				"brief",
				"-q",
				"--fix",
				"all",
				"--offline",
				"A/a.yml",
			},
		},
		{
			name:   "ansible collection lint pass",
			helper: "fakeAnsibleLintExitOk",
			projects: map[check.Base]check.Files{
				{Path: "A", Kind: check.BaseCollection}: {"a.yml"},
			},
			wantCmd: "poetry",
			wantArgs: []string{
				"run",
				"ansible-lint",
				"--project-dir",
				"A",
				"--format",
				"brief",
				"-q",
				"--fix",
				"all",
				"--offline",
				"A/a.yml",
			},
		},
		{
			name:   "ansible collection lint fail",
			helper: "fakeAnsibleLintExitFail",
			projects: map[check.Base]check.Files{
				{Path: "A", Kind: check.BaseCollection}: {"a.yml"},
			},
			wantCmd: "poetry",
			wantArgs: []string{
				"run",
				"ansible-lint",
				"--project-dir",
				"A",
				"--format",
				"brief",
				"-q",
				"--fix",
				"all",
				"--offline",
				"A/a.yml",
			},
			wantErr: check.ErrCheckProject,
		},
		// non-Ansible files
		{
			name:     "files format pass",
			helper:   "fakeFormatExitOk",
			files:    check.Files{"a.yml"},
			wantCmd:  "yamlfmt",
			wantArgs: []string{"a.yml"},
		},
		{
			name:     "files format fail",
			helper:   "fakeFormatExitFail",
			files:    check.Files{"a.yml"},
			wantCmd:  "yamlfmt",
			wantArgs: []string{"a.yml"},
			wantErr:  check.ErrCheckFiles,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			faker := tests.NewCommandFaker(tc.wantCmd, tc.wantArgs...)
			command := faker.Fake(t, tc.helper)
			opts := []check.CheckOpt{check.WithTestCommand(command)}

			err := check.Check(tc.projects, tc.files, opts...)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("unexpected error '%v'; want '%v'", err, tc.wantErr)
			}
		})
	}
}
func TestMain(m *testing.M) {
	helpers := map[string]tests.HelperFunc{
		// keep-sorted start
		"fakeAnsibleLintExitFail":         helperExitFail,
		"fakeAnsibleLintExitOk":           helperExitOk,
		"fakeAnsibleLintValidateEnv":      fakeAnsibleLintValidateEnv,
		"fakeAnsibleLintValidateNoConfig": fakeAnsibleLintValidateNoConfig,
		"fakeFormatExitFail":              helperExitFail,
		"fakeFormatExitOk":                helperExitOk,
		// keep-sorted end
	}
	main := tests.NewTestMain(helpers)
	main(m)
}

func helperExitFail([]string) error { return tests.ErrExitNoOutput }

func helperExitOk([]string) error { return nil }

// fakeAnsibleLintValidateNoConfig fails if ANSIBLE_CONFIG is set. Ansible
// collections do not have ansible.cfg.
func fakeAnsibleLintValidateNoConfig([]string) error {
	envVar := "ANSIBLE_CONFIG"

	if val, ok := os.LookupEnv(envVar); ok {
		return fmt.Errorf("unexpected environment: %s is set to '%s'", envVar, val)
	}
	return nil
}

func fakeAnsibleLintValidateEnv(args []string) error {
	envVar := "ANSIBLE_CONFIG"

	gotVal, ok := os.LookupEnv(envVar)
	if !ok {
		return fmt.Errorf("invalid environment: missing %s", envVar)
	}

	idx := slices.Index(args, "--project-dir")
	if idx == -1 {
		return errors.New("invalid arguments: missing --project-dir")
	}
	idx += 1
	if idx == len(args) {
		return errors.New("invalid arguments: missing --project-dir value")
	}

	projectDir := args[idx]

	wantVal := filepath.Join(projectDir, "ansible.cfg")
	if wantVal != gotVal {
		return fmt.Errorf("unexpected environment: %s is '%s'; want '%s'", envVar, gotVal, wantVal)
	}

	return nil
}
