// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check_test

import (
	"testing"

	"github.com/skhal/lab/check/cmd/check-godoc/internal/check"
)

func TestIsTest(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "foo.go"},
		{name: "foo_test.go", want: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := check.IsTest(tc.name)

			if got != tc.want {
				t.Errorf("check.IsTest(%q) got %v; want %v", tc.name, got, tc.want)
			}
		})
	}
}
