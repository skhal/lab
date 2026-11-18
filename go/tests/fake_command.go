// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tests

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// ErrExitNoOutput exits helper with exit code 1 but does not print error.
var ErrExitNoOutput = errors.New("exit without output")

const helperEnvKey = "GO_TEST_RUN_HELPER"

// CommandFaker fakes [exec.Command] in tests.
type CommandFaker struct {
	wantName string
	wantArgs []string
}

// NewCommandFaker constructs a faker with expected command and arguments.
func NewCommandFaker(wantName string, wantArgs ...string) *CommandFaker {
	return &CommandFaker{wantName, wantArgs}
}

// Fake creates a fake [exec.Command]. Inject it instead of [exec.Command] from
// the tests.
//
// The fake validates used command and arguments, when the code under test calls
// the fake instead of exec.Command(command, args...). It uses
// [testing.T.Fatalf] if validation fails. Therefore, call it from the test
// goroutine.
//
// It returns exec.Cmd with environment marker to run the test executable in
// helper-mode. Do not replace exec.Cmd.Env in the code under test. See
// [NewTestMain].
func (fc *CommandFaker) Fake(t *testing.T, helper string) func(string, ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		t.Helper()

		if fc.wantName != name {
			t.Fatalf("unexpected command %q; want %q", name, fc.wantName)
		}

		if d := cmp.Diff(fc.wantArgs, args); d != "" {
			t.Fatalf("args mismatch (-want +got):\n%s", d)
		}

		testApp, err := os.Executable()
		if err != nil {
			t.Fatal(err)
		}

		args = append([]string{helper}, args...)

		cmd := exec.Command(testApp, args...)
		cmd.Env = append(cmd.Environ(), fmt.Sprintf("%s=1", helperEnvKey))

		return cmd
	}
}

// HelperFunc is the helper function in the _test.go that is supplied to
// [CommandFaker.Fake] and would be called by [NewTestMain].
type HelperFunc func(args []string) error

// NewTestMain generates a TestMain that delegates the call to the helper.
func NewTestMain(helpers map[string]HelperFunc) func(m *testing.M) {
	return func(m *testing.M) {
		flag.Parse()

		if os.Getenv(helperEnvKey) == "" {
			os.Exit(m.Run())
		}

		args := flag.Args()
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "missing helper arguments")
			os.Exit(1)
		}

		helperName := args[0]

		helper, ok := helpers[helperName]
		if !ok {
			fmt.Fprintln(os.Stderr, "unsupported helper", helperName)
			os.Exit(1)
		}

		if err := helper(args[1:]); err != nil {
			if !errors.Is(err, ErrExitNoOutput) {
				fmt.Fprintf(os.Stderr, "%s: %s\n", helperName, err)
			}
			os.Exit(1)
		}

		os.Exit(0)
	}
}
