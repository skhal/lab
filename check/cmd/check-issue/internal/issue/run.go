// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package issue

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
)

var (
	// ErrCheck indicates general error in the check.
	ErrCheck = errors.New("check error")

	// ErrNoIssue indicates missing issue
	ErrNoIssue = errors.New("missing issue")
)

// ReadFileFunc is used to read file.
type ReadFileFunc func(string) ([]byte, error)

// Config configures the check.
type Config struct {
	// ReadFileFn reads file.
	ReadFileFn ReadFileFunc
}

// NewConfig creates a new configuration with os.ReadFile function.
func NewConfig() *Config {
	return &Config{
		ReadFileFn: os.ReadFile,
	}
}

// Run executes the check. It is expected that the check will run as a
// commit-msg git-hook(1). Therefore there should be a single file, else it
// returns an error.
func Run(cfg *Config, files []string) error {
	if len(files) != 1 {
		return ErrCheck
	}
	return Check(cfg, files[0])
}

var (
	noissueRegexp = regexp.MustCompile(`^(?i)no_issue(?:: .*)?$`)
	// lint = keyword [ ":" ] issue .
	// keyword = close | closes | closed | fix | fixes | fised | resolve | resolves | resolved .
	// issue = [ owner "/" repo ] "#" number .
	// ref: https://docs.github.com/en/get-started/writing-on-github/working-with-advanced-formatting/using-keywords-in-issues-and-pull-requests
	//
	// Enforce sub-set of possible options:
	// - keyword: issue, close, fix
	// - issue: local or remote
	issueRegexp   = regexp.MustCompile(`^(?i)(?:issue|close|fix) (?:\w+/\w+)?#\d+$`)
)

// Check validates the file with commit message.
func Check(cfg *Config, file string) error {
	data, err := cfg.ReadFileFn(file)
	if err != nil {
		return err
	}
	s := bufio.NewScanner(bytes.NewReader(data))
	for s.Scan() {
		line := s.Text()
		if noissueRegexp.MatchString(line) {
			return nil
		}
		if issueRegexp.MatchString(line) {
			return nil
		}
	}
	return ErrNoIssue
}
