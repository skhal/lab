// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check_test

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/skhal/lab/check/cmd/check-go-test/internal/build"
	"github.com/skhal/lab/check/cmd/check-go-test/internal/check"
	"github.com/skhal/lab/check/cmd/check-go-test/internal/test"
)

func TestJSONUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		b       string
		want    check.Event
		wantErr error
	}{
		{
			name:    "not json",
			b:       "test",
			wantErr: check.ErrNotJSON,
		},
		{
			name: "build event",
			b:    `{"ImportPath":"test","Action":"build-output"}`,
			want: (*check.BuildEvent)(
				&build.Event{ImportPath: "test", Action: build.ActionOutput},
			),
		},
		{
			name: "test event",
			b:    `{"Package":"test","Action":"output"}`,
			want: &check.TestEvent{Event: &test.Event{Package: "test", Action: test.ActionOutput}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := check.JSONUnmarshal([]byte(tc.b))
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("check.JSONUnmarshal() unexpected error %q; want %q", err, tc.wantErr)
				t.Log(tc.b)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("check.JSONUnmarshal() mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestTester_Test(t *testing.T) {
	tests := []struct {
		name     string
		helper   string
		command  *testCommand
		pkgs     []string
		coverage check.Coverage
		wantErr  error
	}{
		{
			name:    "no packages",
			helper:  "nullOutput",
			command: newTestCommand("go", "test", "-json", "-vet=all"),
		},
		{
			name:    "build pass",
			helper:  "packageBuildPass",
			command: newTestCommand("go", "test", "-json", "-vet=all", "./test"),
			pkgs:    []string{"./test"},
		},
		{
			name:    "build fail",
			helper:  "packageBuildFail",
			command: newTestCommand("go", "test", "-json", "-vet=all", "./test"),
			pkgs:    []string{"./test"},
			wantErr: check.ErrBuild,
		},
		{
			name:    "test pass",
			helper:  "packageTestPass",
			command: newTestCommand("go", "test", "-json", "-vet=all", "./test"),
			pkgs:    []string{"./test"},
		},
		{
			name:    "test fail",
			helper:  "packageTestFail",
			command: newTestCommand("go", "test", "-json", "-vet=all", "./test"),
			pkgs:    []string{"./test"},
			wantErr: check.ErrTest,
		},
		{
			name:     "coverage above threshold",
			helper:   "packageTestCoverage",
			command:  newTestCommand("go", "test", "-json", "-vet=all", "-cover", "./test"),
			pkgs:     []string{"./test"},
			coverage: 25,
		},
		{
			name:     "coverage below threshold",
			helper:   "packageTestCoverage",
			command:  newTestCommand("go", "test", "-json", "-vet=all", "-cover", "./test"),
			pkgs:     []string{"./test"},
			coverage: 75,
			wantErr:  check.ErrCoverage,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			opts := []check.Opt{
				check.WithCoverage(float64(tc.coverage)),
				check.WithCommand(tc.command.New(t, tc.helper)),
			}
			tester := check.NewTester(opts...)

			err := tester.Test(tc.pkgs)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("(*Tester).Test() unexpected error %v; want %v", err, tc.wantErr)
			}
		})
	}
}

const cmdHelperEnvKey = "GO_TEST_RUN_HELPER"

type testCommand struct {
	Name string
	Args []string
}

func newTestCommand(name string, args ...string) *testCommand {
	return &testCommand{name, args}
}

func (c *testCommand) New(t *testing.T, helper string) func(string, ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		t.Helper()
		if c.Name != name {
			t.Fatalf("unexpected command %s; want %s", name, c.Name)
		}
		opts := []cmp.Option{
			cmpopts.SortSlices(strings.Compare),
		}
		if d := cmp.Diff(c.Args, args, opts...); d != "" {
			t.Fatalf("args mismatch (-want +got):\n%s", d)
		}
		testApp, err := os.Executable()
		if err != nil {
			t.Fatal(err)
		}
		cmd := exec.Command(testApp, helper)
		cmd.Env = append(cmd.Environ(), fmt.Sprintf("%s=1", cmdHelperEnvKey))
		return cmd
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	if os.Getenv(cmdHelperEnvKey) == "" {
		os.Exit(m.Run())
	}
	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "missing helper")
		os.Exit(1)
	}
	args := flag.Args()
	name := args[0]
	h, ok := helpers[name]
	if !ok {
		fmt.Fprintf(os.Stderr, "unsupported helper %s\n", name)
		os.Exit(1)
	}
	if err := h(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var helpers = map[string]func() error{
	// keep-sorted start
	"nullOutput":          nullOutput,
	"packageBuildFail":    packageBuildFail,
	"packageBuildPass":    packageBuildPass,
	"packageTestCoverage": packageTestCoverage,
	"packageTestFail":     packageTestFail,
	"packageTestPass":     packageTestPass,
	// keep-sorted end
}

func nullOutput() error {
	return nil
}

func packageBuildFail() error {
	newEvent := func(a build.Action, out string) *build.Event {
		return &build.Event{ImportPath: "test/foo", Action: a, Output: out}
	}
	events := []any{
		newEvent(build.ActionOutput, "test output"),
		newEvent(build.ActionFail, ""),
	}
	for _, e := range events {
		b, err := json.Marshal(e)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	}
	return nil
}

func packageBuildPass() error {
	newEvent := func(a build.Action, out string) *build.Event {
		return &build.Event{ImportPath: "test/foo", Action: a, Output: out}
	}
	events := []any{
		newEvent(build.ActionOutput, "test output"),
	}
	for _, e := range events {
		b, err := json.Marshal(e)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	}
	return nil
}

func packageTestCoverage() error {
	newEvent := func(a test.Action, out string) *test.Event {
		return &test.Event{Package: "test", Test: "foo", Action: a, Output: out}
	}
	events := []any{
		newEvent(test.ActionOutput, "test output"),
		newEvent(test.ActionOutput, "coverage: 50.0% of statements"),
		newEvent(test.ActionPass, ""),
	}
	for _, e := range events {
		b, err := json.Marshal(e)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	}
	return nil
}

func packageTestFail() error {
	newEvent := func(a test.Action, out string) *test.Event {
		return &test.Event{Package: "test", Test: "foo", Action: a, Output: out}
	}
	events := []any{
		newEvent(test.ActionOutput, "test output"),
		newEvent(test.ActionFail, ""),
	}
	for _, e := range events {
		b, err := json.Marshal(e)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	}
	return nil
}

func packageTestPass() error {
	newEvent := func(a test.Action, out string) *test.Event {
		return &test.Event{Package: "test", Test: "foo", Action: a, Output: out}
	}
	events := []any{
		newEvent(test.ActionOutput, "test output"),
		newEvent(test.ActionPass, ""),
	}
	for _, e := range events {
		b, err := json.Marshal(e)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	}
	return nil
}
