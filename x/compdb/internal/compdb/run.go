// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package compdb generates LLVM compilation database for Bazel targets.
package compdb

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/skhal/lab/x/compdb/internal/bazel"
)

// ErrSource is error to extraction the translation unit source.
var ErrSource = errors.New("failed to get source")

// Command defines the translation unit: the main source, where to build,
// and the command to compile.
//
// Ref: https://clang.llvm.org/docs/JSONCompilationDatabase.html
type Command struct {
	File      string   `json:"file"`      // the main translation unit source
	Directory string   `json:"directory"` // the compilation working directory
	Arguments []string `json:"arguments"` // the compile command
}

// GenCommands generates compile commands for a list of targets.
func GenCommands(targets []string) ([]*Command, error) {
	aset, err := bazel.Aquery(targets)
	if err != nil {
		return nil, err
	}
	return MakeCommands(aset)
}

// MakeCommands generates a list of compile commands from a set of Bazel
// actions. It returns an error if the actions do not include source file that
// is set by a pair of arguments ("-c", "file") or can't extract current
// working directory.
func MakeCommands(aset *bazel.ActionSet) ([]*Command, error) {
	if aset == nil || len(aset.Actions) == 0 {
		return nil, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	commands := make([]*Command, 0, len(aset.Actions))
	for _, a := range aset.Actions {
		s, err := GetSource(a)
		if err != nil {
			return nil, err
		}
		commands = append(commands, &Command{
			File:      s,
			Directory: wd,
			Arguments: a.Arguments,
		})
	}
	return commands, nil
}

const compileFlag = "-c"

// GetSource extracts the main translation unit source file from the action.
// It expects the source to be present in the [bazel.Action.Arguments] list,
// right after the compile flag "-c", e.g. "foo.c" source in:
//
//	bazel.Action{Arguments: [1, 2, 3, "-c", "foo.c", 4, 5]}
func GetSource(a *bazel.Action) (string, error) {
	for i, arg := range a.Arguments {
		if arg == compileFlag {
			i += 1
			if i == len(a.Arguments) {
				break
			}
			return a.Arguments[i], nil
		}
	}
	return "", ErrSource
}

// Print dumps commands in JSON format.
func Print(w io.Writer, cc []*Command) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(cc)
}
