// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tests

import (
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const newFileMode = 0644

// GoldenFile holds baseline data for regression tests.
type GoldenFile string

// Diff generates a difference between buf and contents of the golden file. It
// fails the test if it can't read the file.
func (f GoldenFile) Diff(t *testing.T, buf string) string {
	t.Helper()
	data, err := os.ReadFile(string(f))
	if err != nil {
		t.Fatalf("golden file %s: read: %v", f, err)
	}
	splitStrings := cmpopts.AcyclicTransformer("SplitString", func(s string) []string {
		return strings.Split(s, "\n")
	})
	return cmp.Diff(string(data), buf, splitStrings)
}

// Write updates contents of the golden file with data. It fails the test if it
// fails to write data.
func (f GoldenFile) Write(t *testing.T, data string) {
	t.Helper()
	err := os.WriteFile(string(f), []byte(data), newFileMode)
	if err != nil {
		t.Fatalf("golden file %s: write: %v", f, err)
	}
}
