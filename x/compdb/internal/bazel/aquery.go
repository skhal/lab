// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bazel provides integration points with Bazel build system such as
// running actions query.
package bazel

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	// ErrCommand means there is an error running Bazel.
	ErrCommand = errors.New("command error")

	// ErrJSON means there is an error unmarshaling JSON.
	ErrJSON = errors.New("json error")
)

const bazelCommand = "bazel"

// ActionSet is the actions graph.
//
// Ref: https://github.com/bazelbuild/bazel/blob/e2cf1361d07a4e6ecdf241addbe2891111cf08b4/src/main/protobuf/analysis_v2.proto#L25 -- NOLINT
type ActionSet struct {
	Actions []*Action // All actions in the graph.
}

// Action describes a single action in the actions graph.
type Action struct {
	Arguments []string // A command associated with the action.
}

// NewCmdFunc is a constructor for [exec.Cmd].
type NewCmdFunc func(name string, args ...string) *exec.Cmd

// Aquery runs an actions query (aquery) for targets.
func Aquery(targets []string) (*ActionSet, error) {
	return AqueryWithCommand(targets, exec.Command)
}

// AqueryWithCommand runs Bazel actions-query (aquery) for targets. It uses
// the newCmd function to create exec.Cmd to run Bazel.
func AqueryWithCommand(targets []string, newCmd NewCmdFunc) (*ActionSet, error) {
	if len(targets) == 0 {
		return nil, nil
	}
	cmd := newCmd(bazelCommand, newArgs(targets)...)
	cmd.Stderr = os.Stderr
	b, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrCommand, err)
	}
	aset := new(ActionSet)
	if err := json.Unmarshal(b, aset); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrJSON, err)
	}
	return aset, nil

}

var defaultArgs = []string{
	"--output=jsonproto",
	// Disable logging
	"--noshow_progress",
	"--ui_event_filters=-info",
	// Do not include the names of inputs and outputs (possibly large):
	// https://bazel.build/query/aquery#include-artifacts
	"--include_artifacts=false",
	// Do not put compiler args into a param file:
	// https://github.com/bazelbuild/bazel/issues/5163
	//
	// Ref: https://bazel.build/configure/windows
	"--features=-compiler_param_file",
	"--host_features=-compiler_param_file",
	// Disable check that targets only include headers from direct deps:
	// https://bazel.build/docs/bazel-and-cpp#toolchain-features
	//
	// Ref: https://llvm.org/docs/CodingStandards.html#library-layering
	"--features=-layering_check",
	"--host_features=-layering_check",
}

func newArgs(targets []string) []string {
	// Ref: https://github.com/redpanda-data/redpanda/blob/70d49ac8d266e832acf805d718b5df634b58ae94/bazel/compilation_database_generator/main.go#L158 -- NOLINT
	return append([]string{
		"aquery",
		fmt.Sprintf(`mnemonic('CppCompile',%s)`, unionDeps(targets)),
	}, defaultArgs...)
}

func unionDeps(targets []string) string {
	tt := make([]string, 0, len(targets))
	for _, t := range targets {
		tt = append(tt, fmt.Sprintf("deps(%s)", t))
	}
	// Use "+" instead of "union" to shorten args
	return strings.Join(tt, " + ")
}
