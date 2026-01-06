// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
//go:generate stringer -type=Action

package test

import (
	"encoding/json"
	"strings"
	"time"
)

// Event is a TestEvent from test2json (see: go doc test2json).
type TestEvent struct {
	// Time is the time of the event, omitted for cached test results.
	Time        time.Time     `json:",omitzero"`  // RFC3339
	Action      Action        `json:",omitzero"`  // a JSON stream begins with ActionStart
	Package     string        `json:",omitempty"` // identifies the package under test
	Test        string        `json:",omitempty"`
	Elapsed     time.Duration `json:",omitzero"`
	Output      string        `json:",omitempty"`
	FailedBuild string        `json:",omitempty"`
}

func (e *TestEvent) MarshalJSON() ([]byte, error) {
	// Use an alias to the Event to avoid MarshalJSON() infinite loop - json
	// package recursively traverses the structure and uses MarshalJSON for the
	// fields of types with such a method.
	type EventAlias TestEvent
	evt := struct {
		Elapsed float64 `json:",omitempty"`
		*EventAlias
	}{
		Elapsed:    e.Elapsed.Seconds(),
		EventAlias: (*EventAlias)(e),
	}
	return json.Marshal(evt)
}

func (e *TestEvent) UnmarshalJSON(b []byte) error {
	// Use an alias to the Event to avoid MarshalJSON() infinite loop - json
	// package recursively traverses the structure and uses MarshalJSON for the
	// fields of types with such a method.
	type EventAlias TestEvent
	evt := &struct {
		Elapsed float64
		*EventAlias
	}{
		EventAlias: (*EventAlias)(e),
	}
	if err := json.Unmarshal(b, evt); err != nil {
		return err
	}
	e.Elapsed = time.Duration(evt.Elapsed * float64(time.Second))
	return nil
}

// Actin is the event state (see go doc test2json)
type Action int

const (
	ActionUnspecified Action = iota
	ActionStart              // the test binary is about to be executed
	ActionRun                // the test run started
	ActionPause              // the test paused
	ActionContinue           // the test run continued
	ActionPass               // the test passed
	ActionBenchmark          // the benchmarks logs and did not fail
	ActionFail               // the test or benchmark failed
	ActionOutput             // test output prints
	ActionSkip               // the test skipped or no tests in the package
)

// MarshalJSON implements json.Marshaler interface.
func (a Action) MarshalJSON() ([]byte, error) {
	var s string
	switch a {
	default:
		s = "unspecified"
	case ActionStart:
		s = "start"
	case ActionRun:
		s = "run"
	case ActionPause:
		s = "pause"
	case ActionContinue:
		s = "cont"
	case ActionPass:
		s = "pass"
	case ActionBenchmark:
		s = "bench"
	case ActionFail:
		s = "fail"
	case ActionOutput:
		s = "output"
	case ActionSkip:
		s = "skip"
	}
	return json.Marshal(s)
}

// UnmsarshalJSON implements json.Unmarshaler interface.
func (a *Action) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	default:
		*a = ActionUnspecified
	case "start":
		*a = ActionStart
	case "run":
		*a = ActionRun
	case "pause":
		*a = ActionPause
	case "cont":
		*a = ActionContinue
	case "pass":
		*a = ActionPass
	case "bench":
		*a = ActionBenchmark
	case "fail":
		*a = ActionFail
	case "output":
		*a = ActionOutput
	case "skip":
		*a = ActionSkip
	}
	return nil
}
