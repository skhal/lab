// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bazel_test

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
	"github.com/skhal/lab/x/compdb/internal/bazel"
)

const envGoTestRunHelper = "GO_TEST_RUN_HELPER"

func TestAqueryWithCommand(t *testing.T) {
	tests := []struct {
		name        string
		helper      string
		targets     []string
		wantCmd     []string
		wantActions []*bazel.Action
		wantErr     error
	}{
		{
			name:   "no targets",
			helper: "cmdFail",
		},
		{
			name:    "one target",
			helper:  "cmdBazelAquery",
			targets: []string{"//a"},
			wantCmd: []string{
				"bazel",
				"--features=-compiler_param_file",
				"--features=-layering_check",
				"--host_features=-compiler_param_file",
				"--host_features=-layering_check",
				"--include_artifacts=false",
				"--noshow_progress",
				"--output=jsonproto",
				"--ui_event_filters=-info",
				"aquery",
				"mnemonic('CppCompile',deps(//a))",
			},
			wantActions: []*bazel.Action{
				{Arguments: []string{"one", "two", "three"}},
			},
		},
		{
			name:    "bazel error",
			helper:  "cmdFail",
			targets: []string{"//a"},
			wantCmd: []string{
				"bazel",
				"--features=-compiler_param_file",
				"--features=-layering_check",
				"--host_features=-compiler_param_file",
				"--host_features=-layering_check",
				"--include_artifacts=false",
				"--noshow_progress",
				"--output=jsonproto",
				"--ui_event_filters=-info",
				"aquery",
				"mnemonic('CppCompile',deps(//a))",
			},
			wantErr: bazel.ErrCommand,
		},
		{
			name:    "bazel error",
			helper:  "cmdNotJSON",
			targets: []string{"//a"},
			wantCmd: []string{
				"bazel",
				"--features=-compiler_param_file",
				"--features=-layering_check",
				"--host_features=-compiler_param_file",
				"--host_features=-layering_check",
				"--include_artifacts=false",
				"--noshow_progress",
				"--output=jsonproto",
				"--ui_event_filters=-info",
				"aquery",
				"mnemonic('CppCompile',deps(//a))",
			},
			wantErr: bazel.ErrJSON,
		},
		{
			name:    "union two targets",
			helper:  "cmdBazelAquery",
			targets: []string{"//a", "//b"},
			wantCmd: []string{
				"bazel",
				"--features=-compiler_param_file",
				"--features=-layering_check",
				"--host_features=-compiler_param_file",
				"--host_features=-layering_check",
				"--include_artifacts=false",
				"--noshow_progress",
				"--output=jsonproto",
				"--ui_event_filters=-info",
				"aquery",
				"mnemonic('CppCompile',deps(//a) + deps(//b))",
			},
			wantActions: []*bazel.Action{
				{Arguments: []string{"one", "two", "three"}},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wantCmd, wantArgs := func(cmd []string) (string, []string) {
				if len(cmd) < 1 {
					return "", nil
				}
				return cmd[0], cmd[1:]
			}(tc.wantCmd)
			cmdFunc := (&testCommand{tc.helper}).CmdFunc(t, wantCmd, wantArgs...)

			aset, err := bazel.AqueryWithCommand(tc.targets, cmdFunc)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("AqueryWithCommand() unexpected error %v; want %v", err, tc.wantErr)
			}
			var want *bazel.ActionSet
			if tc.wantActions != nil {
				want = &bazel.ActionSet{Actions: tc.wantActions}
			}
			if d := cmp.Diff(want, aset, cmpopts.EquateEmpty()); d != "" {
				t.Errorf("AqueryWithCommand() action set mismatch (-want +got):\n%s", d)
			}
		})
	}
}

type testCommand struct {
	helper string
}

func (tc *testCommand) CmdFunc(t *testing.T, wantName string, wantArgs ...string) bazel.NewCmdFunc {
	return func(name string, args ...string) *exec.Cmd {
		t.Helper()
		if name != wantName {
			t.Fatalf("unexpected command %q; want %q", name, wantName)
		}
		opts := []cmp.Option{
			cmpopts.SortSlices(strings.Compare),
		}
		if d := cmp.Diff(wantArgs, args, opts...); d != "" {
			t.Fatalf("args mismatch (-want +got):\n%s", d)
		}
		p, err := os.Executable()
		if err != nil {
			t.Fatal(err)
		}
		cmd := exec.Command(p, tc.helper)
		cmd.Env = append(cmd.Environ(), fmt.Sprintf("%s=1", envGoTestRunHelper))
		return cmd
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	if _, ok := os.LookupEnv(envGoTestRunHelper); !ok {
		os.Exit(m.Run())
	}
	name := flag.Args()[0]
	h, ok := helpers[name]
	if !ok {
		fmt.Fprintf(os.Stderr, "invalid helper %s\n", name)
		os.Exit(1)
	}
	h()
}

var helpers = map[string]func(){
	// keep-sorted start
	"cmdBazelAquery": cmdBazelAquery,
	"cmdFail":        func() { os.Exit(1) },
	"cmdNotJSON":     func() { fmt.Println("plain text") },
	// keep-sorted end
}

func cmdBazelAquery() {
	aset := bazel.ActionSet{
		Actions: []*bazel.Action{
			{Arguments: []string{"one", "two", "three"}},
		},
	}
	b, err := json.Marshal(aset)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(string(b))
}
