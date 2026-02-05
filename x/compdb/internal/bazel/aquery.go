// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Bazel provides integration points with Bazel build system such as running
// actions query.
package bazel

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

const bazelCommand = "bazel"

// ActionSet is the actions graph.
//
// Ref: https://github.com/bazelbuild/bazel/blob/e2cf1361d07a4e6ecdf241addbe2891111cf08b4/src/main/protobuf/analysis_v2.proto#L25
type ActionSet struct {
	Actions []*Action // All actions in the graph.
}

// Action describes a single action in the actions graph.
type Action struct {
	Arguments []string // A command associated with the action.
}

// Aquery runs an actions query (aquery) for a given target.
func Aquery(target string) (*ActionSet, error) {
	cmd := exec.Command(bazelCommand, newArgs(target)...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	aset := new(ActionSet)
	if err := json.Unmarshal(b, aset); err != nil {
		return nil, err
	}
	return aset, nil
}

var standardArgs = []string{
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

func newArgs(target string) []string {
	// Ref: https://github.com/redpanda-data/redpanda/blob/70d49ac8d266e832acf805d718b5df634b58ae94/bazel/compilation_database_generator/main.go#L158
	return append([]string{
		"aquery",
		fmt.Sprintf(`mnemonic('CppCompile',deps(%s))`, target),
	}, standardArgs...)
}
