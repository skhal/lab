// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"iter"
	"os"
	"os/exec"

	"github.com/skhal/lab/check/cmd/check-go-test/internal/build"
	"github.com/skhal/lab/check/cmd/check-go-test/internal/test"
)

// Tester runs `go test` on packages and groups events by event ids. It also
// keeps track of failed tests for further analysis.
type Tester struct {
	events   map[EventID][]Event
	coverage Coverage
}

// NewTester creates a tester, ready for testing packages.
func NewTester() *Tester {
	return &Tester{
		events: make(map[EventID][]Event),
	}
}

// Test runs a single `go test` for the packages.
func (t *Tester) Test(pkgs []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	args := append([]string{"test", "-json", "-vet=all"}, pkgs...)
	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Stderr = os.Stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	defer cmd.Wait()
	var errs []error
	for id, e := range decodeEvents(stdout) {
		ee := t.events[id]
		t.events[id] = append(ee, e)
		if e.Fail() {
			switch e.(type) {
			case *TestEvent:
				errs = append(errs, testError(ee))
			case *BuildEvent:
				errs = append(errs, buildError(ee))
			}
		}
	}
	return errors.Join(errs...)
}

func decodeEvents(r io.Reader) iter.Seq2[EventID, Event] {
	return func(yield func(EventID, Event) bool) {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			b := scanner.Bytes()
			e, err := JSONUnmarshal(b)
			if err != nil {
				break
			}
			if !yield(e.ID(), e) {
				break
			}
		}
	}
}

var (
	jsonPrefix       = []byte(`{`)
	buildEventMarker = []byte(`"Action":"build-`)
)

// ErrNotJSON means the output line was not a JSON object.
var ErrNotJSON = errors.New("not JSON")

// JSONUnmarshal decodes b into a build or test event.
func JSONUnmarshal(b []byte) (Event, error) {
	// Build output may include non-JSON lines `go help buildjson`
	if !bytes.HasPrefix(b, jsonPrefix) {
		return nil, fmt.Errorf("parse %q: %w", b, ErrNotJSON)
	}
	if bytes.Contains(b, buildEventMarker) {
		return jsonUnmarshalBuildEvent(b)
	}
	return jsonUnmarshalTestEvent(b)
}

func jsonUnmarshalBuildEvent(b []byte) (Event, error) {
	e := new(build.Event)
	if err := json.Unmarshal(b, &e); err != nil {
		return nil, err
	}
	return (*BuildEvent)(e), nil
}

func jsonUnmarshalTestEvent(b []byte) (Event, error) {
	e := new(test.Event)
	if err := json.Unmarshal(b, &e); err != nil {
		return nil, err
	}
	return &TestEvent{Event: e}, nil
}
