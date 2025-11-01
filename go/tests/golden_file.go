// Copyright 2025 Samvel Khalatyan. All rights reserved.

package tests

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
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
	return cmp.Diff(string(data), buf)
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
