// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Compdb generates LLVM compilation database for Bazel targets.
package compdb

import (
	"encoding/json"
	"errors"
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

// Run generates a compilation database for a list of targets.
func Run(targets []string) error {
	commands, err := genCommands(targets)
	if err != nil {
		return err
	}
	return Print(commands)
}

func genCommands(targets []string) ([]*Command, error) {
	aset, err := bazel.Aquery(targets)
	if err != nil {
		return nil, err
	}
	return makeCommands(aset)
}

func makeCommands(aset *bazel.ActionSet) ([]*Command, error) {
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
// right after the compile flag "-c".
//
// For eexample, the following action has source "foo.c":
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
func Print(cc []*Command) error {
	e := json.NewEncoder(os.Stdout)
	return e.Encode(cc)
}
