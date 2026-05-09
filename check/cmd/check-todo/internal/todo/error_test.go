// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package todo_test

import (
	"errors"
	"testing"

	"github.com/skhal/lab/check/cmd/check-todo/internal/todo"
)

func TestTodoError_Is(t *testing.T) {
	err := &todo.TodoError{}

	if !errors.Is(err, todo.ErrTodo) {
		t.Errorf("want TodoError to be equivalent to ErrTodo")
	}
}
