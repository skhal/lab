// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "errors"

type command struct{}

func newCommand() *command {
	return new(command)
}

// Run executes the command.
func (cmd *command) Run() error {
	return errors.New("not implemented")
}
