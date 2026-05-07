// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// ErrInvalidAction means an error in parsing a build action.
var ErrInvalidAction = errors.New("invalid action")

// Event descrbies a build state.
//
// See: go help buildjson.
type Event struct {
	ImportPath string `json:",omitempty"` // package path
	Output     string `json:",omitempty"` // build line
	Action     Action `json:",omitzero"`  // build event kind
}

// Action describes a build event. See: go help buildjson.
//
//go:generate stringer -type=Action -linecomment
type Action int

const (
	_ Action = iota
	// keep-sorted start
	ActionFail   // build-fail
	ActionOutput // build-output
	// keep-sorted end
)

// MarshalJSON implements [json.Marshaler].
func (a Action) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

var actionByName = make(map[string]Action)

func init() {
	for a := ActionFail; a <= ActionOutput; a++ {
		actionByName[a.String()] = a
	}
}

// UnmarshalJSON implements [json.Unmarshaler].
func (a *Action) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	action, ok := actionByName[strings.ToLower(s)]
	if !ok {
		return fmt.Errorf("unmarshal action %q: %w", s, ErrInvalidAction)
	}
	*a = action
	return nil
}
