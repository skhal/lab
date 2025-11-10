// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
//go:generate stringer -type=Action

package test

import (
	"encoding/json"
	"strings"
	"time"
)

// Event is a TestEvent from test2json (see: go doc test2json).
type Event struct {
	// Time is the time of the event, omitted for cached test results.
	Time        time.Time // RFC3339
	Action      Action    // a JSON stream begins with ActionStart
	Package     string    // identifies the package under test
	Test        string
	Elapsed     Seconds
	Output      string
	FailedBuild string
}

type Seconds float64

// Actin is the event state (see go doc test2json)
type Action int

const (
	ActionUnspecified Action = iota
	ActionStart              // the test bainry is about to be executed
	ActionRun                // the test run started
	ActionPause              // the test paused
	ActionContinue           // the test run continued
	ActionPass               // the test passed
	ActionBenchmark          // the benchmaks logs and did not fail
	ActionFail               // the test or benchamrk failed
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
