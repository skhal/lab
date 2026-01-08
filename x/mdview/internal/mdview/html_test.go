// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mdview_test

import (
	"flag"
	"os"
	"testing"

	"github.com/skhal/lab/go/tests"
	"github.com/skhal/lab/x/mdview/internal/mdview"
)

var update = flag.Bool("update", false, "update golden files")

func TestToHTML(t *testing.T) {
	markdown := mustReadFile(t, "testdata/basic.md")
	golden := tests.GoldenFile("testdata/basic_golden.html")

	got, err := mdview.ToHTML(markdown)

	if err != nil {
		t.Fatalf("ToHTML() unexpected error: %s", err)
	}
	if *update {
		golden.Write(t, string(got))
	}
	if diff := golden.Diff(t, string(got)); diff != "" {
		t.Errorf("ToHTML() mismatch (-want, +got):\n%s", diff)
	}
}

func mustReadFile(t *testing.T, filename string) []byte {
	t.Helper()
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
