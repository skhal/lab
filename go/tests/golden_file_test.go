// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tests_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/skhal/lab/go/tests"
)

const newFileMode = 0644

func mustCreateFile(t *testing.T, file string, data string) string {
	t.Helper()
	tmpfile := filepath.Join(t.TempDir(), file)
	if err := os.WriteFile(tmpfile, []byte(data), newFileMode); err != nil {
		t.Fatalf("create file %s: %v", file, err)
		t.Logf("data:\n%s", data)
	}
	return tmpfile
}

func mustReadFile(t *testing.T, file string) string {
	t.Helper()
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("read file %s: %v", file, err)
	}
	return string(data)
}

func TestGodlenFile_Diff_empty(t *testing.T) {
	tmpfile := mustCreateFile(t, "golden.txt", "test-data")
	golden := tests.GoldenFile(tmpfile)
	data := "test-data"

	diff := golden.Diff(t, data)

	if diff != "" {
		t.Errorf("(tests.GoldenFile).Diff(%q) got diff; want empty diff", data)
		t.Logf("diff:\n%s", diff)
	}
}

func TestGodlenFile_Diff_notEmpty(t *testing.T) {
	tmpfile := mustCreateFile(t, "golden.txt", "test-data")
	golden := tests.GoldenFile(tmpfile)
	data := "test-data new"

	diff := golden.Diff(t, data)

	if diff == "" {
		t.Errorf("(tests.GoldenFile).Diff(%q) got empty diff; want diff", data)
	}
}

func TestGoldenFile_Write(t *testing.T) {
	tmpfile := filepath.Join(t.TempDir(), "golden.txt")
	golden := tests.GoldenFile(tmpfile)
	data := "test-data"

	golden.Write(t, data)

	if got := mustReadFile(t, tmpfile); got != data {
		t.Errorf("(tests.GoldenFile).Write(%q) unexpected write: %q", data, got)
	}
}

func TestGoldenFile_Write_fileExists(t *testing.T) {
	tmpfile := mustCreateFile(t, "golden.txt", "test-data exists")
	golden := tests.GoldenFile(tmpfile)
	data := "test-data"

	golden.Write(t, data)

	if got := mustReadFile(t, tmpfile); got != data {
		t.Errorf("(tests.GoldenFile).Write(%q) unexpected write: %q", data, got)
	}
}
