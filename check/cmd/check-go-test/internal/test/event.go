// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package test

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// TestEvent is a go-test output from test2json (see: go doc test2json).
type TestEvent struct {
	// Time is the event's time in RFC3339. It is missing for cached results.
	Time        time.Time     `json:",omitzero"`
	Package     string        `json:",omitempty"` // package under test
	Test        string        `json:",omitempty"` // test name
	Output      string        `json:",omitempty"` // output line
	FailedBuild string        `json:",omitempty"` // build output on error
	Action      Action        `json:",omitzero"`  // event kind
	Elapsed     time.Duration `json:",omitzero"`  // time since the test started
}

// MarshalJSON implements [json.Marshaler].
func (e *TestEvent) MarshalJSON() ([]byte, error) {
	// Use an alias to the Event to avoid MarshalJSON() infinite loop - json
	// package recursively traverses the structure and uses MarshalJSON for the
	// fields of types with such a method.
	type EventAlias TestEvent
	evt := struct {
		*EventAlias
		Elapsed float64 `json:",omitempty"`
	}{
		Elapsed:    e.Elapsed.Seconds(),
		EventAlias: (*EventAlias)(e),
	}
	return json.Marshal(evt)
}

// UnmarshalJSON implements [json.Unmarshaler].
func (e *TestEvent) UnmarshalJSON(b []byte) error {
	// Use an alias to the Event to avoid MarshalJSON() infinite loop - json
	// package recursively traverses the structure and uses MarshalJSON for the
	// fields of types with such a method.
	type EventAlias TestEvent
	evt := &struct {
		*EventAlias
		Elapsed float64
	}{
		EventAlias: (*EventAlias)(e),
	}
	if err := json.Unmarshal(b, evt); err != nil {
		return err
	}
	e.Elapsed = time.Duration(evt.Elapsed * float64(time.Second))
	return nil
}

// Action is the event state (see go doc test2json)
//
//go:generate stringer -type=Action -linecomment
type Action int

const (
	_ Action = iota
	// keep-sorted start
	ActionBenchmark // bench
	ActionContinue  // cont
	ActionFail      // fail
	ActionOutput    // output
	ActionPass      // pass
	ActionPause     // pause
	ActionRun       // run
	ActionSkip      // skip
	ActionStart     // start
	// keep-sorted end
)

// MarshalJSON implements [json.Marshaler].
func (a Action) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

var actionByName = make(map[string]Action)

func init() {
	for a := ActionBenchmark; a <= ActionStart; a++ {
		actionByName[a.String()] = a
	}
}

// UnmarshalJSON implements [json.Unmarshaler] interface.
func (a *Action) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	action, ok := actionByName[strings.ToLower(s)]
	if !ok {
		return fmt.Errorf("invalid action %s", s)
	}
	*a = action
	return nil
}
