// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/skhal/lab/infra/cmd/certval/pb"
	"google.golang.org/protobuf/encoding/prototext"
)

// ErrConfigUnmarshal means parse configuration as text-proto failed.
var ErrConfigUnmarshal = errors.New("unmarshal error")

// ParseConfig parses text-proto configuration from the file at path. It
// returns an error if the configuration can't be accessed or parsed.
func ParseConfig(path string) (*pb.Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, newConfigError(path, err)
	}
	cfg := new(pb.Config)
	if err := prototext.Unmarshal(b, cfg); err != nil {
		err = fmt.Errorf("%w: %s", ErrConfigUnmarshal, err)
		return nil, newConfigError(path, err)
	}
	return cfg, nil
}

type configError struct {
	err  error
	path string
}

func newConfigError(path string, err error) *configError {
	return &configError{
		path: path,
		err:  err,
	}
}

// Error implements [builtin.error] interface.
func (e *configError) Error() string {
	return fmt.Sprintf("config %s: %s", e.path, e.err.Error())
}

// Unwrap returns wrapped error for [errors.Unwrap].
func (e *configError) Unwrap() error {
	return e.err
}
