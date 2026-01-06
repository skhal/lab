// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package check_test

import (
	"fmt"
	"testing"

	"github.com/skhal/lab/check/cmd/check-go-test/internal/check"
	"github.com/skhal/lab/check/cmd/check-go-test/internal/test"
)

func TestNewEventID(t *testing.T) {
	tests := []struct {
		name  string
		event *test.TestEvent
		want  check.EventID
	}{
		{
			name:  "no package",
			event: &test.TestEvent{Test: "test"},
			want:  check.EventID("test"),
		},
		{
			name:  "with package",
			event: &test.TestEvent{Package: "package", Test: "test"},
			want:  check.EventID("package/test"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := check.NewEventID(tc.event)

			if tc.want != got {
				t.Errorf("test.NewEventID(%s) = %s; want %s", (*printableEvent)(tc.event), got, tc.want)
			}
		})
	}
}

type printableEvent test.TestEvent

func (e *printableEvent) String() string {
	return fmt.Sprintf("{package:%q test:%q}", e.Package, e.Test)
}
