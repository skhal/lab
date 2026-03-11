// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package todo

import "fmt"

// TodoError contains a reference to the invalid todo-comment in the file. It
// holds the file name, line number, and the line content.
type TodoError struct {
	File string // file name
	Line int    // line number
	Text string // line text
}

// Error implements [builtin.error] interface.
func (e *TodoError) Error() string {
	return fmt.Sprintf("%s:%d %s", e.File, e.Line, e.Text)
}

// Is implements interface for [errors.Is].
func (e *TodoError) Is(target error) bool {
	x, ok := target.(*TodoError)
	if !ok {
		return false
	}
	return *e == *x
}
