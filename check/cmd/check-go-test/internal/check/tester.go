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

var defaultArgs = []string{"test", "-json", "-vet=all"}

// Tester runs go-test on a list of packages and processes JSON output. It
// groups events by event ids, which is a package and optional test case.
type Tester struct {
	coverage Coverage
}

// Test runs go-test for multiple packages.
func (t *Tester) Test(pkgs []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	args := t.args(pkgs)
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
	return t.processOutput(stdout)
}

func (t *Tester) args(pos []string) []string {
	args := defaultArgs
	if t.coverage > 0 {
		args = append(args, "-cover")
	}
	return append(args, pos...)
}

func (t *Tester) processOutput(r io.Reader) error {
	var (
		buildEvents = make(map[EventID][]*BuildEvent)
		testEvents  = make(map[EventID][]*TestEvent)
		errs        []error
	)
	for id, event := range decodeEvents(r) {
		switch e := event.(type) {
		case *TestEvent:
			testEvents[id] = append(testEvents[id], e)
			if e.Fail() {
				errs = append(errs, TestError(testEvents[id]))
			}
			if err := t.checkCoverage(e); err != nil {
				errs = append(errs, err)
			}
		case *BuildEvent:
			buildEvents[id] = append(buildEvents[id], e)
			if e.Fail() {
				errs = append(errs, BuildError(buildEvents[id]))
			}
		}
	}
	return errors.Join(errs...)
}

func (t *Tester) checkCoverage(e *TestEvent) error {
	if e.Coverage == nil || *e.Coverage >= t.coverage {
		return nil
	}
	return &CoverageError{
		Package: e.Package,
		Got:     *e.Coverage,
		Want:    t.coverage,
	}
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
	return NewTestEvent(e)
}
