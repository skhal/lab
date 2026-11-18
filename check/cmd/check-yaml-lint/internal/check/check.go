// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	labslices "github.com/skhal/lab/go/slices"
)

var (
	// ErrEmptyProject means the Ansible project or collection does not have files
	// to check.
	ErrEmptyProject = errors.New("empty project")

	// ErrCheckProject means there was error running ansible-lint.
	ErrCheckProject = errors.New("lint ansible project error")

	// ErrCheckFiles means there was error running yamlfmt.
	ErrCheckFiles = errors.New("format files error")
)

type checkProjectError struct {
	err   error
	base  Base
	files Files
}

func newCheckProjectError(base Base, files Files, err error) *checkProjectError {
	return &checkProjectError{base: base, files: files, err: err}
}

// Error implements [builtin.error] interface.
func (cpe *checkProjectError) Error() string {
	files := strings.Join(cpe.files, ",")
	return fmt.Sprintf("%s %s %s: %v", ErrCheckProject, cpe.base, files, cpe.err)
}

// Is implements [errors.Is] API.
func (cpe *checkProjectError) Is(err error) bool {
	return err == ErrCheckProject
}

type checkFilesError struct {
	err   error
	files Files
}

func newCheckFilesError(files Files, err error) *checkFilesError {
	return &checkFilesError{files: files, err: err}
}

// Error implements [builtin.error] interface.
func (cfe *checkFilesError) Error() string {
	files := strings.Join(cfe.files, ",")
	return fmt.Sprintf("%s %s: %v", ErrCheckFiles, files, cfe.err)
}

// Is implements [errors.Is] API.
func (cfe *checkFilesError) Is(err error) bool {
	return err == ErrCheckFiles
}

// CheckOpt is a [Check] option that configures checker in some way.
type CheckOpt func(*checker)

// CommandFunc represents [exec.Command] API.
type CommandFunc func(cmd string, args ...string) *exec.Cmd

// WithTestCommand configures checker with fake command for tests.
func WithTestCommand(f CommandFunc) CheckOpt {
	return func(chr *checker) {
		chr.commandFunc = f
	}
}

// Check runs Ansible linter on Ansible projects and collections and yamlfmt
// on everything else. It reports a consolidated error if any of the checks
// fail.
func Check(ansibleProjects map[Base]Files, nonAnsibleFiles Files, opts ...CheckOpt) error {
	var ee []error

	chr := newChecker(opts...)
	if err := chr.CheckProjects(ansibleProjects); err != nil {
		ee = append(ee, err)
	}
	if err := chr.CheckFiles(nonAnsibleFiles); err != nil {
		ee = append(ee, err)
	}

	return errors.Join(ee...)
}

type checker struct {
	commandFunc CommandFunc
}

func newChecker(opts ...CheckOpt) checker {
	chr := checker{commandFunc: exec.Command}
	for _, opt := range opts {
		opt(&chr)
	}
	return chr
}

// CheckProjects runs ansible-lint on Ansible projects and collections.
func (c checker) CheckProjects(projects map[Base]Files) error {
	var ee []error
	sortByPath := func(a, b Base) int {
		return strings.Compare(a.Path, b.Path)
	}
	for _, base := range slices.SortedFunc(maps.Keys(projects), sortByPath) {
		files := projects[base]
		if err := c.checkProject(base, files); err != nil {
			ee = append(ee, err)
		}
	}
	return errors.Join(ee...)
}

func (chr checker) checkProject(base Base, files Files) error {
	if len(files) == 0 {
		return fmt.Errorf("%s: %w", base, ErrEmptyProject)
	}

	mapFunc := func(f string) string { return filepath.Join(base.Path, f) }
	files = labslices.MapFunc(files, mapFunc)

	args := append([]string{
		"run",
		"ansible-lint",
		"--project-dir",
		base.Path,
		"--format",
		"brief",
		"-q", // quiet
		"--fix",
		"all",
		"--offline", // no schema refresh, do not install requirements.yml
	}, files...)

	cmd := chr.commandFunc("poetry", args...)

	cmd.Env = append(cmd.Environ(), "ANSIBLE_DEPRECATION_WARNINGS=False")

	if base.Kind == BaseProject {
		p := filepath.Join(base.Path, ansibleProjectConfig)
		cmd.Env = append(cmd.Environ(), fmt.Sprintf("ANSIBLE_CONFIG=%s", p))
	}

	// pass the output through for pre-commit to control it.
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return newCheckProjectError(base, files, err)
	}
	return nil
}

// CheckFiles runs yamlfmt on non-Ansible files.
func (chr checker) CheckFiles(files Files) error {
	if len(files) == 0 {
		return nil
	}

	cmd := chr.commandFunc("yamlfmt", files...)

	// pass the output through for pre-commit to control it.
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return newCheckFilesError(files, err)
	}
	return nil
}
