// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check

import (
	"path/filepath"

	"github.com/skhal/lab/check/cmd/check-go-test/internal/build"
	"github.com/skhal/lab/check/cmd/check-go-test/internal/test"
)

// EventID identifies an Event. It includes the Event's package and test names.
type EventID string

// Event corresponds to a single output line of Go test command. It can be a
// build or test event.
type Event interface {
	// Fail returns true if the event represents ActionFail.
	Fail() bool

	// ID returns an event identifier
	ID() EventID
}

// BuildEvent is an event from the build stage.
type BuildEvent build.Event

// Fail returns true if the event corresponds to the fail action.
func (e *BuildEvent) Fail() bool {
	return e.Action == build.ActionFail
}

// ID returns a unique event identifier within a go-test run.
func (e *BuildEvent) ID() EventID {
	id := e.ImportPath
	return EventID(id)
}

// TestEvent is an even from the test stage.
type TestEvent test.TestEvent

// Fail returns true if the event corresponds t the fail action.
func (e *TestEvent) Fail() bool {
	return e.Action == test.ActionFail
}

// ID returns a unique event identiifer within a go-test run.
func (e *TestEvent) ID() EventID {
	// filepath ignores empty elements (e.Package).
	return EventID(filepath.Join(e.Package, e.Test))
}
