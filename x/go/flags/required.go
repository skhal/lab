// Copyright 2025 Samvel Khalatyan. All rights reserved.

package flags

// RequiredString is a flag that must be set from the command line.
// https://cs.opensource.google/go/go/+/master:src/cmd/go/internal/base/flag.go;l=33-38;drc=96d8d3eb3294e85972aed190aec1806ef3c30712
type RequiredString struct {
	val    string
	// parsed keeps track whether the flag was set.
	parsed bool
}

func (s *RequiredString) Set(v string) error {
	s.val = v
	s.parsed = true
	return nil
}

func (s *RequiredString) String() string {
	if s.parsed {
		return s.val
	}
	return "<not parsed>"
}
