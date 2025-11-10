// Copyright 2025 Samvel Khalatyan. All rights reserved.

package test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"iter"
	"os/exec"
	"path/filepath"
)

// EventID identifies an Event. It includes the Event's package and test names.
type EventID string

// Tester runs `go test` on packages and groups events by event ids. It also
// keeps track of failed tests for further analysis.
type Tester struct {
	events map[EventID][]*Event
	fails  []EventID
}

// NewTester creates a tester, ready for testing packages.
func NewTester() *Tester {
	return &Tester{
		events: make(map[EventID][]*Event),
	}
}

// Test runs `go test` on a single package and collects test output, groupped
// by test id. It keeps track of failed tests.
func (t *Tester) Test(pkg string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "test", "-json", pkg)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	defer cmd.Wait()
	for id, e := range decodeEvents(stdout) {
		if e.Action == ActionFail {
			t.fails = append(t.fails, id)
		}
		ee := t.events[id]
		t.events[id] = append(ee, e)
	}
	return nil
}

func decodeEvents(r io.Reader) iter.Seq2[EventID, *Event] {
	return func(yield func(EventID, *Event) bool) {
		dec := json.NewDecoder(r)
		for {
			e := new(Event)
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
func NewEventID(e *Event) EventID {
	id := e.Test
	if e.Package != "" {
		id = filepath.Join(e.Package, e.Test)
	}
	return EventID(id)
}

// FailedTest holds failed test package, name and output of `go test` for the
// the given test.
type FailedTest struct {
	Package string
	Test    string

	Output []byte
}

func newFailedTest(ee []*Event) *FailedTest {
	buf := new(bytes.Buffer)
	var pkg, test string
	for _, e := range ee {
		switch e.Action {
		case ActionFail:
			pkg = e.Package
			test = e.Test
		case ActionOutput:
			buf.WriteString(e.Output)
		}
	}
	return &FailedTest{
		Package: pkg,
		Test:    test,
		Output:  buf.Bytes(),
	}
}
