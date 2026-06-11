// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package query

import (
	"errors"
	"fmt"
	"strings"

	"github.com/skhal/lab/book/irex/pb"
	"github.com/skhal/lab/book/irex/query/feature/ping"
	"github.com/skhal/lab/book/irex/query/feature/plot"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	// ErrNoCommand means the query does not include a command.
	ErrNoCommand = errors.New("missing command")

	// ErrInvalidCommand means the query has unsupported command.
	ErrInvalidCommand = errors.New("invalid command")
)

type parseFunc func(q string) (*pb.Intent, error)

// commandParsers maps commands to a parser to handle parameters.
var commandParsers = map[string]parseFunc{
	"plot": dispatch(pb.E_PlotIntent_PlotIntent, plot.Parse),
	"ping": dispatch(pb.E_PingIntent_PingIntent, ping.Parse),
}

// Understand parses query q and generates an intent that best describes the
// query.
// It returns an error if query understanding fails.
func Understand(q string) (*pb.Intent, error) {
	trimmedQuery := strings.TrimSpace(strings.ToLower(q))
	if trimmedQuery == "" {
		return nil, ErrNoCommand
	}

	var cmdName, params string
	switch n := strings.Index(trimmedQuery, " "); n {
	case -1:
		// command only, no params
		cmdName = trimmedQuery
	default:
		cmdName = trimmedQuery[0:n]
		params = strings.TrimSpace(trimmedQuery[n:])
	}

	cmd, ok := commandParsers[cmdName]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrInvalidCommand, q)
	}
	return cmd(params)
}

// dispatch wraps a strongly-typed command parser that returns an intent of type
// Ret and wraps it inside the pb.Intent under extension ext.
func dispatch[Ret proto.Message](ext protoreflect.ExtensionType, f func(q string) (Ret, error)) parseFunc {
	return func(q string) (*pb.Intent, error) {
		msg, err := f(q)
		if err != nil {
			return nil, err
		}
		intent := pb.Intent_builder{}.Build()
		proto.SetExtension(intent, ext, msg)
		return intent, nil
	}
}
