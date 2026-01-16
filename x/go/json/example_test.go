// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package json_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type Foo struct {
	F string
}

type Bar struct {
	B string
}

// Example_streamOfMixedTypes decodes a stream of JSON messages, one per line,
// of different types.
//
// Since JSON decoder skips unmatched fields, there is no way to tell whether
// Unamrshal of a valid JSON string into a given data structure was successful,
// unless the decoder is explicitly asked to disallow unknown fields, but this
// behavior is undesired due to backwards incompatibility.
//
// Moreover, the decoder operates on the reader, automatically advancing current
// position in the buffer. Instead of using alternative techniques to "unread"
// the reader and move the position back and forth to try Unmarshal different
// type of the message, this solution assumes that a single JSON message is
// serialized in the line.
//
// The example assumes that the two message types, Foo and Bar, have non-
// overlapping structures, e.g. no fields with the same name. It allows to
// check whether the "decoded" value has zero value, and as such try the next
// message type.
//
// NOTE: `go test -json ...`. The tool mixes build and test JSON
// streams in the standard output. See `go help buildjson` and
// `go doc test2json`. Unfortunately, the two events, BuildEvent and TestEvent,
// share the same fields with the same type, Action and Output, rendering this
// technique to detect failed JSON Unmarshal invalid.
func Example_streamOfMixedTypes() {
	str := `
{"F":"test-foo"}
{"B":"test-bar"}
`
	r := strings.NewReader(str)
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Bytes()
		if len(line) == 0 {
			continue
		}
		f, b, err := decode(line)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		switch {
		case f != nil:
			fmt.Println("Foo:", f.F)
		case b != nil:
			fmt.Println("Bar:", b.B)
		}
	}
	// Output:
	// Foo: test-foo
	// Bar: test-bar
}

func decode(b []byte) (*Foo, *Bar, error) {
	foo := new(Foo)
	err := json.Unmarshal(b, foo)
	if err == nil && !reflect.ValueOf(*foo).IsZero() {
		return foo, nil, nil
	}
	bar := new(Bar)
	err = json.Unmarshal(b, bar)
	if err == nil && !reflect.ValueOf(*bar).IsZero() {
		return nil, bar, nil
	}
	return nil, nil, err
}
