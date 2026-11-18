// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tests_test

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/skhal/lab/go/tests"
)

func TestCommandFaker(t *testing.T) {
	testCases := []struct {
		name       string
		helper     string
		cmd        string
		args       []string
		wantCmd    string
		wantArgs   []string
		wantErr    bool
		wantErrMsg string // substring expected in helper stderr
	}{
		{
			name:     "helper exit ok",
			helper:   "testExitOk",
			cmd:      "test-command",
			args:     []string{"a", "b"},
			wantCmd:  "test-command",
			wantArgs: []string{"a", "b"},
		},
		{
			name:     "helper exit fail",
			helper:   "testExitFail",
			cmd:      "test-command",
			args:     []string{"a", "b"},
			wantCmd:  "test-command",
			wantArgs: []string{"a", "b"},
			wantErr:  true,
		},
		{
			name:       "invalid helper",
			helper:     "non-existent helper",
			cmd:        "test-command",
			args:       []string{"a", "b"},
			wantCmd:    "test-command",
			wantArgs:   []string{"a", "b"},
			wantErr:    true,
			wantErrMsg: "unsupported helper",
		},
		{
			name:     "helper gets args",
			helper:   "testHelperArgs",
			cmd:      "test-command",
			args:     []string{"a", "b"},
			wantCmd:  "test-command",
			wantArgs: []string{"a", "b"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			faker := tests.NewCommandFaker(tc.wantCmd, tc.wantArgs...)
			command := faker.Fake(t, tc.helper)
			var stderr bytes.Buffer

			cmd := command(tc.cmd, tc.args...)
			cmd.Stderr = &stderr
			err := cmd.Run()

			if tc.wantErr {
				if err == nil {
					t.Errorf("missing error from exec.Cmd.Run()")
				} else if _, ok := errors.AsType[*exec.ExitError](err); !ok {
					t.Errorf("unexpected error %T (%v); want *exec.ExitError", err, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error '%v'; stderr: %s", err, stderr.String())
				}
			}
			if !strings.Contains(stderr.String(), tc.wantErrMsg) {
				t.Errorf("stderr %q does not contain %q", stderr.String(), tc.wantErrMsg)
			}
		})
	}
}

var helpers = map[string]tests.HelperFunc{
	"testExitFail":   testExitFail,
	"testExitOk":     testExitOk,
	"testHelperArgs": testHelperArgs,
}

func TestMain(m *testing.M) {
	main := tests.NewTestMain(helpers)
	main(m)
}

func testExitFail([]string) error { return tests.ErrExitNoOutput }

func testExitOk([]string) error { return nil }

func testHelperArgs(args []string) error {
	want := []string{"a", "b"}

	if d := cmp.Diff(want, args); d != "" {
		return fmt.Errorf("arguments mismatch (-want +got):\n%s", d)
	}
	return nil
}
