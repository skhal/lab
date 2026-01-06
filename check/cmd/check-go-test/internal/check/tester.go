// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"iter"
	"os/exec"
	"path/filepath"

	"github.com/skhal/lab/check/cmd/check-go-test/internal/test"
)

// EventID identifies an Event. It includes the Event's package and test names.
type EventID string

// Tester runs `go test` on packages and groups events by event ids. It also
// keeps track of failed tests for further analysis.
type Tester struct {
	events map[EventID][]*test.TestEvent
	fails  []EventID
}

// NewTester creates a tester, ready for testing packages.
func NewTester() *Tester {
	return &Tester{
		events: make(map[EventID][]*test.TestEvent),
	}
}

// Test runs `go test` on a single package and collects test output, grouped
// by test id. It keeps track of failed tests.
func (t *Tester) Test(pkg string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "test", "-json", "-vet=all", pkg)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	defer cmd.Wait()
	for id, e := range decodeEvents(stdout) {
		if e.Action == test.ActionFail {
			t.fails = append(t.fails, id)
		}
		ee := t.events[id]
		t.events[id] = append(ee, e)
	}
	return nil
}

func decodeEvents(r io.Reader) iter.Seq2[EventID, *test.TestEvent] {
	return func(yield func(EventID, *test.TestEvent) bool) {
		dec := json.NewDecoder(r)
		for {
			e := new(test.TestEvent)
			if err := dec.Decode(e); err == io.EOF {
				break
			} else if err != nil {
				break
			}
			if !yield(NewEventID(e), e) {
				break
			}
		}
	}
}

// VisitFails calls f on failed tests.
func (t *Tester) VisitFails(f func(*FailedTest)) {
	for _, id := range t.fails {
		ft := newFailedTest(t.events[id])
		f(ft)
	}
}

// NewEventID constructs an EventID for a given event.
func NewEventID(e *test.TestEvent) EventID {
	id := e.Test
	if e.Package != "" {
		id = filepath.Join(e.Package, e.Test)
	}
	return EventID(id)
}

// FailedTest holds failed test package, name and output of `go test` for a
// given test.
type FailedTest struct {
	Package string
	Test    string

	Output []byte
}

func newFailedTest(ee []*test.TestEvent) *FailedTest {
	buf := new(bytes.Buffer)
	var pkg, t string
	for _, e := range ee {
		switch e.Action {
		case test.ActionFail:
			pkg = e.Package
			t = e.Test
		case test.ActionOutput:
			buf.WriteString(e.Output)
		}
	}
	return &FailedTest{
		Package: pkg,
		Test:    t,
		Output:  buf.Bytes(),
	}
}
