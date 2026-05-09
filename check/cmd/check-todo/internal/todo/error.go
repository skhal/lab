// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package todo

import (
	"errors"
	"fmt"
)

// ErrTodo means a todo comment is malformed. It must have a reference to an
// isuse.
var ErrTodo = errors.New("todo error")

// TodoError contains a reference to the invalid todo-comment in the file. It
// holds the file name, line number, and the line content.
type TodoError struct {
	File string // file name
	Text string // todo-comment
	Line int    // line number
}

// Error implements [builtin.error] interface.
func (e *TodoError) Error() string {
	return fmt.Sprintf("%s:%d %s", e.File, e.Line, e.Text)
}

// Is implements interface for [errors.Is].
func (*TodoError) Is(err error) bool {
	return err == ErrTodo
}
