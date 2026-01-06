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

var ErrInvalidAction = errors.New("invalid actin")

// Ref: `go help buildjson`.
type Event struct {
	ImportPath string `json:",omitempty"`
	Action     Action `json:",omitzero"`
	Output     string `json:",omitempty"`
}

//go:generate stringer -type=Action
type Action int

const (
	_            Action = iota
	ActionOutput        // the toolchain printed output
	ActionFail          // the build failed
)

func (a Action) MarshalJSON() ([]byte, error) {
	var s string
	switch a {
	default:
		return nil, fmt.Errorf("action %d: %w", a, ErrInvalidAction)
	case ActionOutput:
		s = "build-output"
	case ActionFail:
		s = "build-fail"
	}
	return json.Marshal(s)
}

func (a *Action) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	default:
		return fmt.Errorf("unmarshal action %q: %w", s, ErrInvalidAction)
	case "build-output":
		*a = ActionOutput
	case "build-fail":
		*a = ActionFail
	}
	return nil
}
